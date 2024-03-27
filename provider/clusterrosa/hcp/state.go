package hcp

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-redhat/terraform-provider-rhcs/provider/clusterrosa/sts"
	"github.com/terraform-redhat/terraform-provider-rhcs/provider/proxy"
)

type ClusterRosaHcpState struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	DomainPrefix   types.String `tfsdk:"domain_prefix"`
	ExternalID     types.String `tfsdk:"external_id"`
	Private        types.Bool   `tfsdk:"private"`
	APIURL         types.String `tfsdk:"api_url"`
	ConsoleURL     types.String `tfsdk:"console_url"`
	ChannelGroup   types.String `tfsdk:"channel_group"`
	EtcdEncryption types.Bool   `tfsdk:"etcd_encryption"`
	Properties     types.Map    `tfsdk:"properties"`
	OCMProperties  types.Map    `tfsdk:"ocm_properties"`
	State          types.String `tfsdk:"state"`

	// AWS fields
	AWSAccountID        types.String `tfsdk:"aws_account_id"`
	AWSBillingAccountID types.String `tfsdk:"aws_billing_account_id"`
	AWSSubnetIDs        types.List   `tfsdk:"aws_subnet_ids"`
	Sts                 *sts.HcpSts  `tfsdk:"sts"`
	CloudRegion         types.String `tfsdk:"cloud_region"`
	KMSKeyArn           types.String `tfsdk:"kms_key_arn"`
	EtcdKmsKeyArn       types.String `tfsdk:"etcd_kms_key_arn"`
	Tags                types.Map    `tfsdk:"tags"`

	// Network fields
	Domain      types.String `tfsdk:"domain"`
	PodCIDR     types.String `tfsdk:"pod_cidr"`
	MachineCIDR types.String `tfsdk:"machine_cidr"`
	ServiceCIDR types.String `tfsdk:"service_cidr"`
	HostPrefix  types.Int64  `tfsdk:"host_prefix"`
	Proxy       *proxy.Proxy `tfsdk:"proxy"`

	// Standard machine pools fields
	ComputeMachineType types.String `tfsdk:"compute_machine_type"`
	Replicas           types.Int64  `tfsdk:"replicas"`
	AvailabilityZones  types.List   `tfsdk:"availability_zones"`

	// Version/Upgrade fields
	Version        types.String `tfsdk:"version"`
	CurrentVersion types.String `tfsdk:"current_version"`
	UpgradeAcksFor types.String `tfsdk:"upgrade_acknowledgements_for"`

	// Meta fields - not related to cluster spec
	DisableWaitingInDestroy        types.Bool  `tfsdk:"disable_waiting_in_destroy"`
	DestroyTimeout                 types.Int64 `tfsdk:"destroy_timeout"`
	WaitForCreateComplete          types.Bool  `tfsdk:"wait_for_create_complete"`
	WaitForStdComputeNodesComplete types.Bool  `tfsdk:"wait_for_std_compute_nodes_complete"`
}
