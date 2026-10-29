# Backend For Gloomhaven Companion App

Currently very much a work in progress and documentation is lacking as a result.
This backend service is one of many used by my [React app FE](https://github.com/Zane-T-Rice/apps).

The goal is to make something to clear up some table space while playing Gloomhaven 2E.
(And also learn a bunch of technologies I wanted to play with.)

## Prerequisites

### For Just Local
- [Auth0](https://auth0.com)
- [aws-cli](https://aws.amazon.com/cli/). The seed script uses it against DynamoDB Local, so no AWS account is needed,
  but the CLI still wants credentials and a region configured. Any dummy values will do.
- [Go](https://go.dev/)
- [Nix](https://nixos.org/download/), which `scripts/start_dynamodb_local.sh` uses to get DynamoDB Local. If you do not
  use Nix, run DynamoDB Local on port 8000 any way you like instead.

A `.env` file with the following defined

```
AUDIENCE="YOUR_AUTH0_API_IDENTIFIER"
ISSUER="https://YOUR_ISSUER.us.auth0.com/"
LOCAL_SERVICE_PORT="7575"
WEBSITE_DOMAIN="http://localhost:3100"
API_GATEWAY_BASE_PATH="/gloomhaven-companion-service"
LOCAL_DATABASE_ENDPOINT="http://localhost:8000/"
GLOOMHAVEN_COMPANION_SERVICE_URL="http://localhost:7575"
WEB_SOCKETS_URL="http://localhost:8080"

# For seeding dynamodb table (see scripts/seed_dynamodb.sh).
# SEED_PLAYER_ID is optional. Set it to your Auth0 user id to also see the fake campaigns when logged in.
CLIENT_ID="YOUR_M2M_APPLICATION_CLIENT_ID"
CLIENT_SECRET="YOUR_M2M_APPLICATION_CLIENT_SECRET"
SEED_PLAYER_ID="google-oauth2|123456789"
```

### For Deploying
- [aws account](https://aws.amazon.com/) and credentials for it
- [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)


## Local Development

There are scripts in the /scripts directory. Run them from the root of the repository.

- build.sh  : Builds the Go binaries and puts them in the dist directory.
- local_websocket.sh : Runs websocket-local which is strictly for testing the service locally. websocket-local is never deployed to AWS.
- local.sh : Runs the service proper locally.
- seed_dynamodb.sh : Creates the service's required dynamodb table in the local DynamoDB and seeds the templates (the service needs to be running).
- start_dynamodb_local.sh : Starts DynamoDB Local (from the nixpkgs dynamodb-local package, no docker) in memory on port 8000.

To get everything up and running locally run the following from the root directory.

```
scripts/start_dynamodb_local.sh
scripts/local.sh
scripts/local_websocket.sh
# Once the others are up, create the table and seed the templates.
scripts/seed_dynamodb.sh
```

all_local.sh does the same in one go: it builds, then starts DynamoDB Local, the websocket service and the built
service binary each in its own alacritty window (change `TERMINAL` at the top of the script to use another terminal),
and seeds after a short wait.
If DynamoDB Local is not up in time, run seed_dynamodb.sh again.

## Templates

Templates are the stats that figures are made from. Each one is in an edition (a named group, i.e. "Gloomhaven 2E"
or a player's "Homebrew") and is either a system template or belongs to the player that made it.

- System templates are made by admins with `POST /templates?scope=system`. Only admins can change or delete them.
- Any player can make templates (`POST /templates`) and editions (`POST /templates/editions`), and change or delete their own.
- `GET /templates` and `GET /templates/editions` give the system's, the player's own and those of every other player
  in any of the player's campaigns, which are read only. Each has a `scope` of `system`, `mine` or `shared`.
- Every template needs an edition, and a template can only be in a system edition or one of its owner's own. A figure or a
  scenario made from a template remembers the edition (`editionId`), which is how same named templates in different
  editions are told apart.
- An edition that still has templates in it can not be deleted.

internal/routers/templates_test.go checks these rules. It needs DynamoDB Local running and is skipped without it.

## Deployment

The deployments/terraform directory has all the terraform I use to deploy the service
(both the service proper and the websockets service) to AWS. Start by reading
main.tf because there are some resources you need to create manually or edit after
the deploy for it all to work.

There is a script in the /scripts directory for deploying via the Terraform configuration.
Run the script from the root of the repository.

- deploy.sh : Runs build and terraform to deploy to AWS.

```
scripts/deploy.sh
```