aws dynamodb create-table \
    --table-name gloomhaven-companion-service \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
        AttributeName=entity,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH AttributeName=entity,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST \
    --table-class STANDARD \
    --endpoint-url http://localhost:8000