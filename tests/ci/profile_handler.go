package ci

***REMOVED***
***REMOVED***
	"os"
	"strings"

	// . "github.com/onsi/gomega"

	CON "github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/constants"
	EXE "github.com/terraform-redhat/terraform-provider-rhcs/tests/utils/exec"
***REMOVED***

// Profile Provides profile struct for cluster creation be matrix
type Profile struct {
	ClusterName           string `ini:"cluster_name,omitempty" json:"cluster_name,omitempty"`
	ProductID             string `ini:"product_id,omitempty" json:"product_id,omitempty"`
	Version               string `ini:"version,omitempty" json:"version,omitempty"` //Version supports indicated version started with openshift-v or minor-1
	ChannelGroup          string `ini:"channel_group,omitempty" json:"channel_group,omitempty"`
	CloudProvider         string `ini:"cloud_provider,omitempty" json:"cloud_provider,omitempty"`
	Region                string `ini:"region,omitempty" json:"region,omitempty"`
	InstanceType          string `ini:"instance_type,omitempty" json:"instance_type,omitempty"`
	Zones                 string `ini:"zones,omitempty" json:"zones,omitempty"`           // zones should be like a,b,c,d
	StorageLB             bool   `ini:"storage_lb,omitempty" json:"storage_lb,omitempty"` // the unit is GIB, don't support unit set
	Tagging               bool   `ini:"tagging,omitempty" json:"tagging,omitempty"`
	Labeling              bool   `ini:"labeling,omitempty" json:"labeling,omitempty"`
	Etcd                  bool   `ini:"etcd,omitempty" json:"etcd,omitempty"`
	FIPS                  bool   `ini:"fips,omitempty" json:"fips,omitempty"`
	CCS                   bool   `ini:"ccs,omitempty" json:"ccs,omitempty"`
	STS                   bool   `ini:"sts,omitempty" json:"sts,omitempty"`
	Autoscale             bool   `ini:"autoscale,omitempty" json:"autoscale,omitempty"`
	MultiAZ               bool   `ini:"multi_az,omitempty" json:"multi_az,omitempty"`
	BYOVPC                bool   `ini:"byovpc,omitempty" json:"byovpc,omitempty"`
	PrivateLink           bool   `ini:"private_link,omitempty" json:"private_link,omitempty"`
	Private               bool   `ini:"private,omitempty" json:"private,omitempty"`
	BYOK                  bool   `ini:"byok,omitempty" json:"byok,omitempty"`
	ETCDKMS               bool   `ini:"etcd_kms,omitempty" json:"etcd_kms,omitempty"`
	NetWorkingSet         bool   `ini:"networking_set,omitempty" json:"networking_set,omitempty"`
	Proxy                 bool   `ini:"proxy,omitempty" json:"proxy,omitempty"`
	Hypershift            bool   `ini:"hypershift,omitempty" json:"hypershift,omitempty"`
	OIDCConfig            string `ini:"oidc_config,omitempty" json:"oidc_config,omitempty"`
	ProvisionShard        string `ini:"provisionShard,omitempty" json:"provisionShard,omitempty"`
	Ec2MetadataHttpTokens string `ini:"imdsv2,omitempty" json:"imdsv2,omitempty"`
	AuditLogForward       bool   `ini:"auditlog_forward,omitempty" json:"auditlog_forward,omitempty"`
	AdminEnabled          bool   `ini:"admin_enabled,omitempty" json:"admin_enabled,omitempty"`
	ManagedPolicies       bool   `ini:"managed_policies,omitempty" json:"managed_policies,omitempty"`
	VolumeSize            int    `ini:"volume_size,omitempty" json:"volume_size,omitempty"`
	ManifestsDIR          string `ini:"manifests_dir,omitempty" json:"manifests_dir,omitempty"`
}

