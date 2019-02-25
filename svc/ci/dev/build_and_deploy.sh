set -e

GOOS=linux GOARCH=386 go build

gcloud compute scp svc canon-dev:~
gcloud compute ssh canon-dev --command="sudo mv svc /usr/sbin/canon"

gcloud compute ssh canon-dev --command="sudo systemctl restart canon.service"