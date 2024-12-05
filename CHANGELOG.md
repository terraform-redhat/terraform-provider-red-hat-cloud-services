## 1.6.7 (5 Dec, 2024)
ENHANCEMENTS:
* Bug fixes:
  * Adjust HCP cluster readiness timeout duration
  * Warn instead of error for rosa_creator_arn property

## 1.6.6 (Oct 31, 2024)
FEATURES:
* Validates machine pool custom disk sizes
* Support for additional security groups on the initial machine pool through the HCP cluster spec
* Includes indicator field to ignore machine pool deletion errors (Only to be used when cluster is defined within the same management file as the machine pools as running 'terraform destroy' over all resources)

ENHANCEMENTS:
* Bug fixes:
  * Add mutex to guarantee only one kubelet config for classic clusters
  * Updated to allow IAM `role` for creator arn
  * Allow initial machine pools to be imported when tags are empty ('{}')
  * Warn instead of error when patching rosa_creator_arn

## 1.6.5 (Oct 14, 2024)
FEATURES:
 * Support for Registry Config
 * Support for Hosted Control Plane Root Disk Size

## 1.6.4 (Sept 17, 2024)
FEATURES:
 * Allow yaml for tuning configs
 * Replace thumbprint fetching to OCM API call
 * Include sensitive to password/secret attributes
ENHANCEMENTS:
* Bug fixes:
  * Check if tags are unknown prior to internally consider null
  * Validate negative durations for autoscaler settings
  * Adjust some warning messages for HCP clusters
  * Adjust regex for rosa creator arn property
* Documentation:
  * Adjust upgrading HCP clusters guide

## 1.6.3 (Aug 15, 2024)
FEATURES:
* Include trusted IP list resource
* Include IMDSv2 option for ROSA Hosted Control Plane cluster resource
* Include ROSA Hosted Control Plane cluster machine pool additional security groups on creation
* Allow longer cluster names
* Allow to update htpasswd idp resource
* Include support of component routes for default ingress on ROSA classic clusters
* Support cluster admin creation for ROSA Hosted Control Plane clusters
* Support of user defined AWS tags on machine pool resource
* Support Kubelet Config resource on ROSA Hosted Control Plane cluster resource

ENHANCEMENTS:
* Bug fixes:
  * Fix to version ID validation
  * Validation for cluster admin user name
  * Validation that cluster id is not empty for subresources
  * Does not allow to edit tags on machine pool of a hosted control plane
  * Allow proxy settings to be reset
  * Allow tuning configs to be reset
  * Disallow to edit STS fields of a cluster
  * Validation that name and subnet ID on machine pools are not empty
  * Validation of AWS account IDs
  * Allows to set etcd encryption as false if no value is supplied to KMS key ARN
  * Forwards AWS billing account ID changes
  * Disables to edit cloud region of cluster resource
  * Disables to edit machine pool instance and subnet ID of a machine pool resource
  * Validates ROSA creator ARN property
  * Disables editing the etcd KMS key ARN
  * Allows assumed-role into ROSA creator ARN property
  * Typo on error message for OIDC
* Documentation:
  * Adjust replicas description
  * Includes more example cases for resources
  * Adjusts ROSA Hosted Control Plane autoscaler docs to mention it is not fully available

## 1.6.2 (May 02, 2024)
FEATURES:
* Adjustment to include RH SRE support role into policies data source

## 1.6.1 (Apr 11, 2024)
FEATURES:
* Adjustments to documentation and HCP control plane readiness timeout

## 1.6.0 (Mar 21, 2024)
FEATURES:
* HCP Resources (OCM-5758)

## 1.5.1 (Feb 29, 2024)
FEATURES:
* Add Machine Pool data source (OCM-5428)
* Add Cluster data source (OCM-3259)
* Add auto generated ClusterAdmin password (OCM-5632, OCM-6092)

ENHANCEMENTS:
* Move shared code to common library (OCM-5798, OCM-5397)
* Add availability_zones and aws_submet_ids lists for MachinePool output (OCM-5414, OCM-5409)
* Bug fixes:
  * Cluster wait timeout should finish with error (OCM-5753, OCM-5399)


