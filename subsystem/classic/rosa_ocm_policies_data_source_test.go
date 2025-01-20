/*
Copyright (c) 2024 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package classic

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2/dsl/core"             // nolint
	. "github.com/onsi/gomega"                         // nolint
	. "github.com/onsi/gomega/ghttp"                   // nolint
	. "github.com/openshift-online/ocm-sdk-go/testing" // nolint
	. "github.com/terraform-redhat/terraform-provider-rhcs/subsystem/framework"
)

const (
	// This is the cluster that will be returned by the server when asked to retrieve a cluster with ID 123
	getStsPoliciesRequests = `{
		"items": [
		  {
			"kind": "STSPolicy",
			"id": "openshift_cloud_credential_operator_cloud_credential_operator_iam_ro_creds_policy",
			"details": "{}",
			"arn": "arn:aws:iam::000000000000:policy/ROSACloudCredentialOperator",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "openshift_cloud_network_config_controller_cloud_credentials_policy",
			"details":  "{}",
			"arn": "arn:aws:iam::000000000000:policy/ROSACloudNetworkConfigOperator",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "openshift_cluster_csi_drivers_ebs_cloud_credentials_policy",
			"details":  "{}",
			"arn": "arn:aws:iam::000000000000:policy/ROSAClusterCSIDriversEBSOperator",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "openshift_image_registry_installer_cloud_credentials_policy",
			"details":  "{}",
			"arn": "arn:aws:iam::000000000000:policy/ROSAImageRegistryOperator",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "openshift_ingress_operator_cloud_credentials_policy",
			"details":  "{}",
			"arn": "arn:aws:iam::000000000000:policy/ROSACloudIngressOperator",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "openshift_machine_api_aws_cloud_credentials_policy",
			"details":  "{}",
			"arn": "arn:aws:iam::000000000000:policy/ROSAMachineAPIOperator",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "shared_vpc_openshift_ingress_operator_cloud_credentials_policy",
			"details":  "{}",
			"arn": "",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "sts_installer_permission_policy",
			"details":  "{}",
			"arn": "arn:aws:iam::000000000000:policy/ROSAInstallerPolicy",
			"type": "AccountRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "sts_instance_controlplane_permission_policy",
			"details":  "{}",
			"arn": "arn:aws:iam::000000000000:policy/ROSAControlPlanePolicy",
			"type": "AccountRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "sts_instance_worker_permission_policy",
			"details":  "{}",
			"arn": "arn:aws:iam::000000000000:policy/ROSAWorkerPolicy",
			"type": "AccountRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "sts_support_permission_policy",
			"details":  "{}",
			"arn": "arn:aws:iam::000000000000:policy/ROSASRESupportPolicy",
			"type": "AccountRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "openshift_hcp_capa_controller_manager_credentials_policy",
			"details": "{}",
			"arn": "arn:aws:iam::aws:policy/service-role/ROSANodePoolManagementPolicy",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "openshift_hcp_cloud_network_config_controller_cloud_credentials_policy",
			"details": "{}",
			"arn": "arn:aws:iam::aws:policy/service-role/ROSACloudNetworkConfigOperatorPolicy",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "openshift_hcp_cluster_csi_drivers_ebs_cloud_credentials_policy",
			"details": "{}",
			"arn": "arn:aws:iam::aws:policy/service-role/ROSAAmazonEBSCSIDriverOperatorPolicy",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "openshift_hcp_control_plane_operator_credentials_policy",
			"details": "{}",
			"arn": "arn:aws:iam::aws:policy/service-role/ROSAControlPlaneOperatorPolicy",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "openshift_hcp_image_registry_installer_cloud_credentials_policy",
			"details": "{}",
			"arn": "arn:aws:iam::aws:policy/service-role/ROSAImageRegistryOperatorPolicy",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "openshift_hcp_ingress_operator_cloud_credentials_policy",
			"details": "{}",
			"arn": "arn:aws:iam::aws:policy/service-role/ROSAIngressOperatorPolicy",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "openshift_hcp_kms_provider_credentials_policy",
			"details": "{}",
			"arn": "arn:aws:iam::aws:policy/service-role/ROSAKMSProviderPolicy",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "openshift_hcp_kube_controller_manager_credentials_policy",
			"details": "{}",
			"arn": "arn:aws:iam::aws:policy/service-role/ROSAKubeControllerPolicy",
			"type": "OperatorRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "sts_hcp_installer_permission_policy",
			"details": "{}",
			"arn": "arn:aws:iam::aws:policy/service-role/ROSAInstallerPolicy",
			"type": "AccountRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "sts_hcp_instance_worker_permission_policy",
			"details": "{}",
			"arn": "arn:aws:iam::aws:policy/service-role/ROSAWorkerInstancePolicy",
			"type": "AccountRole"
		  },
		  {
			"kind": "STSPolicy",
			"id": "sts_hcp_support_permission_policy",
			"details": "{}",
			"arn": "arn:aws:iam::aws:policy/service-role/ROSASRESupportPolicy",
			"type": "AccountRole"
		  },
		  {
			"kind":"STSPolicy",
			"id":"sts_support_trust_policy",
			"details":"{\"Version\": \"2012-10-17\", \"Statement\": [{\"Action\": [\"sts:AssumeRole\"], \"Effect\": \"Allow\", \"Principal\": {\"AWS\": [\"arn:aws:iam::12345678912:role/RH-Technical-Support-12345678\"]}}]}",
			"arn":"",
			"type":"AccountRole"
		  }	  
		],
		"kind": "STSPoliciesList",
		"page": 1,
		"size": 23,
		"total": 23
	  }`
)

var _ = Describe("OCM policies data source", func() {

	It("Can list OCM policies", func() {
		// Prepare the server:
		TestServer.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/aws_inquiries/sts_policies"),
				RespondWithJSON(http.StatusOK, getStsPoliciesRequests),
			),
		)

		// Run the apply command:
		Terraform.Source(`
		  data "rhcs_policies" "my_policies" {
		  }
		`)
		runOutput := Terraform.Apply()
		Expect(runOutput.ExitCode).To(BeZero())

		// Check the state:
		resource := Terraform.Resource("rhcs_policies", "my_policies").(map[string]interface{})
		Expect(fmt.Sprint(resource["attributes"])).To(Equal(fmt.Sprint(
			map[string]interface{}{
				"operator_role_policies": map[string]interface{}{
					"openshift_cloud_credential_operator_cloud_credential_operator_iam_ro_creds_policy": "{}",
					"openshift_cloud_network_config_controller_cloud_credentials_policy":                "{}",
					"openshift_cluster_csi_drivers_ebs_cloud_credentials_policy":                        "{}",
					"openshift_image_registry_installer_cloud_credentials_policy":                       "{}",
					"openshift_ingress_operator_cloud_credentials_policy":                               "{}",
					"shared_vpc_openshift_ingress_operator_cloud_credentials_policy":                    "{}",
					"openshift_machine_api_aws_cloud_credentials_policy":                                "{}",
					"openshift_aws_vpce_operator_avo_aws_creds_policy":                                  nil,
				},
				"account_role_policies": map[string]interface{}{
					"sts_installer_permission_policy":             "{}",
					"sts_support_permission_policy":               "{}",
					"sts_instance_worker_permission_policy":       "{}",
					"sts_instance_controlplane_permission_policy": "{}",
					"sts_support_rh_sre_role":                     "arn:aws:iam::12345678912:role/RH-Technical-Support-12345678",
				},
			},
		)))
	})
})
