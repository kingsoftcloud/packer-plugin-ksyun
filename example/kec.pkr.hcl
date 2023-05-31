packer {
  required_plugins {
    ksyun = {
      version = "0.0.13"
      source  = "github.com/kingsoftcloud/ksyun"
    }
  }
}

variable ak {
  type    = string
  default = "${env("KSYUN_ACCESS_KEY")}"
}

variable sk {
  type    = string
  default = "${env("KSYUN_SECRET_KEY")}"
}

source "ksyun-kec" "test" {
  access_key      = var.ak
  secret_key      = var.sk
  region          = "cn-shanghai-2"
  image_name      = "packer_test"
  source_image_id = "IMG-12112384-c3d3-4d42-8882-58234825ba1c"
  instance_type   = "N3.1B"
  ssh_username    = "root"

  # 此参数用于跳过ssh
  # communicator                = "none"

  # 如需使用ssh，须保证网络能通，如果不在同一个网络环境下就要挂公网ip
  associate_public_ip_address = true

  ssh_clear_authorized_keys = true

  # 此参数为true时，data_disks的硬盘不会打快照加入镜像
  # image_ignore_data_disks = true

  data_disks {
    data_disk_type = "SSD3.0"
    data_disk_size = 50
  }

  # 复制镜像到以下region
  image_copy_regions = ["cn-beijing-6", "cn-guangzhou-1"]

  # 镜像复制后的名称, 不命名则使用原镜像的名称
  image_copy_names = ["copy-test"]

}

build {
  sources = ["source.ksyun-kec.test"]
  provisioner "shell" {
    inline = ["sleep 10", "df -h"]
  }
}