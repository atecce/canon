provider "cloudflare" {
    email = "root@atec.pub"
    token = "${file("/keybase/private/atec/etc/cloudflare/token")}"
}

resource "cloudflare_record" "canon" {
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

resource "google_compute_firewall" "provisioner" {

    name = "provisioner"
    network = "default"
    target_tags = ["provisioner"]

    source_ranges = ["0.0.0.0/0"]
    allow = {
        protocol = "tcp"
        ports = ["22","80","443"]
    }
}

resource "google_compute_address" "static" {
    name = "canon"
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
    machine_type = "n1-standard-8"
    boot_disk = {
        initialize_params = {
            image = "centos-cloud/centos-7"
            size = 100
        }
    }

    tags = ["http-server", "https-server", "provisioner"]

    provisioner "remote-exec" {
        connection = {
            type = "ssh"
            user = "atec"
            private_key = "${file("~/.ssh/google_compute_engine")}"
            timeout = "120s"
        }
        inline = [
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
        source = "/keybase/public/atec/bin/proxy"
        destination = "/usr/sbin/proxy"
    }

    provisioner "file" {
        connection = {
            type = "ssh"
            user = "root"
            private_key = "${file("~/.ssh/google_compute_engine")}"
            timeout = "120s"
        }
        source = "/keybase/public/atec/etc/proxy.service"
        destination = "/etc/systemd/system/proxy.service"
    }

    provisioner "remote-exec" {
        connection = {
            type = "ssh"
            user = "root"
            private_key = "${file("~/.ssh/google_compute_engine")}"
            timeout = "120s"
        }
        inline = [
            "yum install -y perl-Digest-SHA java-sdk",

            "wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-6.5.1.rpm",
            "wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-6.5.1.rpm.sha512",
            "shasum -a 512 -c elasticsearch-6.5.1.rpm.sha512",
            "rpm --install elasticsearch-6.5.1.rpm",
            "yum install -y elasticsearch",

            "systemctl daemon-reload",
            "systemctl enable proxy.service",
            "systemctl enable elasticsearch.service",

            "chmod 755 /usr/sbin/proxy",
            "systemctl start proxy.service",
            "systemctl start elasticsearch.service",
        ]
    }
    depends_on = ["google_compute_firewall.provisioner"]
}