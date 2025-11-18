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

for i in $(ls templates/*); do \
  curl --request POST \
    --url $GLOOMHAVEN_COMPANION_SERVICE_URL/templates \
    --header "Authorization: Bearer $ACCESS_TOKEN" \
    --header 'Content-Type: application/json' \
    --data "@$i";
done