#!/bin/bash

gcloud compute instances create canon \
    --machine-type n1-standard-8 \
    --image-family centos-7 --image-project centos-cloud \
    --address canon --tags http-server,https-server

gcloud compute disks resize canon --size 100GB

gcloud compute instances reset canon