## 1.5.0 (Jan 3, 2024)
FEATURES:
* Allow additional security group day1 and day2 (OCM-4716, OCM-4717, OCM-4718)
* Add support for Cluster Autoscaler configuration update (OCM-214)
* Add Root Volume size configuration fro default and none default MachinePool (OCM-171)
* Add Kubelet Config Resource to allow Pids Limit configuration in ROSA Cluster (OCM-4766)
* Allow Ingress day2 configuration in ROSA Cluster (OCM-2643)
* Default MachinePool manipulation after cluster creation can be done only from MachinePool resource (OCM-2648)

ENHANCEMENTS:
* Bug fixes
  * Updating email key validation in ldap IDP (OCM-4473)
  * Fix behavior when user try to import IDP of none existing cluster (OCM-5420)
  * Prevent users from adding empty list to MachinePool Labels (OCM-5416, OCM-5285)
  * Remove any RequiresReplace plan modifier and validate value not changed during update instead (OCM-4928)
* Docs:
  * Add note for default MachinePool in cluster rosa classic - irrelevant after create (OCM-5418)
  * Add Onboarding page for new contributers (OCM-4473)

## 1.4.2 (Nov 28, 2023)
FEATURES:
* Add `infra_id` attribute to `cluster_rosa_classic` resource

