#/bin/bash
docker run -d \
    -e PYSERVER_ALLOWED_ORIGIN="*" \
    -e PYSERVER_CLIENTID="" \
    -e GOOGLE_APPLICATION_CREDENTIALS="/creds.json" \
    --mount type=bind,source="${PWD}"/creds.json,target="/creds.json" \
    -p 5000:5000 \
    --name py-server-dev \
    py-server:latest