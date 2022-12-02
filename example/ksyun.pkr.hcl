source "ksyun-kec" "example" {
  access_key                  = var.access_key
  secret_key                  = var.secret_key
  region                      = "cn-beijing-6"
  availability_zone           = "cn-beijing-6a"
  instance_charge_type        = "HourlyInstantSettlement"
  image_name                  = "packer_test"
  source_image_id             = "IMG-dd1f8324-1f27-46e0-ad6b-b41d8c8ff025"
  instance_type               = "C4.2B"
  ssh_username                = "root"
  associate_public_ip_address = true
  public_ip_charge_type       = "DailyPaidByTransfer"
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
