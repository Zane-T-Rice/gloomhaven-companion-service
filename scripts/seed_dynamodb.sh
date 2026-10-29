source ./.env

AWS_PAGER="" aws dynamodb create-table \
    --table-name gloomhaven-companion-service \
    --attribute-definitions \
        AttributeName=parent,AttributeType=S \
        AttributeName=entity,AttributeType=S \
    --key-schema AttributeName=parent,KeyType=HASH AttributeName=entity,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST \
    --table-class STANDARD \
    --endpoint-url http://localhost:8000 \
    --global-secondary-indexes \
        "[
            {
                \"IndexName\": \"entity-index\",
                \"KeySchema\": [{\"AttributeName\":\"entity\",\"KeyType\":\"HASH\"},
                                {\"AttributeName\":\"parent\",\"KeyType\":\"RANGE\"}],
                \"Projection\":{
                    \"ProjectionType\":\"KEYS_ONLY\"
                }
            }
        ]"

TOKEN_RESPONSE=`curl --request POST \
  --url "$ISSUER"oauth/token \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data grant_type=client_credentials \
  --data client_id=$CLIENT_ID \
  --data client_secret=$CLIENT_SECRET \
  --data audience=$AUDIENCE`;

ACCESS_TOKEN=`echo $TOKEN_RESPONSE | grep -Eo '"access_token"[^,]*' | grep -Eo '[^:]*$' | tr -d '"'`;

for i in $(ls templates/enemies/*); do \
  curl --request POST \
    --url "$GLOOMHAVEN_COMPANION_SERVICE_URL/templates?scope=system" \
    --header "Authorization: Bearer $ACCESS_TOKEN" \
    --header 'Content-Type: application/json' \
    --data "@$i";
done

for i in $(ls templates/characters/*); do \
  curl --request POST \
    --url "$GLOOMHAVEN_COMPANION_SERVICE_URL/templates?scope=system" \
    --header "Authorization: Bearer $ACCESS_TOKEN" \
    --header 'Content-Type: application/json' \
    --data "@$i";
done

for i in $(ls templates/summons/*); do \
  curl --request POST \
    --url "$GLOOMHAVEN_COMPANION_SERVICE_URL/templates?scope=system" \
    --header "Authorization: Bearer $ACCESS_TOKEN" \
    --header 'Content-Type: application/json' \
    --data "@$i";
done

for i in $(ls templates/bosses/*); do \
  curl --request POST \
    --url "$GLOOMHAVEN_COMPANION_SERVICE_URL/templates?scope=system" \
    --header "Authorization: Bearer $ACCESS_TOKEN" \
    --header 'Content-Type: application/json' \
    --data "@$i";
done


# Fake campaigns, each with a party of four players and a few fake scenarios, to have something
# to look at locally. The players are added as allies to every scenario, like adding a player with
# the Add Player button does, using the character template stats for the scenario level.
# The campaigns are created as the machine-to-machine application the token above is for.
# So that they also show up when logged in to the website, set SEED_PLAYER_ID in .env to
# your Auth0 user id (i.e. SEED_PLAYER_ID="google-oauth2|123456789") and the campaigns are
# linked to that player too. Running this again creates the campaigns again.

# Prints the id at the end of the "entity" of a created item.
created_id() {
  grep -Eo "\"entity\":\"#$1#[^\"]*\"" | grep -Eo "[^#]*\"$" | tr -d "\""
}

# The ally for a character class (i.e. "Silent Knife") at a character level.
character_ally() {
  node -e "
    const [file, level] = process.argv.slice(1);
    const { character } = require(file).stats[level];
    console.log(JSON.stringify({ ...character, alignment: \"ally\" }));
  " "$PWD/templates/characters/$(echo "$1" | tr "A-Z " "a-z_").json" "$2"
}

# seed_campaign "campaign name" "comma separated player classes" scenario...
# Each scenario is "scenario name|level|comma separated enemy groups".
seed_campaign() {
  local name=$1
  local players=$2
  shift 2
  local campaign_id=`curl -s --request POST \
    --url $GLOOMHAVEN_COMPANION_SERVICE_URL/campaigns \
    --header "Authorization: Bearer $ACCESS_TOKEN" \
    --header "Content-Type: application/json" \
    --data "{\"name\": \"$name\"}" | created_id CAMPAIGN`;

  if [ -n "$SEED_PLAYER_ID" ]; then
    AWS_PAGER="" aws dynamodb put-item \
      --table-name gloomhaven-companion-service \
      --endpoint-url http://localhost:8000 \
      --item "{\"parent\": {\"S\": \"#PLAYER#$SEED_PLAYER_ID\"}, \"entity\": {\"S\": \"#CAMPAIGN#$campaign_id\"}}";
  fi

  for scenario in "$@"; do
    IFS="|" read -r scenario_name scenario_level scenario_groups <<< "$scenario"
    local scenario_id=`curl -s --request POST \
      --url $GLOOMHAVEN_COMPANION_SERVICE_URL/campaigns/$campaign_id/scenarios \
      --header "Authorization: Bearer $ACCESS_TOKEN" \
      --header "Content-Type: application/json" \
      --data "{\"name\": \"$scenario_name\", \"scenarioLevel\": $scenario_level, \"groups\": \"$scenario_groups\"}" | created_id SCENARIO`;

    # Characters have levels 1 to 9. Level 0 scenarios just use level 1 characters.
    local character_level=$(( scenario_level > 0 ? scenario_level : 1 ))
    IFS="," read -ra classes <<< "$players"
    for class in "${classes[@]}"; do
      class=$(echo "$class" | xargs)
      curl -s --request POST \
        --url $GLOOMHAVEN_COMPANION_SERVICE_URL/campaigns/$campaign_id/scenarios/$scenario_id/figures \
        --header "Authorization: Bearer $ACCESS_TOKEN" \
        --header "Content-Type: application/json" \
        --data "$(character_ally "$class" $character_level)" > /dev/null;
    done
  done
}

seed_campaign "Frosthaven Crew" "Berserker, Bladeswarm, Quartermaster, Sawbones" \
  "1. Black Barrow|1|Living Bones, Living Corpse, Cave Bear" \
  "2. Barrow Lair|2|Living Spirit, Ooze, Giant Viper" \
  "3. Inox Encampment|2|Inox Guard, Inox Archer, Inox Priest"

seed_campaign "Sunday Night Gloomhaven" "Bruiser, Tinkerer, Spellweaver, Mindthief" \
  "1. Gloomhaven Warehouse|1|City Guard, City Archer, Hound" \
  "2. Diamond Mine|3|Earth Demon, Cave Bear, Crystal Rot" \
  "3. Gloomhaven Square|4|Cultist, Night Demon, Flame Demon" \
  "4. Temple of the Seer|5|Lurker Soldier, Deep Terror, Frost Demon"

seed_campaign "Solo Practice" "Cragheart, Doomstalker, Sunkeeper, Soothsinger" \
  "1. Practice Run|0|Bandit Scout, Bandit Archer" \
  "2. Rematch|1|Black Imp, Forest Imp, Chaos Demon"
