#!/bin/bash

gcloud compute ssh canon --command="sudo systemctl daemon-reload"
gcloud compute ssh canon --command="sudo systemctl enable elasticsearch.service"
