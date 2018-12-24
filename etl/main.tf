provider "google" {
    credentials = "${file("/keybase/private/atec/etc/gcp/telos.json")}"
    project = "telos-162721"
    region = "us-east1"
    zone = "us-east1-b"
}


resource "google_compute_instance" "default" {

    name = "test"

    zone = "us-east1-b"

    machine_type = "n1-standard-1"

    boot_disk {
        initialize_params {
            image = "debian-cloud/debian-9"
        }
    }

    network_interface {
        network = "default"
    }

    # TODO
    # provisioner "remote-exec" {
    #     inline = [
    #         "wget https://atec.keybase.pub/bin/canon/etl > /usr/local/bin/etl"
    #     ]
    # }
}