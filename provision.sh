#!/bin/bash

# TODO need proxy to port 9000
gcloud compute instances create canon \
    --image-family centos-7 --image-project centos-cloud \
    --address canon --tags http-server,https-server \
    --create-disk=size=100GB