variable "access_key" {
  type    = string
  default = env("KSYUN_ACCESS_KEY")
}

variable "secret_key" {
  type    = string
  default = env("KSYUN_SECRET_KEY")
}