ENHANCEMENTS:
* Update the documentation files
* Fix the validators to support unknown values for proxy and private_hosted_zone. For more details please read the [issue-363](https://github.com/terraform-redhat/terraform-provider-rhcs/issues/363)

## 1.4.1 (Nov 20, 2023)
ENHANCEMENTS:
* Fix the error message of `private_hosted_zone` validator in `cluster_rosa_classic` resource  
* Fix a bug in `identity_provider` resource from type openid - `cannot reflect tftypes.List[tftypes.String] into a map, must be a map` 

## 1.4.0 (Oct 19, 2023)
FEATURES:
* Add wait attribute in `cluster_rosa_classic` resource for waiting cluster readiness in the creation flow
* Added new `rhcs_info` data source for OCM account details.

ENHANCEMENTS:
* Docs - adjust descriptions
* Upgrade framework version - update `terraform-plugin-framework` to `v1.3.5` (and not `v1.4.0` due to issues with Terraform CLI version `1.6.0`)
* Provider Attributes changes
  * Remove unused attributes
  * Remove all attributes but Token from the docs - internal attributes
  * Add option for environment variables for all string attributes
  * Add "Authentication and configuration" section in the main index and remove attributes section
* Bug fixes:
  * Add verification that cluster exists for machine pool resource
  * Add validation for IDP htpassward with duplicate username
  * missing availability zones in region validation
  * Allow specifying pool subnet even for 1AZ clusters in machine pool resource

## 1.3.0 (Sep 27, 2023)
FEATURES:
* Private cluster - add new variable "private" indicates if the cluster has private connection
* Add support for creating cluster with pre-defined shared VPC
* Added new "rhcs_dns_domain" resource to allow reserving base domain before cluster creation.
* Support resources reconciliation - if a resource was removed without the use of the Terraform provider, executing "terraform apply" should prompt its recreation.
* Htpasswd identity provider - allow creating with multiple users
* Support MachinePool import into the terraform state

ENHANCEMENTS:
* Bug fixes
  * Adding http tokens default to terraform state in case its not returned
  * Terraform run or import failing after configuring 'additional-trust-bundle-file'
  * Provider produced inconsistent result after apply - additional_trust_bundle
  * Day one MachinePool - fix auto scaling/replicas validations
* Docs:
  * Add s3 missing permission for OIDC provider

## 1.2.4 (Sep 4, 2023)
ENHANCEMENTS:
* Fix for "Provider produced inconsistent result after apply" error when setting proxy.additional_trust_bundle

## 1.2.3 (Aug 24, 2023)
ENHANCEMENTS:
* Fixed a bug in cluster_rosa_resource -Terraform provider panic after adding additional CA bundle to ROSA cluster

## 1.2.2 (Aug 3, 2023)
ENHANCEMENTS:
* Update the documentation files to point the correct links.

## 1.2.1 (Aug 3, 2023)
ENHANCEMENTS:
* Update the documentation files to point the correct links.
* Fix the default value of openshift_version in the examples

## 1.2.0 (Aug 1, 2023)
FEATURES:
* Enable creating cluster admin in cluster create
* Add support for cluster properties update and delete

ENHANCEMENTS:
* Update the documentation files
* identity_provider resource can be imported by terraform import command
* rosa_cluster resource can be imported by terraform import command
* Remove AWS validations from rosa_cluster resource
* Recreate IDP tf resource if it was deleted not from tf
* Recreate rosa_cluster tf resource if it was deleted not from tf
* Recreate MachinePool tf resource if it was deleted not from tf
* Bug fixes:
  * populate rosa_rf_version with cluster properties
  * Cluster properties are now allowed to be added in Day1 and be changed in Day 2
  * TF-provider support creating a single-az machinepool for multi-az cluster
  * Improve error message: replica or autoscaling should be required parameters for creating additional machinepools
  * Validate OCP version in create_account_roles module
  * Support in generated account_role_prefix by terraform provider

## 1.1.0 (Jul 5, 2023)
ENHANCEMENTS:
* Update the documentation files
* Openshift Cluster upgrade improvements
* Add support for cluster properties update and delete

## 1.0.5 (Jun 29, 2023)
FEATURES:
* Add an options to set version in oidc clusters
* Add update/remove taints from machine pool
* Support edit/delete labels of secondary machine pool
* Create new topics on Terraform vars and modifying machine pools.
* Support upgrade cluster 

ENHANCEMENTS:
* Rename all resources prefix to start with `rhcs` (instead of `ocm`)
* Rename "terraform-provider-ocm" to "terraform-provider-rhcs"
* Improve examples
* Remove mandatory openshift-v prefix from create cluster version attribute
* Update the documentation files
* Update CI files
* Fix path also to be used for the operator roles creation
* Fix use_spot_instances attribute usage in machinepool resource

## 1.0.2 (Jun 21, 2023)
FEATURES:
* Added GitHub IDP provider support
* Added Google IDP provider support
* Adding support for http_tokens_state field.
* Added day 2 proxy settings
* Support cluster update/upgrade

ENHANCEMENTS:
* Add and improve documentations and examples
* Improve tests coverage
* Adjust rosa_cluster_resource to support OIDC config ID as an input attribute
* Improve the provider logger

## 1.0.0 (April 4, 2023)
ENHANCEMENTS:
* Bug fixes - Validate that the cluster version is compatible to the selected account roles' version in `cluster_rosa_classic` resource 

## 0.0.3 (Mar 28, 2023)
FEATURES:
* Add `ocm_policies` data source for getting the account role policies and operator role policies from OCM API.

ENHANCEMENTS:
* Add domain attribute to `cluster_rosa_classic` resource
* Bug fixes
  * update the descriptions of several attributes in `cluster_rosa_classic` resource.
  * Stop waiting when the cluster encounters an error state
* Add end-to-end test


## 0.0.2 (Feb 21, 2023)
FEATURES:
* Add `cluster_waiter` resource for addressing the cluster state in cluster creation scenario

ENHANCEMENTS:
* Add BYO OIDC support in `cluster_rosa_classic` resource
* Address cluster state while destroying the cluster_rosa_classic resource
* Add gitlab `identity_provider` resource


## 0.0.1 (Feb 12, 2023)
RESOURCES:
* ocm_cluster
* ocm_cluster_rosa_classic
* ocm_group_membership
* ocm_identity_provider
* ocm_machine_pool

DATA SOURCES: 
* ocm_cloud_providers
* ocm_rosa_operator_roles
* ocm_groups
* ocm_machine_types
* ocm_versions

ENHANCEMENTS:
* Move to a new GitHub organization `terraform-redhat`
* Update the documentation files to be generated by `tfplugindocs` tool
