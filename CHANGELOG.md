# Latest Release

Please refer to [releases](https://github.com/kingsoftcloud/packer-plugin-ksyun/releases) for the latest CHANGELOG information.

---
## 0.0.13 (Jun 28, 2023)

### FEATURES

* Support `image_share_account` and `image_warm_up`, sharing a custom image to other accounts and support setting the image as warm-up that will accelerate instance boot. 

## 0.0.12 (May 31, 2023)

### FEATURES

* Support `image_copy_regions` and `image_copy_names`, copying a custom image helps you quickly deploy the operating environment of the current KEC instance to another region.


## 0.0.11 (May 5, 2023)

### FEATURES
* Support creating kec image include all snapshots of data disks
* Support `image_ignore_data_disks`, if the value is true, the key image created will not include any snapshot of data disks.
* update packer-plugin-sdk version: 0.0.14 -> v0.4.0
