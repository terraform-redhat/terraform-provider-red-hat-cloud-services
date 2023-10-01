package exec

***REMOVED***
	"context"

	CON "github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/constants"
***REMOVED***
***REMOVED***

type ClusterCreationArgs struct {
	AccountRolePrefix    string            `json:"account_role_prefix,omitempty"`
	OCMENV               string            `json:"rhcs_environment,omitempty"`
	ClusterName          string            `json:"cluster_name,omitempty"`
	OperatorRolePrefix   string            `json:"operator_role_prefix,omitempty"`
	OpenshiftVersion     string            `json:"openshift_version,omitempty"`
	Token                string            `json:"token,omitempty"`
	URL                  string            `json:"url,omitempty"`
	AWSRegion            string            `json:"aws_region,omitempty"`
	AWSAvailabilityZones []string          `json:"aws_availability_zones,omitempty"`
	Replicas             int               `json:"replicas,omitempty"`
	ChannelGroup         string            `json:"channel_group,omitempty"`
	AWSHttpTokensState   string            `json:"aws_http_tokens_state,omitempty"`
	PrivateLink          string            `json:"private_link,omitempty"`
	Private              string            `json:"private,omitempty"`
	AWSSubnetIDs         []string          `json:"aws_subnet_ids,omitempty"`
	ComputeMachineType   string            `json:"compute_machine_type,omitempty"`
	DefaultMPLabels      map[string]string `json:"default_mp_labels,omitempty"`
	DisableSCPChecks     bool              `json:"disable_scp_checks,omitempty"`
	MultiAZ              bool              `json:"multi_az,omitempty"`
	MachineCIDR          string            `json:"machine_cidr,omitempty"`
	OIDCConfigID         string            `json:"oidc_config_id,omitempty"`
}

// Just a placeholder, not research what to output yet.
type ClusterOutout struct {
	ClusterID string `json:"cluster_id,omitempty"`
}

// ******************************************************
// RHCS test cases used
const (

	// MaxExpiration in unit of hour
	MaxExpiration = 168

	// MaxNodeNumber means max node number per cluster/machinepool
	MaxNodeNumber = 180

	// MaxNameLength means cluster name will be trimed when request certificate
	MaxNameLength = 15

	MaxIngressNumber = 2
***REMOVED***

// version channel_groups
const (
	FastChannel      = "fast"
	StableChannel    = "stable"
	NightlyChannel   = "nightly"
	CandidateChannel = "candidate"
***REMOVED***

type ClusterService struct {
	CreationArgs *ClusterCreationArgs
	ManifestDir  string
	Context      context.Context
}

func (creator *ClusterService***REMOVED*** Init(manifestDir string***REMOVED*** error {
	creator.ManifestDir = CON.GrantClusterManifestDir(manifestDir***REMOVED***
	ctx := context.TODO(***REMOVED***
	creator.Context = ctx
	err := runTerraformInit(ctx, creator.ManifestDir***REMOVED***
	if err != nil {
		return err
	}
	return nil

}

func (creator *ClusterService***REMOVED*** Create(createArgs *ClusterCreationArgs, extraArgs ...string***REMOVED*** error {
	args := combineStructArgs(createArgs, extraArgs...***REMOVED***
	_, err := runTerraformApplyWithArgs(creator.Context, creator.ManifestDir, args***REMOVED***
	if err != nil {
		return err
	}
	return nil
}

func (creator *ClusterService***REMOVED*** Output(***REMOVED*** (string, error***REMOVED*** {
	out, err := runTerraformOutput(creator.Context, creator.ManifestDir***REMOVED***
	if err != nil {
		return "", err
	}
	clusterObj := out["cluster_id"]
	clusterID := h.DigString(clusterObj, "value"***REMOVED***
	return clusterID, nil
}

func (creator *ClusterService***REMOVED*** Destroy(createArgs *ClusterCreationArgs, extraArgs ...string***REMOVED*** error {
	args := combineStructArgs(createArgs, extraArgs...***REMOVED***
	err := runTerraformDestroyWithArgs(creator.Context, creator.ManifestDir, args***REMOVED***
	return err
}

func NewClusterService(manifestDir string***REMOVED*** (*ClusterService, error***REMOVED*** {
	sc := &ClusterService{}
	err := sc.Init(manifestDir***REMOVED***
	return sc, err
}
