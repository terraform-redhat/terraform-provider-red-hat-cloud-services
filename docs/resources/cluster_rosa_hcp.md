---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rhcs_cluster_rosa_hcp Resource - terraform-provider-rhcs"
subcategory: ""
description: |-
  OpenShift managed cluster using rosa sts.
---

# rhcs_cluster_rosa_hcp (Resource)

OpenShift managed cluster using rosa sts.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `aws_account_id` (String) Identifier of the AWS account. After the creation of the resource, it is not possible to update the attribute value.
- `aws_billing_account_id` (String) Identifier of the AWS account for billing. After the creation of the resource, it is not possible to update the attribute value.
- `cloud_region` (String) Cloud region identifier, for example 'us-east-1'.
- `name` (String) Name of the cluster. Cannot exceed 15 characters in length. After the creation of the resource, it is not possible to update the attribute value.

### Optional

- `autoscaling_enabled` (Boolean) Enable autoscaling for the initial worker pool. This attribute is specifically applies for the Worker Machine Pool and becomes irrelevant once the resource is created. Any modifications to the default Machine Pool should be made through the Terraform imported Machine Pool resource. For more details, refer to [Worker Machine Pool in ROSA Cluster](../guides/worker-machine-pool.md)
- `availability_zones` (List of String) Availability zones. This attribute is specifically applies for the Worker Node Pool and becomes irrelevant once the resource is created. Any modifications to the default Machine Pool should be made through the Terraform imported Machine Pool resource. For more details, refer to [Worker Node Pool in ROSA Cluster](../guides/worker-machine-pool.md)
- `aws_private_link` (Boolean) Provides private connectivity from your cluster's VPC to Red Hat SRE, without exposing traffic to the public internet. After the creation of the resource, it is not possible to update the attribute value.
- `aws_subnet_ids` (List of String) AWS subnet IDs. After the creation of the resource, it is not possible to update the attribute value.
- `channel_group` (String) Name of the channel group where you select the OpenShift cluster version, for example 'stable'. For ROSA, only 'stable' is supported. After the creation of the resource, it is not possible to update the attribute value.
- `compute_machine_type` (String) Identifies the machine type used by the initial worker nodes, for example `m5.xlarge`. Use the `rhcs_machine_types` data source to find the possible values. This attribute is specifically applies for the Worker Machine Pool and becomes irrelevant once the resource is created. Any modifications to the default Machine Pool should be made through the Terraform imported Machine Pool resource. For more details, refer to [Worker Machine Pool in ROSA Cluster](../guides/worker-machine-pool.md)
- `destroy_timeout` (Number) This value sets the maximum duration in minutes to allow for destroying resources. Default value is 60 minutes.
- `disable_waiting_in_destroy` (Boolean) Disable addressing cluster state in the destroy resource. Default value is false, and so a `destroy` will wait for the cluster to be deleted.
- `etcd_encryption` (Boolean) Encrypt etcd data. Note that all AWS storage is already encrypted. After the creation of the resource, it is not possible to update the attribute value.
- `external_id` (String) Unique external identifier of the cluster. After the creation of the resource, it is not possible to update the attribute value.
- `host_prefix` (Number) Length of the prefix of the subnet assigned to each node. After the creation of the resource, it is not possible to update the attribute value.
- `kms_key_arn` (String) The key ARN is the Amazon Resource Name (ARN) of a AWS Key Management Service (KMS) Key. It is a unique, fully qualified identifier for the AWS KMS Key. A key ARN includes the AWS account, Region, and the key ID(optional). After the creation of the resource, it is not possible to update the attribute value.
- `machine_cidr` (String) Block of IP addresses for nodes. After the creation of the resource, it is not possible to update the attribute value.
- `pod_cidr` (String) Block of IP addresses for pods. After the creation of the resource, it is not possible to update the attribute value.
- `properties` (Map of String) User defined properties.
- `proxy` (Attributes) proxy (see [below for nested schema](#nestedatt--proxy))
- `replicas` (Number) Number of worker/compute nodes to provision. Single zone clusters need at least 2 nodes, multizone clusters need at least 3 nodes. This attribute is specifically applies for the Worker Machine Pool and becomes irrelevant once the resource is created. Any modifications to the default Machine Pool should be made through the Terraform imported Machine Pool resource. For more details, refer to [Worker Machine Pool in ROSA Cluster](../guides/worker-machine-pool.md)
- `service_cidr` (String) Block of IP addresses for the cluster service network. After the creation of the resource, it is not possible to update the attribute value.
- `sts` (Attributes) STS configuration. (see [below for nested schema](#nestedatt--sts))
- `tags` (Map of String) Apply user defined tags to all cluster resources created in AWS. After the creation of the resource, it is not possible to update the attribute value.
- `upgrade_acknowledgements_for` (String) Indicates acknowledgement of agreements required to upgrade the cluster version between minor versions (e.g. a value of "4.12" indicates acknowledgement of any agreements required to upgrade to OpenShift 4.12.z from 4.11 or before).
- `version` (String) Desired version of OpenShift for the cluster, for example '4.11.0'. If version is greater than the currently running version, an upgrade will be scheduled.
- `wait_for_create_complete` (Boolean) Wait until the cluster is either in a ready state or in an error state. The waiter has a timeout of 60 minutes, with the default value set to false

### Read-Only

- `api_url` (String) URL of the API server.
- `console_url` (String) URL of the console.
- `current_version` (String) The currently running version of OpenShift on the cluster, for example '4.11.0'.
- `domain` (String) DNS domain of cluster.
- `id` (String) Unique identifier of the cluster.
- `ocm_properties` (Map of String) Merged properties defined by OCM and the user defined 'properties'.
- `state` (String) State of the cluster.

<a id="nestedatt--proxy"></a>
### Nested Schema for `proxy`

Optional:

- `additional_trust_bundle` (String) A string containing a PEM-encoded X.509 certificate bundle that will be added to the nodes' trusted certificate store.
- `http_proxy` (String) HTTP proxy.
- `https_proxy` (String) HTTPS proxy.
- `no_proxy` (String) No proxy.


<a id="nestedatt--sts"></a>
### Nested Schema for `sts`

Required:

- `instance_iam_roles` (Attributes) Instance IAM Roles (see [below for nested schema](#nestedatt--sts--instance_iam_roles))
- `operator_role_prefix` (String) Operator IAM Role prefix
- `role_arn` (String) Installer Role
- `support_role_arn` (String) Support Role

Optional:

- `oidc_config_id` (String) OIDC Configuration ID
- `oidc_endpoint_url` (String) OIDC Endpoint URL

Read-Only:

- `thumbprint` (String) SHA1-hash value of the root CA of the issuer URL

<a id="nestedatt--sts--instance_iam_roles"></a>
### Nested Schema for `sts.instance_iam_roles`

Required:

- `worker_role_arn` (String) Worker/Compute Node Role ARN