package provider

***REMOVED***
	"github.com/hashicorp/terraform-plugin-framework/types"
***REMOVED***

type ClusterWaiterState struct {
	Cluster types.String `tfsdk:"cluster"`
	Ready   types.Bool   `tfsdk:"ready"`
	Timeout types.Int64  `tfsdk:"timeout"`
}
