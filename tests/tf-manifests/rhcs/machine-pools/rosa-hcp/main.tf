terraform {
  required_providers {
    rhcs = {
      version = ">= 1.1.0"
      source  = "terraform.local/local/rhcs"
    }
  }
}

provider "rhcs" {
}

locals {
  autoscaling = var.autoscaling_enabled == null && var.min_replicas == null && var.max_replicas == null ? null : {
    enabled      = var.autoscaling_enabled,
    min_replicas = var.min_replicas,
    max_replicas = var.max_replicas
  }
  aws_node_pool = {
    instance_type                 = var.machine_type,
    additional_security_group_ids = var.additional_security_groups,
    ec2_metadata_http_tokens      = var.ec2_metadata_http_tokens,
    tags                          = var.tags
    disk_size                     = var.disk_size
  }
}

resource "rhcs_hcp_machine_pool" "mps" {
  count                        = var.mp_count
  cluster                      = var.cluster
  name                         = var.mp_count == 1 ? var.name : "${var.name}-${count.index}"
  subnet_id                    = var.subnet_id
  labels                       = var.labels
  replicas                     = var.replicas
  taints                       = var.taints
  tuning_configs               = var.tuning_configs
  auto_repair                  = var.auto_repair
  upgrade_acknowledgements_for = var.upgrade_acknowledgements_for
  version                      = var.openshift_version
  autoscaling                  = local.autoscaling
  aws_node_pool                = local.aws_node_pool
  kubelet_configs              = var.kubelet_configs
}
