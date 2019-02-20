GOOS=linux GOARCH=386 go build

gcloud compute scp svc canon:~
gcloud compute ssh canon --command="sudo mv svc /usr/sbin/canon"

gcloud compute ssh canon --command="sudo systemctl restart canon.service"