variable "token" {
  type      = string
  sensitive = true
}

variable "operator_role_prefix" {
  type = string
}

variable "url" {
  type    = string
  default = "https://api.openshift.com"
}

variable "account_role_prefix" {
  type = string
}

variable "cluster_name" {
  type    = string
  default = "my-cluster"
}

variable "cloud_region" {
  type    = string
  default = "us-east-2"
}

variable "availability_zones" {
  type    = list(string***REMOVED***
  default = ["us-east-2a"]
}
