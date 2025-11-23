TERMINAL=alacritty
$TERMINAL -e scripts/build.sh
$TERMINAL -e scripts/start_dynamodb_local.sh &
$TERMINAL -e scripts/local_websocket.sh &
$TERMINAL -e ./dist/gloomhaven-companion-service/bootstrap &

# This one is tricky, but for now just using a sleep. It can always be re-run
# if the dynamodb local does not come up in time.
sleep 5
$TERMINAL -e scripts/seed_dynamodb.sh &