func PrepareVPC(region string, privateLink bool, multiZone bool, azIDs []string, name ...string***REMOVED*** ([]string, []string, []string***REMOVED*** {

	vpcService := EXE.NewVPCService(***REMOVED***
	vpcArgs := &EXE.VPCArgs{
		AWSRegion: region,
		MultiAZ:   multiZone,
		VPCCIDR:   CON.DefaultVPCCIDR,
	}

	if len(azIDs***REMOVED*** != 0 {
		vpcArgs.AZIDs = azIDs
	}
	if len(name***REMOVED*** == 1 {
		vpcArgs.Name = name[0]
	}
	err := vpcService.Create(vpcArgs***REMOVED***
	if err != nil {
		vpcService.Destroy(***REMOVED***
	}
	privateSubnets, publicSubnets, zones, err := vpcService.Output(***REMOVED***
	// privateSubnets, publicSubnets, zones, err := EXE.CreateAWSVPC(vpcArgs***REMOVED***
	// Expect(err***REMOVED***.ToNot(HaveOccurred(***REMOVED******REMOVED***
	if err != nil {
		vpcService.Destroy(***REMOVED***
		return nil, nil, nil
	}
	return privateSubnets, publicSubnets, zones
}

func PrepareAccountRoles(***REMOVED*** {
}

func PrepareProxy(***REMOVED*** {}

func PrepareKMSKey(***REMOVED*** {}

func PrepareRoute53(***REMOVED*** {}

func GenerateClusterCreationArgsByProfile(profile *Profile***REMOVED*** (clusterArgs *EXE.ClusterCreationArgs, manifestsDir string, err error***REMOVED*** {
	clusterArgs = &EXE.ClusterCreationArgs{
		Token: os.Getenv(CON.TokenENVName***REMOVED***,
	}
	if profile.ClusterName != "" {
		clusterArgs.ClusterName = profile.ClusterName
	} else {
		clusterArgs.ClusterName = "rhcs-tf" // Generate random chars later
	}
	if profile.AdminEnabled {
		// clusterArgs.
	}
	if profile.Region != "" {
		clusterArgs.AWSRegion = profile.Region
	} else {
		clusterArgs.AWSRegion = CON.DefaultAWSRegion
	}

	if profile.STS {
		accService := EXE.NewAccountRoleService(***REMOVED***
		acctPrefix := clusterArgs.ClusterName
		majorVersion := ""
		if profile.Version != "" {
			majorVersion = strings.Join(strings.Split(profile.Version, "."***REMOVED***[0:2], "."***REMOVED***
***REMOVED***
		accountRoleArgs := EXE.AccountRolesArgs{
			AccountRolePrefix: acctPrefix,
			OpenshiftVersion:  majorVersion,
			ChannelGroup:      profile.ChannelGroup,
			Token:             os.Getenv(CON.TokenENVName***REMOVED***,
***REMOVED***
		err = accService.Create(&accountRoleArgs***REMOVED***
		if err != nil {
			defer accService.Destroy(&accountRoleArgs***REMOVED***
			return
***REMOVED***
		clusterArgs.AccountRolePrefix = acctPrefix
		if profile.OIDCConfig != "" {
			clusterArgs.OIDCConfig = profile.OIDCConfig
***REMOVED***
		clusterArgs.OperatorRolePrefix = clusterArgs.ClusterName

	}
	if profile.Region == "" {
		profile.Region = CON.DefaultAWSRegion
	}

	if profile.AuditLogForward {

	}
	if profile.MultiAZ {
		clusterArgs.MultiAZ = true
	}

	if profile.BYOVPC {
		var zones []string
		if profile.Zones != "" {
			zones = strings.Split(profile.Zones, ","***REMOVED***
***REMOVED***

		privateSubnets, publicSubnets, zones := PrepareVPC(profile.Region, profile.PrivateLink, profile.MultiAZ, zones, clusterArgs.ClusterName***REMOVED***
		clusterArgs.AWSAvailabilityZones = zones
		if privateSubnets == nil {
			err = fmt.Errorf("error when creating the vpc, check the previous log. The created resources had been destroyed"***REMOVED***
			return
***REMOVED***
		if profile.PrivateLink {
			clusterArgs.AWSSubnetIDs = privateSubnets
***REMOVED*** else {
			clusterArgs.AWSSubnetIDs = append(privateSubnets, publicSubnets...***REMOVED***
***REMOVED***
	}

	if profile.Version != "" {
		clusterArgs.OpenshiftVersion = profile.Version
	}

	if profile.ChannelGroup != "" {
		clusterArgs.ChannelGroup = profile.ChannelGroup
	}

	if profile.Tagging {

	}
	if profile.Labeling {
		clusterArgs.DefaultMPLabels = map[string]string{
			"test1": "testdata1",
***REMOVED***
	}
	if profile.NetWorkingSet {
		clusterArgs.MachineCIDR = CON.DefaultVPCCIDR
	}

	if profile.Proxy {
	}

	return clusterArgs, profile.ManifestsDIR, err
}

func CreateRHCSClusterByProfile(profile *Profile***REMOVED*** (string, error***REMOVED*** {
	creationArgs, manifests_dir, err := GenerateClusterCreationArgsByProfile(profile***REMOVED***
	if err != nil {
		return "", err
	}
	clusterService := EXE.NewClusterService(manifests_dir***REMOVED***
	err = clusterService.Create(creationArgs***REMOVED***
	if err != nil {
		clusterService.Destroy(creationArgs***REMOVED***
		return "", err
	}
	clusterID, err := clusterService.Output(***REMOVED***
	return clusterID, err
}
