#!/bin/bash

gcloud compute ssh canon --command="sudo systemctl start proxy.service"
gcloud compute ssh canon --command="sudo systemctl start elasticsearch.service"
