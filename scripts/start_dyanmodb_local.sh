(
    cd ../dynamo-local-admin-docker;
    docker pull instructure/dynamo-local-admin
    docker run -p 8000:8000 -it --rm instructure/dynamo-local-admin
)