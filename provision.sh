#!/bin/bash

gcloud compute instances create canon \
    --image-family centos-7 --image-project centos-cloud \
    --address canon --tags http-server,https-server #\
# TODO    --create-disk=size=100GB