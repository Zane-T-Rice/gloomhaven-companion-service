# Backend For Gloomhaven Companion App

Currently very much a work in progress and documentation is lacking as a result.
This backend service is one of many used by my [React app FE](https://github.com/Zane-T-Rice/apps).

The goal is to make something to clear up some table space while playing Gloomhaven 2E.
(And also learn a bunch of technologies I wanted to play with.)

## Prerequisites

### For Just Local
- [Auth0](https://auth0.com)
- [aws account](https://aws.amazon.com/)
- [aws-cli](https://aws.amazon.com/cli/)
- [Go](https://go.dev/)

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