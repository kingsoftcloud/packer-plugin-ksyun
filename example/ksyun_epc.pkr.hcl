source "ksyun-kec" "example_epc" {
  access_key                  = var.access_key
  secret_key                  = var.secret_key
  region                      = "cn-beijing-6"
  source_image_id             = "eb8c0428-476e-49af-8ccb-9fad2455a54c"
  host_type                   = "EC-I-III-II"
  availability_zone           = "cn-beijing-6c"
  raid                        = "Raid1"
  ssh_username                = "root"
  ssh_clear_authorized_keys   = true
  associate_public_ip_address = true
}
build {
  name    = "ksyun_epc"
  sources = ["source.ksyun-kec.example_epc"]
  provisioner "shell" {
    inline = [
      "sleep 30",
      "yum install mysql -y"
    ]
  }
}
