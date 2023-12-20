terraform {
  required_providers {
    rhcs = {
      version = ">= 1.1.0"
      source  = "terraform.local/local/rhcs"
    }
  }
}

provider "rhcs" {
  url   = var.url
}

resource "rhcs_cluster_rosa_classic" "rosa_sts_cluster_import" {}
resource "rhcs_cluster_rosa_classic" "rosa_import_no_cluster_id" {}
resource "rhcs_identity_provider" "idp_google_import" {}
resource "rhcs_identity_provider" "idp_gitlab_import" {}
resource "rhcs_machine_pool" "mp_import" {}
