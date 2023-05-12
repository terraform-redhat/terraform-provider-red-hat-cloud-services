/*
Copyright (c) 2021 Red Hat, Inc.

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

package provider

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2/dsl/core"             // nolint
	. "github.com/onsi/gomega"                         // nolint
	. "github.com/onsi/gomega/ghttp"                   // nolint
	. "github.com/openshift-online/ocm-sdk-go/testing" // nolint
)

var _ = Describe("Cluster creation", func() {
	// This is the cluster that will be returned by the server when asked to create or retrieve
	// a cluster.
	const template = `{
	  "id": "123",
	  "name": "my-cluster",
	  "region": {
	    "id": "us-west-1"
	  },
	  "multi_az": true,
	  "api": {
	    "url": "https://my-api.example.com"
	  },
	  "console": {
	    "url": "https://my-console.example.com"
	  },
	  "network": {
	    "machine_cidr": "10.0.0.0/16",
	    "service_cidr": "172.30.0.0/16",
	    "pod_cidr": "10.128.0.0/14",
	    "host_prefix": 23
	  },
	  "version": {
		  "id": "openshift-4.8.0"
	  }
	}`

	const templateReadyState = `{
	  "id": "123",
	  "name": "my-cluster",
	  "state": "ready",
	  "region": {
	    "id": "us-west-1"
	  },
	  "multi_az": true,
	  "api": {
	    "url": "https://my-api.example.com"
	  },
	  "console": {
	    "url": "https://my-console.example.com"
	  },
	  "network": {
	    "machine_cidr": "10.0.0.0/16",
	    "service_cidr": "172.30.0.0/16",
	    "pod_cidr": "10.128.0.0/14",
	    "host_prefix": 23
	  },
	  "version": {
		  "id": "openshift-4.8.0"
	  }
	}`

	const versionListPage1 = `{
	"kind": "VersionList",
	"page": 1,
	"size": 2,
	"total": 2,
	"items": [{
			"kind": "Version",
			"id": "openshift-v4.10.1",
			"href": "/api/clusters_mgmt/v1/versions/openshift-v4.10.1",
			"raw_id": "4.11.1"
		},
		{
			"kind": "Version",
			"id": "openshift-v4.10.1",
			"href": "/api/clusters_mgmt/v1/versions/openshift-v4.11.1",
			"raw_id": "4.11.1"
		}
	]
}`

	It("Creates basic cluster", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/versions"),
				RespondWithJSON(http.StatusOK, versionListPage1),
			),
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(`.name`, "my-cluster"),
				VerifyJQ(`.cloud_provider.id`, "aws"),
				VerifyJQ(`.region.id`, "us-west-1"),
				VerifyJQ(`.product.id`, "rosa"),
				RespondWithPatchedJSON(http.StatusCreated, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "thumbprint": "111111",
							  "role_arn": "",
							  "support_role_arn": "",
							  "instance_iam_roles" : {
								"master_role_arn" : "",
								"worker_role_arn" : ""
							  },
							  "operator_role_prefix" : "test"
						  }
					  }
					},
					{
					  "op": "add",
					  "path": "/nodes",
					  "value": {
						"compute": 3,
						"compute_machine_type": {
							"id": "r5.xlarge"
						}
					  }
					}]`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster_rosa_classic" "my_cluster" {
		    name           = "my-cluster"	
		    cloud_region   = "us-west-1"
			aws_account_id = "123"
			sts = {
				operator_role_prefix = "test"
				role_arn = "",
				support_role_arn = "",
				instance_iam_roles = {
					master_role_arn = "",
					worker_role_arn = "",
				}
			}
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())
	})

	Context("Test destroy cluster", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/versions"),
					RespondWithJSON(http.StatusOK, versionListPage1),
				),
				CombineHandlers(
					VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
					VerifyJQ(`.name`, "my-cluster"),
					VerifyJQ(`.cloud_provider.id`, "aws"),
					VerifyJQ(`.region.id`, "us-west-1"),
					VerifyJQ(`.product.id`, "rosa"),
					VerifyJQ(`.aws.sts.instance_iam_roles.master_role_arn`, ""),
					VerifyJQ(`.aws.sts.instance_iam_roles.worker_role_arn`, ""),
					VerifyJQ(`.aws.sts.operator_role_prefix`, "test"),
					VerifyJQ(`.aws.sts.role_arn`, ""),
					VerifyJQ(`.aws.sts.support_role_arn`, ""),
					VerifyJQ(`.aws.account_id`, "123"),
					RespondWithPatchedJSON(http.StatusCreated, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "thumbprint": "111111",
							  "role_arn": "",
							  "support_role_arn": "",
							  "instance_iam_roles" : {
								"master_role_arn" : "",
								"worker_role_arn" : ""
							  },
							  "operator_role_prefix" : "test"
						  }
					  }
					},
					{
					  "op": "add",
					  "path": "/nodes",
					  "value": {
						"compute": 3,
						"compute_machine_type": {
							"id": "r5.xlarge"
						}
					  }
					}]`),
				),
				CombineHandlers(
					VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/clusters/123"),
					RespondWithJSON(http.StatusOK, templateReadyState),
				),
				CombineHandlers(
					VerifyRequest(http.MethodDelete, "/api/clusters_mgmt/v1/clusters/123"),
					RespondWithJSON(http.StatusOK, templateReadyState),
				),
			)
		})

		It("Disable waiting in destroy resource", func() {
			terraform.Source(`
				  resource "ocm_cluster_rosa_classic" "my_cluster" {
					name           = "my-cluster"	
					cloud_region   = "us-west-1"
					aws_account_id = "123"
					disable_waiting_in_destroy = true
					sts = {
						operator_role_prefix = "test"
						role_arn = "",
						support_role_arn = "",
						instance_iam_roles = {
							master_role_arn = "",
							worker_role_arn = "",
						}
					}
				  }
			`)

			// it should return a warning so exit code will be "0":
			Expect(terraform.Apply()).To(BeZero())
			Expect(terraform.Destroy()).To(BeZero())

		})

		It("Wait in destroy resource but use the default timeout", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/clusters/123"),
					RespondWithJSON(http.StatusNotFound, template),
				),
			)
			terraform.Source(`
				  resource "ocm_cluster_rosa_classic" "my_cluster" {
					name           = "my-cluster"	
					cloud_region   = "us-west-1"
					aws_account_id = "123"
					sts = {
						operator_role_prefix = "test"
						role_arn = "",
						support_role_arn = "",
						instance_iam_roles = {
							master_role_arn = "",
							worker_role_arn = "",
						}
					}
				  }
			`)

			// it should return a warning so exit code will be "0":
			Expect(terraform.Apply()).To(BeZero())
			Expect(terraform.Destroy()).To(BeZero())
		})

		It("Wait in destroy resource and set timeout to a negative value", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/clusters/123"),
					RespondWithJSON(http.StatusNotFound, template),
				),
			)
			terraform.Source(`
				  resource "ocm_cluster_rosa_classic" "my_cluster" {
					name           = "my-cluster"	
					cloud_region   = "us-west-1"
					aws_account_id = "123"
					destroy_timeout = -1
					sts = {
						operator_role_prefix = "test"
						role_arn = "",
						support_role_arn = "",
						instance_iam_roles = {
							master_role_arn = "",
							worker_role_arn = "",
						}
					}
				  }
			`)

			// it should return a warning so exit code will be "0":
			Expect(terraform.Apply()).To(BeZero())
			Expect(terraform.Destroy()).To(BeZero())
		})

		It("Wait in destroy resource and set timeout to a positive value", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/clusters/123"),
					RespondWithJSON(http.StatusNotFound, template),
				),
			)
			terraform.Source(`
				  resource "ocm_cluster_rosa_classic" "my_cluster" {
					name           = "my-cluster"	
					cloud_region   = "us-west-1"
					aws_account_id = "123"
					destroy_timeout = 10
					sts = {
						operator_role_prefix = "test"
						role_arn = "",
						support_role_arn = "",
						instance_iam_roles = {
							master_role_arn = "",
							worker_role_arn = "",
						}
					}
				  }
			`)

			// it should return a warning so exit code will be "0":
			Expect(terraform.Apply()).To(BeZero())
			Expect(terraform.Destroy()).To(BeZero())
		})
	})

	It("Creates cluster with http proxy", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/versions"),
				RespondWithJSON(http.StatusOK, versionListPage1),
			),
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(`.name`, "my-cluster"),
				VerifyJQ(`.cloud_provider.id`, "aws"),
				VerifyJQ(`.region.id`, "us-west-1"),
				VerifyJQ(`.product.id`, "rosa"),
				VerifyJQ(`.proxy.http_proxy`, "http://proxy.com"),
				VerifyJQ(`.proxy.https_proxy`, "http://proxy.com"),
				VerifyJQ(`.additional_trust_bundle`, "123"),
				RespondWithPatchedJSON(http.StatusOK, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "thumbprint": "111111",
							  "role_arn": "",
							  "support_role_arn": "",
							  "instance_iam_roles" : {
								"master_role_arn" : "",
								"worker_role_arn" : ""
							  },
							  "operator_role_prefix" : "test"
						  }
					  }
					},
					{
					  "op": "add",
					  "path": "/proxy",
					  "value": {						  
						  "http_proxy" : "http://proxy.com",
						  "https_proxy" : "http://proxy.com"
					  }
					},
					{
					  "op": "add",
					  "path": "/",
					  "value": {
						  "additional_trust_bundle" : "123"
					  }
					},
					{
					  "op": "add",
					  "path": "/nodes",
					  "value": {
						"compute": 3,
						"compute_machine_type": {
							"id": "r5.xlarge"
						}
					  }
					}]`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster_rosa_classic" "my_cluster" {
		    name           = "my-cluster"	
		    cloud_region   = "us-west-1"
			aws_account_id = "123"
			proxy = {
				http_proxy = "http://proxy.com",
				https_proxy = "http://proxy.com",
				additional_trust_bundle = "123",
			}
			sts = {
				operator_role_prefix = "test"
				role_arn = "",
				support_role_arn = "",
				instance_iam_roles = {
					master_role_arn = "",
					worker_role_arn = "",
				}
			}
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())
	})

	It("Creates cluster with aws subnet ids & private link", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/versions"),
				RespondWithJSON(http.StatusOK, versionListPage1),
			),
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(`.name`, "my-cluster"),
				VerifyJQ(`.cloud_provider.id`, "aws"),
				VerifyJQ(`.region.id`, "us-west-1"),
				VerifyJQ(`.product.id`, "rosa"),
				VerifyJQ(`.aws.subnet_ids.[0]`, "id1"),
				VerifyJQ(`.aws.private_link`, true),
				VerifyJQ(`.nodes.availability_zones.[0]`, "az1"),
				VerifyJQ(`.api.listening`, "internal"),
				RespondWithPatchedJSON(http.StatusOK, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
						  "private_link": true,
						  "subnet_ids": ["id1", "id2", "id3"],	
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "thumbprint": "111111",
							  "role_arn": "",
							  "support_role_arn": "",
							  "instance_iam_roles" : {
								"master_role_arn" : "",
								"worker_role_arn" : ""
							  },
							  "operator_role_prefix" : "test"
						  }
					  }
					},
					{
						"op": "add",
						"path": "/availability_zones",
						"value": ["az1", "az2", "az3"]
					},
					{
					  "op": "add",
					  "path": "/nodes",
					  "value": {
						"compute": 3,
						"compute_machine_type": {
							"id": "r5.xlarge"
						}
					  }
					}]`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster_rosa_classic" "my_cluster" {
		    name           = "my-cluster"
		    cloud_region   = "us-west-1"
			aws_account_id = "123"
			availability_zones = ["az1","az2","az3"]
			aws_private_link = true
			aws_subnet_ids = [
				"id1", "id2", "id3"
			]
			sts = {
				operator_role_prefix = "test"
				role_arn = "",
				support_role_arn = "",
				instance_iam_roles = {
					master_role_arn = "",
					worker_role_arn = "",
				}
			}
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())
	})

	It("Creates cluster when private link is false", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/versions"),
				RespondWithJSON(http.StatusOK, versionListPage1),
			),
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(`.name`, "my-cluster"),
				VerifyJQ(`.cloud_provider.id`, "aws"),
				VerifyJQ(`.region.id`, "us-west-1"),
				VerifyJQ(`.product.id`, "rosa"),
				VerifyJQ(`.aws.private_link`, false),
				VerifyJQ(`.api.listening`, nil),
				RespondWithPatchedJSON(http.StatusOK, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
						  "private_link": false,
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "thumbprint": "111111",
							  "role_arn": "",
							  "support_role_arn": "",
							  "instance_iam_roles" : {
								"master_role_arn" : "",
								"worker_role_arn" : ""
							  },
							  "operator_role_prefix" : "test"
						  }
					  }
					},
					{
					  "op": "add",
					  "path": "/nodes",
					  "value": {
						"compute": 3,
						"compute_machine_type": {
							"id": "r5.xlarge"
						}
					  }
					}]`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster_rosa_classic" "my_cluster" {
		    name           = "my-cluster"	
		    cloud_region   = "us-west-1"
			aws_account_id = "123"
			aws_private_link = false
			sts = {
				operator_role_prefix = "test"
				role_arn = "",
				support_role_arn = "",
				instance_iam_roles = {
					master_role_arn = "",
					worker_role_arn = "",
				}
			}
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())
	})

	It("Creates rosa sts cluster with autoscaling and update the default machine pool ", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/versions"),
				RespondWithJSON(http.StatusOK, versionListPage1),
			),
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(`.name`, "my-cluster"),
				VerifyJQ(`.cloud_provider.id`, "aws"),
				VerifyJQ(`.region.id`, "us-west-1"),
				VerifyJQ(`.product.id`, "rosa"),
				VerifyJQ(`.aws.sts.role_arn`, ""),
				VerifyJQ(`.aws.sts.support_role_arn`, ""),
				VerifyJQ(`.aws.sts.instance_iam_roles.master_role_arn`, ""),
				VerifyJQ(`.aws.sts.instance_iam_roles.worker_role_arn`, ""),
				VerifyJQ(`.aws.sts.operator_role_prefix`, "terraform-operator"),
				VerifyJQ(`.nodes.autoscale_compute.kind`, "MachinePoolAutoscaling"),
				VerifyJQ(`.nodes.autoscale_compute.max_replicas`, float64(4)),
				VerifyJQ(`.nodes.autoscale_compute.min_replicas`, float64(2)),
				VerifyJQ(`.nodes.compute_labels.label_key1`, "label_value1"),
				VerifyJQ(`.nodes.compute_labels.label_key2`, "label_value2"),
				RespondWithPatchedJSON(http.StatusOK, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "thumbprint": "111111",
							  "role_arn": "",
							  "support_role_arn": "",
							  "instance_iam_roles" : {
								"master_role_arn" : "",
								"worker_role_arn" : ""
							  },
							  "operator_role_prefix" : "terraform-operator"
						  }
					  }
					},
					{
					  "op": "add",
					  "path": "/nodes",
					  "value": {
						"autoscale_compute": {
							"min_replicas": 2,
							"max_replicas": 4
						},
						"compute_machine_type": {
							"id": "r5.xlarge"
						},
						"compute_labels": {
							"label_key1": "label_value1",
				    		"label_key2": "label_value2"
						}
					  }
					}
				  ]`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		resource "ocm_cluster_rosa_classic" "my_cluster" {
			name           = "my-cluster"	
			cloud_region   = "us-west-1"
			aws_account_id = "123"
			autoscaling_enabled = "true"
			min_replicas = "2"
			max_replicas = "4"
			default_mp_labels = {
				"label_key1" = "label_value1", 
				"label_key2" = "label_value2"
			}
			sts = {
				role_arn = "",
				support_role_arn = "",
				instance_iam_roles = {
				  master_role_arn = "",
				  worker_role_arn = ""
				},
				"operator_role_prefix" : "terraform-operator"
			}
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())

		// apply for update the min_replica from 2 to 3
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/clusters/123"),
				RespondWithPatchedJSON(http.StatusOK, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "thumbprint": "111111",
							  "role_arn": "arn:aws:iam::account-id:role/ManagedOpenShift-Installer-Role",
							  "support_role_arn": "arn:aws:iam::account-id:role/ManagedOpenShift-Support-Role",
							  "instance_iam_roles" : {
								"master_role_arn" : "arn:aws:iam::account-id:role/ManagedOpenShift-ControlPlane-Role",
								"worker_role_arn" : "arn:aws:iam::account-id:role/ManagedOpenShift-Worker-Role"
							  },
							  "operator_role_prefix" : "terraform-operator"
						  }
					  }
					},
					{
					  "op": "add",
					  "path": "/nodes",
					  "value": {
						"autoscale_compute": {
							"min_replicas": 2,
							"max_replicas": 4
						},
						"compute_machine_type": {
							"id": "r5.xlarge"
						},
						"compute_labels": {
							"label_key1": "label_value1",
				    		"label_key2": "label_value2"
						}
					  }
					}
				  ]`),
			),
			CombineHandlers(
				VerifyRequest(http.MethodPatch, "/api/clusters_mgmt/v1/clusters/123"),
				VerifyJQ(`.nodes.autoscale_compute.kind`, "MachinePoolAutoscaling"),
				VerifyJQ(`.nodes.autoscale_compute.max_replicas`, float64(4)),
				VerifyJQ(`.nodes.autoscale_compute.min_replicas`, float64(3)),
				RespondWithPatchedJSON(http.StatusOK, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "thumbprint": "111111",
							  "role_arn": "arn:aws:iam::account-id:role/ManagedOpenShift-Installer-Role",
							  "support_role_arn": "arn:aws:iam::account-id:role/ManagedOpenShift-Support-Role",
							  "instance_iam_roles" : {
								"master_role_arn" : "arn:aws:iam::account-id:role/ManagedOpenShift-ControlPlane-Role",
								"worker_role_arn" : "arn:aws:iam::account-id:role/ManagedOpenShift-Worker-Role"
							  },
							  "operator_role_prefix" : "terraform-operator"
						  }
					  }
					},
					{
					  "op": "add",
					  "path": "/nodes",
					  "value": {
						"autoscale_compute": {
							"min_replicas": 3,
							"max_replicas": 4
						},
						"compute_machine_type": {
							"id": "r5.xlarge"
						},
						"compute_labels": {
							"label_key1": "label_value1",
				    		"label_key2": "label_value2"
						}
					  }
					}
				  ]`),
			),
		)
		// Run the apply command:
		terraform.Source(`
		resource "ocm_cluster_rosa_classic" "my_cluster" {
			name           = "my-cluster"	
			cloud_region   = "us-west-1"
			aws_account_id = "123"
			autoscaling_enabled = "true"
			min_replicas = "3"
			max_replicas = "4"
			default_mp_labels = {
				"label_key1" = "label_value1", 
				"label_key2" = "label_value2"
			}
			sts = {
				role_arn = "arn:aws:iam::account-id:role/ManagedOpenShift-Installer-Role",
				support_role_arn = "arn:aws:iam::account-id:role/ManagedOpenShift-Support-Role",
				instance_iam_roles = {
				  master_role_arn = "arn:aws:iam::account-id:role/ManagedOpenShift-ControlPlane-Role",
				  worker_role_arn = "arn:aws:iam::account-id:role/ManagedOpenShift-Worker-Role"
				},
				"operator_role_prefix" : "terraform-operator"
			}
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())

		// apply for update the autoscaling group to compute nodes
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/clusters/123"),
				RespondWithPatchedJSON(http.StatusOK, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "thumbprint": "111111",
							  "role_arn": "arn:aws:iam::account-id:role/ManagedOpenShift-Installer-Role",
							  "support_role_arn": "arn:aws:iam::account-id:role/ManagedOpenShift-Support-Role",
							  "instance_iam_roles" : {
								"master_role_arn" : "arn:aws:iam::account-id:role/ManagedOpenShift-ControlPlane-Role",
								"worker_role_arn" : "arn:aws:iam::account-id:role/ManagedOpenShift-Worker-Role"
							  },
							  "operator_role_prefix" : "terraform-operator"
						  }
					  }
					},
					{
					  "op": "add",
					  "path": "/nodes",
					  "value": {
						"autoscale_compute": {
							"min_replicas": 3,
							"max_replicas": 4
						},
						"compute_machine_type": {
							"id": "r5.xlarge"
						},
						"compute_labels": {
							"label_key1": "label_value1",
				    		"label_key2": "label_value2"
						}
					  }
					}
				  ]`),
			),
			CombineHandlers(
				VerifyRequest(http.MethodPatch, "/api/clusters_mgmt/v1/clusters/123"),
				VerifyJQ(`.nodes.compute`, float64(4)),
				RespondWithPatchedJSON(http.StatusOK, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "thumbprint": "111111",
							  "role_arn": "arn:aws:iam::account-id:role/ManagedOpenShift-Installer-Role",
							  "support_role_arn": "arn:aws:iam::account-id:role/ManagedOpenShift-Support-Role",
							  "instance_iam_roles" : {
								"master_role_arn" : "arn:aws:iam::account-id:role/ManagedOpenShift-ControlPlane-Role",
								"worker_role_arn" : "arn:aws:iam::account-id:role/ManagedOpenShift-Worker-Role"
							  },
							  "operator_role_prefix" : "terraform-operator"
						  }
					  }
					},
					{
					  "op": "add",
					  "path": "/nodes",
					  "value": {
						"compute": 4,
						"compute_machine_type": {
							"id": "r5.xlarge"
						},
						"compute_labels": {
							"label_key1": "label_value1",
				    		"label_key2": "label_value2"
						}
					  }
					}
				  ]`),
			),
		)
		// Run the apply command:
		terraform.Source(`
		resource "ocm_cluster_rosa_classic" "my_cluster" {
			name           = "my-cluster"	
			cloud_region   = "us-west-1"
			aws_account_id = "123" 
			replicas = 4
			default_mp_labels = {
				"label_key1" = "label_value1", 
				"label_key2" = "label_value2"
			}
			sts = {
				role_arn = "arn:aws:iam::account-id:role/ManagedOpenShift-Installer-Role",
				support_role_arn = "arn:aws:iam::account-id:role/ManagedOpenShift-Support-Role",
				instance_iam_roles = {
				  master_role_arn = "arn:aws:iam::account-id:role/ManagedOpenShift-ControlPlane-Role",
				  worker_role_arn = "arn:aws:iam::account-id:role/ManagedOpenShift-Worker-Role"
				},
				"operator_role_prefix" : "terraform-operator"
			}
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())
	})

	It("Creates rosa sts cluster with OIDC Configuration ID", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/versions"),
				RespondWithJSON(http.StatusOK, versionListPage1),
			),
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(`.name`, "my-cluster"),
				VerifyJQ(`.cloud_provider.id`, "aws"),
				VerifyJQ(`.region.id`, "us-west-1"),
				VerifyJQ(`.product.id`, "rosa"),
				VerifyJQ(`.aws.sts.role_arn`, ""),
				VerifyJQ(`.aws.sts.support_role_arn`, ""),
				VerifyJQ(`.aws.sts.instance_iam_roles.master_role_arn`, ""),
				VerifyJQ(`.aws.sts.instance_iam_roles.worker_role_arn`, ""),
				VerifyJQ(`.aws.sts.operator_role_prefix`, "terraform-operator"),
				VerifyJQ(`.aws.sts.oidc_config.id`, "aaa"),
				RespondWithPatchedJSON(http.StatusOK, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "oidc_config": {
								"id": "aaa",
								"secret_arn": "aaa",
								"issuer_url": "https://oidc_endpoint_url",
								"reusable": true,
								"managed": false
							  },
							  "thumbprint": "111111",
							  "role_arn": "",
							  "support_role_arn": "",
							  "instance_iam_roles" : {
								"master_role_arn" : "",
								"worker_role_arn" : ""
							  },
							  "operator_role_prefix" : "terraform-operator"
						  }
					  }
					},
					{
						"op": "add",
						"path": "/nodes",
						"value": {
						  "compute": 3,
						  "compute_machine_type": {
							  "id": "r5.xlarge"
						  }
						}
					  }
				  ]`),
			),
		)
		// Run the apply command:
		terraform.Source(`
		resource "ocm_cluster_rosa_classic" "my_cluster" {
			name           = "my-cluster"
			cloud_region   = "us-west-1"
			aws_account_id = "123"
			sts = {
				role_arn = "",
				support_role_arn = "",
				instance_iam_roles = {
				  master_role_arn = "",
				  worker_role_arn = ""
				},
				"operator_role_prefix" : "terraform-operator",
				"oidc_config_id" = "aaa"
			}
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())
	})

	It("Fails to create cluster with incompatible account role's version and fail", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/versions"),
				RespondWithJSON(http.StatusOK, versionListPage1),
			),
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(`.name`, "my-cluster"),
				VerifyJQ(`.cloud_provider.id`, "aws"),
				VerifyJQ(`.region.id`, "us-west-1"),
				VerifyJQ(`.product.id`, "rosa"),
				RespondWithPatchedJSON(http.StatusCreated, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "thumbprint": "111111",
							  "role_arn": "arn:aws:iam::765374464689:role/terr-account-Installer-Role",
							  "support_role_arn": "",
							  "instance_iam_roles" : {
								"master_role_arn" : "",
								"worker_role_arn" : ""
							  },
							  "operator_role_prefix" : "test"
						  }
					  }
					},
					{
					  "op": "add",
					  "path": "/nodes",
					  "value": {
						"compute": 3,
						"compute_machine_type": {
							"id": "r5.xlarge"
						}
					  }
					}]`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster_rosa_classic" "my_cluster" {
		    name           = "my-cluster"	
		    cloud_region   = "us-west-1"
			aws_account_id = "123"
			version = "openshift-v4.12"
			sts = {
				operator_role_prefix = "test"
				role_arn = "arn:aws:iam::765374464689:role/terr-account-Installer-Role",
				support_role_arn = "",
				instance_iam_roles = {
					master_role_arn = "",
					worker_role_arn = "",
				}
			}
		  }
		`)
		// expect to get an error
		Expect(terraform.Apply()).ToNot(BeZero())
	})

	It("Create cluster with http token", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/versions"),
				RespondWithJSON(http.StatusOK, versionListPage1),
			),
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(`.aws.http_tokens_state`, "required"),
				RespondWithPatchedJSON(http.StatusCreated, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
                          "aws_http_tokens_state" : "required",
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "thumbprint": "111111",
							  "role_arn": "",
							  "support_role_arn": "",
							  "instance_iam_roles" : {
								"master_role_arn" : "",
								"worker_role_arn" : ""
							  },
							  "operator_role_prefix" : "test"
						  }
					  }
					},
					{
					  "op": "add",
					  "path": "/nodes",
					  "value": {
						"compute": 3,
						"compute_machine_type": {
							"id": "r5.xlarge"
						}
					  }
					}]`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster_rosa_classic" "my_cluster" {
		    name           = "my-cluster"	
		    cloud_region   = "us-west-1"
			aws_account_id = "123"
			aws_http_tokens_state = "required"
			sts = {
				operator_role_prefix = "test"
				role_arn = "",
				support_role_arn = "",
				instance_iam_roles = {
					master_role_arn = "",
					worker_role_arn = "",
				}
			}
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())
	})
	It("Fails to create cluster with http tokens and not supported version", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/versions"),
				RespondWithJSON(http.StatusOK, versionListPage1),
			),
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(`.aws.http_tokens_state`, "required"),
				RespondWithPatchedJSON(http.StatusCreated, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
                          "aws_http_tokens_state" : "required",
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "thumbprint": "111111",
							  "role_arn": "",
							  "support_role_arn": "",
							  "instance_iam_roles" : {
								"master_role_arn" : "",
								"worker_role_arn" : ""
							  },
							  "operator_role_prefix" : "test"
						  }
					  }
					},
					{
					  "op": "add",
					  "path": "/nodes",
					  "value": {
						"compute": 3,
						"compute_machine_type": {
							"id": "r5.xlarge"
						}
					  }
					}]`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster_rosa_classic" "my_cluster" {
		    name           = "my-cluster"	
		    cloud_region   = "us-west-1"
			aws_account_id = "123"
			aws_http_tokens_state = "required"
			version = "openshift-v4.10"
			sts = {
				operator_role_prefix = "test"
				role_arn = "",
				support_role_arn = "",
				instance_iam_roles = {
					master_role_arn = "",
					worker_role_arn = "",
				}
			}
		  }
		`)
		// expect to get an error
		Expect(terraform.Apply()).ToNot(BeZero())
	})
	It("Fails to create cluster with http tokens with not supported value", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/versions"),
				RespondWithJSON(http.StatusOK, versionListPage1),
			),
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(`.aws.http_tokens_state`, "bad_string"),
				RespondWithPatchedJSON(http.StatusCreated, template, `[
					{
					  "op": "add",
					  "path": "/aws",
					  "value": {
                          "aws_http_tokens_state" : "bad_string",
						  "sts" : {
							  "oidc_endpoint_url": "https://oidc_endpoint_url",
							  "thumbprint": "111111",
							  "role_arn": "",
							  "support_role_arn": "",
							  "instance_iam_roles" : {
								"master_role_arn" : "",
								"worker_role_arn" : ""
							  },
							  "operator_role_prefix" : "test"
						  }
					  }
					},
					{
					  "op": "add",
					  "path": "/nodes",
					  "value": {
						"compute": 3,
						"compute_machine_type": {
							"id": "r5.xlarge"
						}
					  }
					}]`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster_rosa_classic" "my_cluster" {
		    name           = "my-cluster"	
		    cloud_region   = "us-west-1"
			aws_account_id = "123"
			aws_http_tokens_state = "bad_string"
			version = "openshift-v4.12"
			sts = {
				operator_role_prefix = "test"
				role_arn = "",
				support_role_arn = "",
				instance_iam_roles = {
					master_role_arn = "",
					worker_role_arn = "",
				}
			}
		  }
		`)
		// expect to get an error
		Expect(terraform.Apply()).ToNot(BeZero())
	})
})
