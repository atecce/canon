provider "cloudflare" {
    email = "root@atec.pub"
    token = "${file("/keybase/private/atec/etc/cloudflare/token")}"
}

resource "cloudflare_record" "subdomain" {
    name = "canon-dev"
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

resource "google_compute_firewall" "canon-dev-ssh" {

    name = "canon-dev-ssh"
    network = "default"
    target_tags = ["canon-dev"]

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

resource "google_compute_firewall" "canon-dev" {
    
    name = "canon-dev"
    network = "default"
    target_tags = ["canon-dev"]

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
    name = "canon-dev"
}

resource "google_compute_instance" "default" {
 
    name = "canon-dev"
    zone = "us-east1-b"

    network_interface = {
        network = "default"
        access_config = {
            nat_ip = "${google_compute_address.static.address}"
        }
    }
    machine_type = "g1-small"
    boot_disk = {
        initialize_params = {
            image = "centos-cloud/centos-7"
        }
    }

    tags = ["canon-dev"]

    provisioner "remote-exec" {
        connection = {
            type = "ssh"
            user = "atec"
            private_key = "${file("~/.ssh/google_compute_engine")}"
            timeout = "120s"
        }
        inline = [
            "sudo yum install -y wget",

            "sudo yum install -y https://prerelease.keybase.io/keybase_i386.rpm",
            "run_keybase",
            
            "sudo mkdir -p /etc/canon",

            "wget https://atec.keybase.pub/etc/sshd_config",
            "sudo mv sshd_config /etc/ssh/sshd_config",
            "sudo systemctl restart sshd.service",

            "sudo mkdir -p /var/canon",
            "rsync -ah --progress /keybase/public/atec/data/gutenberg/entities.tar.gz .",
            "sudo tar -xvf entities.tar.gz -C /var/canon/"
        ]
    }

    provisioner "file" {
        connection = {
            type = "ssh"
            user = "root"
            private_key = "${file("~/.ssh/google_compute_engine")}"
            timeout = "120s"
        }
        source = "/keybase/private/atec/etc/server.crt"
        destination = "/etc/canon/server.crt"
    }

    provisioner "file" {
        connection = {
            type = "ssh"
            user = "root"
            private_key = "${file("~/.ssh/google_compute_engine")}"
            timeout = "120s"
        }
        source = "/keybase/private/atec/etc/server.key"
        destination = "/etc/canon/server.key"
    }

    provisioner "file" {
        connection = {
            type = "ssh"
            user = "root"
            private_key = "${file("~/.ssh/google_compute_engine")}"
            timeout = "120s"
        }
        source = "../../canon.service"
        destination = "/etc/systemd/system/canon.service"
    }

    depends_on = ["google_compute_firewall.canon-dev"]
}
