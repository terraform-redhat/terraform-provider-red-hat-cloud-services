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
    red-hat-cloud-services = {
      version = ">= 1.0.1"
      source  = "terraform-redhat/red-hat-cloud-services"
    }
  }
}

provider "red-hat-cloud-services" {
  token = var.token
  url   = var.url
}

resource "ocm_identity_provider" "google_idp" {
  cluster = var.cluster_id
  name    = "Google"
  google = {
    client_id     = var.google_client_id
    client_secret = var.google_client_secret
    hosted_domain = var.google_hosted_domain
  }
}
