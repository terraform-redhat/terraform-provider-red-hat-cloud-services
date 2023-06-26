#
# Copyright (c***REMOVED*** 2023 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License"***REMOVED***;
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.20.0"
    }
    ocm = {
      version = ">=1.0.1"
      source  = "terraform-redhat/ocm"
    }
  }
}
provider "ocm" {
  token = var.token
  url   = var.url
}

provider "aws" {
  region = var.cloud_region
}

# Generates the OIDC config resources' names
resource "ocm_rosa_oidc_config_input" "oidc_input" {
  count = var.managed ? 0 : 1

  region = var.cloud_region
}

# Create the OIDC config resources on AWS
module "oidc_config_input_resources" {
  count = var.managed ? 0 : 1

  source  = "terraform-redhat/rosa-sts/aws"
  version = "0.0.8"

  create_oidc_config_resources = true

  bucket_name             = one(ocm_rosa_oidc_config_input.oidc_input[*].bucket_name***REMOVED***
  discovery_doc           = one(ocm_rosa_oidc_config_input.oidc_input[*].discovery_doc***REMOVED***
  jwks                    = one(ocm_rosa_oidc_config_input.oidc_input[*].jwks***REMOVED***
  private_key             = one(ocm_rosa_oidc_config_input.oidc_input[*].private_key***REMOVED***
  private_key_file_name   = one(ocm_rosa_oidc_config_input.oidc_input[*].private_key_file_name***REMOVED***
  private_key_secret_name = one(ocm_rosa_oidc_config_input.oidc_input[*].private_key_secret_name***REMOVED***
}

resource "ocm_rosa_oidc_config" "oidc_config" {
  managed            = var.managed
  secret_arn         = one(module.oidc_config_input_resources[*].secret_arn***REMOVED***
  issuer_url         = one(ocm_rosa_oidc_config_input.oidc_input[*].issuer_url***REMOVED***
  installer_role_arn = var.installer_role_arn
}

data "ocm_rosa_operator_roles" "operator_roles" {
  operator_role_prefix = var.operator_role_prefix
  account_role_prefix  = var.account_role_prefix
}

module "operator_roles_and_oidc_provider" {
  source  = "terraform-redhat/rosa-sts/aws"
  version = "0.0.8"

  create_operator_roles = true
  create_oidc_provider  = true

  cluster_id                  = ""
  rh_oidc_provider_thumbprint = ocm_rosa_oidc_config.oidc_config.thumbprint
  rh_oidc_provider_url        = ocm_rosa_oidc_config.oidc_config.oidc_endpoint_url
  operator_roles_properties   = data.ocm_rosa_operator_roles.operator_roles.operator_iam_roles
  tags                        = var.tags
  path                        = var.path
}
