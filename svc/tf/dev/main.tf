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

resource "google_compute_firewall" "canon-dev" {
    
    name = "canon-dev"
    network = "default"
    target_tags = ["canon-dev"]

    source_ranges = ["141.158.1.238"]
    allow = {
        protocol = "tcp"
        ports = ["22","443"]
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
            "sudo mkdir -p /etc/canon",

            "sudo yum install -y wget",

            "wget https://atec.keybase.pub/etc/sshd_config",
            "sudo mv sshd_config /etc/ssh/sshd_config",
            "sudo systemctl restart sshd.service",
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
        source = "/keybase/public/atec/bin/canon"
        destination = "/usr/sbin/canon"
    }

    provisioner "file" {
        connection = {
            type = "ssh"
            user = "root"
            private_key = "${file("~/.ssh/google_compute_engine")}"
            timeout = "120s"
        }
        source = "canon.service"
        destination = "/etc/systemd/system/canon.service"
    }

    provisioner "remote-exec" {
        connection = {
            type = "ssh"
            user = "root"
            private_key = "${file("~/.ssh/google_compute_engine")}"
            timeout = "120s"
        }
        inline = [
            "chmod 755 /usr/sbin/canon",
            "systemctl start canon.service",
        ]
    }

    depends_on = ["google_compute_firewall.canon-dev"]
}