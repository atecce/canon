#!/bin/bash

# TODO distribute release?

gcloud compute ssh atec@canon --command="sudo yum install -y wget"
gcloud compute ssh atec@canon --command="sudo yum install -y perl-Digest-SHA"
gcloud compute ssh atec@canon --command="sudo yum install -y java-sdk"

gcloud compute ssh atec@canon --command="wget https://atec.keybase.pub/bin/proxy"
gcloud compute ssh atec@canon --command="wget https://atec.keybase.pub/etc/proxy.service"
gcloud compute ssh atec@canon --command="sudo mv /home/atec/proxy.service /etc/systemd/system/"
gcloud compute ssh atec@canon --command="sudo mv /home/atec/proxy /usr/sbin/"

gcloud compute ssh atec@canon --command="wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-6.5.1.rpm"
gcloud compute ssh atec@canon --command="wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-6.5.1.rpm.sha512"
gcloud compute ssh atec@canon --command="shasum -a 512 -c elasticsearch-6.5.1.rpm.sha512"
gcloud compute ssh atec@canon --command="sudo rpm --install elasticsearch-6.5.1.rpm"
gcloud compute ssh atec@canon --command="sudo yum install -y elasticsearch"

gcloud compute ssh atec@canon --command="sudo systemctl daemon-reload"
gcloud compute ssh atec@canon --command="sudo systemctl enable proxy.service"
gcloud compute ssh atec@canon --command="sudo systemctl enable elasticsearch.service"