aws dynamodb create-table \
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