source "ksyun-kec" "example" {
  access_key                  = var.access_key
  secret_key                  = var.secret_key
  region                      = "cn-shanghai-2"
  image_name                  = "packer_test"
  source_image_id             = "IMG-dd1f8324-1f27-46e0-ad6b-b41d8c8ff025"
  instance_type               = "N3.1B"
  ssh_username                = "root"
  associate_public_ip_address = true
}
build {
  name    = "ksyun"
  sources = ["source.ksyun-kec.example"]
  provisioner "shell" {
    inline = [
      "sleep 30",
      "yum install mysql -y"
    ]
  }
}
