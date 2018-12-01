#!/bin/bash

gcloud compute ssh canon --command="sudo yum install -y wget"
gcloud compute ssh canon --command="wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-6.5.1.rpm"
gcloud compute ssh canon --command="wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-6.5.1.rpm.sha512"
# TODO gcloud compute ssh canon --command="shasum -a 512 -c elasticsearch-6.5.1.rpm.sha512"
gcloud compute ssh canon --command="sudo rpm --install elasticsearch-6.5.1.rpm"
gcloud compute ssh canon --command="sudo yum install -y java-sdk"
gcloud compute ssh canon --command="sudo yum install -y elasticsearch"
