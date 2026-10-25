#!/usr/bin/env bash
# Runs DynamoDB Local on port 8000 from the nixpkgs package, so docker is not needed.
# The package has an unfree license, so allow unfree packages for just this command.
# The data is only kept in memory, so run seed_dynamodb.sh each time this is started.
NIXPKGS_ALLOW_UNFREE=1 nix-shell -p dynamodb-local --run \
    'dynamodb-local -inMemory -disableTelemetry -port 8000'
