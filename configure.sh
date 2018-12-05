#!/bin/bash

# TODO distribute release?
GOOS=linux GOARCH=386 go build github.com/atecce/proxy

gcloud compute scp proxy ~/go/src/github.com/atecce/proxy/proxy.service atec@canon:~
gcloud compute ssh atec@canon --command="sudo mv /home/atec/proxy.service /etc/systemd/system/"
gcloud compute ssh atec@canon --command="sudo mv /home/atec/proxy /usr/sbin/"

gcloud compute ssh atec@canon --command="sudo yum install -y wget"
gcloud compute ssh atec@canon --command="sudo yum install -y perl-Digest-SHA"
gcloud compute ssh atec@canon --command="wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-6.5.1.rpm"
gcloud compute ssh atec@canon --command="wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-6.5.1.rpm.sha512"
gcloud compute ssh atec@canon --command="shasum -a 512 -c elasticsearch-6.5.1.rpm.sha512"
gcloud compute ssh atec@canon --command="sudo yum install -y java-sdk"
gcloud compute ssh atec@canon --command="sudo rpm --install elasticsearch-6.5.1.rpm"
gcloud compute ssh atec@canon --command="sudo yum install -y elasticsearch"

gcloud compute ssh atec@canon --command="sudo systemctl daemon-reload"
gcloud compute ssh atec@canon --command="sudo systemctl enable proxy.service"
gcloud compute ssh atec@canon --command="sudo systemctl enable elasticsearch.service"