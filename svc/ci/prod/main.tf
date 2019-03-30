provider "cloudflare" {
    email = "root@atec.pub"
    token = "${file("/keybase/private/atec/etc/cloudflare/token")}"
}

resource "cloudflare_record" "subdomain" {
    name = "canon"
    domain = "atec.pub"
    type = "A"
    value = "${google_compute_address.static.address}"
    proxied = true
}

provider "google" {
    credentials = "${file("/keybase/private/atec/etc/gcp/telos.json")}"
    project = "telos-162721"
    region = "us-east1"
    zone = "us-east1-b"
}

resource "google_compute_firewall" "canon-ssh" {

    name = "canon-ssh"
    network = "default"
    target_tags = ["canon"]

    source_ranges = [
        "141.158.1.238",
        "100.19.46.101",
        "65.207.13.162"
    ]

    allow = {
        protocol = "tcp"
        ports = ["22"]
    }
}

resource "google_compute_firewall" "canon" {
    
    name = "canon"
    network = "default"
    target_tags = ["canon"]

    # https://www.cloudflare.com/ips-v4
    source_ranges = [
        "173.245.48.0/20",
        "103.21.244.0/22",
        "103.22.200.0/22",
        "103.31.4.0/22",
        "141.101.64.0/18",
        "108.162.192.0/18",
        "190.93.240.0/20",
        "188.114.96.0/20",
        "197.234.240.0/22",
        "198.41.128.0/17",
        "162.158.0.0/15",
        "104.16.0.0/12",
        "172.64.0.0/13",
        "131.0.72.0/22"
    ]
    allow = {
        protocol = "tcp"
        ports = ["443"]
    }
}

resource "google_compute_address" "static" {
    name = "canon"
}

variable "kb_key" {
    type = "string"
}

resource "google_compute_instance" "default" {
 
    name = "canon"
    zone = "us-east1-b"

    network_interface = {
        network = "default"
        access_config = {
            nat_ip = "${google_compute_address.static.address}"
        }
    }
    machine_type = "n1-standard-1"
    boot_disk = {
        initialize_params = {
            image = "centos-cloud/centos-7"
        }
    }

    tags = ["canon"]

    provisioner "remote-exec" {
        connection = {
            type = "ssh"
            user = "atec"
            private_key = "${file("~/.ssh/google_compute_engine")}"
            timeout = "120s"
        }
        inline = [
            "sudo yum install -y https://prerelease.keybase.io/keybase_i386.rpm",

            "run_keybase",
            "keybase oneshot -u atec --paperkey '${var.kb_key}'",

            "sudo mkdir -p /etc/canon/",

            "cp /keybase/private/atec/etc/server.crt . && sudo cp server.crt /etc/canon/",
            "cp /keybase/private/atec/etc/server.key . && sudo cp server.key /etc/canon/",

            "cp /keybase/public/atec/etc/canon.service . && sudo cp canon.service /etc/systemd/system/",
            "cp /keybase/public/atec/etc/yum.repos.d/mongodb-org-4.0.repo . && sudo cp mongodb-org-4.0.repo /etc/yum.repos.d/",

            "sudo yum install -y mongodb-org",
            "sudo systemctl start mongod",
            
            "rsync -ah --progress /keybase/public/atec/data/gutenberg/entities.bson.gz .",
            "mongorestore --archive=entities.bson.gz --gzip -vvvvv",

            "mongo << EOF",
            "use canon",
            "db.entities.createIndex({ author: 1 })",
            "EOF"

            # TODO build and deploy
        ]
    }
}
