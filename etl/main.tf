provider "google" {
    credentials = "${file("/keybase/private/atec/etc/gcp/telos.json")}"
    project = "telos-162721"
    region = "us-east1"
    zone = "us-east1-b"
}

resource "google_compute_firewall" "etl" {

    name = "etl"
    network = "default"
    target_tags = ["etl"]

    source_ranges = ["0.0.0.0/0"]
    allow = {
        protocol = "tcp"
        ports = ["22"]
    }
}

resource "google_compute_address" "static" {
    name = "etl"
}
  
resource "google_compute_instance" "default" {

    name = "etl"
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

    tags = ["etl"]

    provisioner "remote-exec" {
        connection = {
            type = "ssh"
            user = "atec"
            private_key = "${file("~/.ssh/google_compute_engine")}"
            timeout = "120s"
        }
        inline = [
            "sudo yum install -y perl-Digest-SHA java-sdk wget",

            "wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-6.5.1.rpm",
            "wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-6.5.1.rpm.sha512",
            "shasum -a 512 -c elasticsearch-6.5.1.rpm.sha512",
            "sudo rpm --install elasticsearch-6.5.1.rpm",
            "sudo yum install -y elasticsearch",

            "sudo systemctl daemon-reload",
            "sudo systemctl enable elasticsearch.service",

            "sudo systemctl start elasticsearch.service",
        ]
    }
}
