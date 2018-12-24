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
        ports = ["22"]
    }
}
    
resource "google_compute_instance" "default" {

    name = "etl"
    zone = "us-east1-b"
    tags = ["provisioner"]

    network_interface = {
        network = "default"
        access_config = {}
    }
    machine_type = "n1-standard-1"
    boot_disk = {
        initialize_params = {
            image = "debian-cloud/debian-9"
        }
    }

    provisioner "remote-exec" {
        connection = {
            type = "ssh"
            user = "atec"
            private_key = "${file("~/.ssh/google_compute_engine")}"
            timeout = "120s"
        }
        inline = [
            "wget https://atec.keybase.pub/bin/canon/etl",
            # TODO
            # "chmod 755 etl",
            # "./etl 2>etl.log &"
        ]
    }

    depends_on = ["google_compute_firewall.provisioner"]
}