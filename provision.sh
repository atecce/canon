#!/bin/bash

gcloud compute instances create canon \
    --machine-type n1-standard-8 --boot-disk-size 100GB \
    --image-family centos-7 --image-project centos-cloud \
    --address canon --tags http-server,https-server