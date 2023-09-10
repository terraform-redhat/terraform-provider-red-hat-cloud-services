/*
Copyright (c***REMOVED*** 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License"***REMOVED***;
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clusterrosaclassic

***REMOVED***
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
***REMOVED***
***REMOVED***
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/onsi/ginkgo/v2/dsl/core" // nolint
***REMOVED***             // nolint
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	"github.com/terraform-redhat/terraform-provider-rhcs/build"
	"github.com/terraform-redhat/terraform-provider-rhcs/provider/common"
***REMOVED***

type MockHttpClient struct {
	response *http.Response
}

func (c MockHttpClient***REMOVED*** Get(url string***REMOVED*** (resp *http.Response, err error***REMOVED*** {
	return c.response, nil
}

const (
	clusterId         = "1n2j3k4l5m6n7o8p9q0r"
	clusterName       = "my-cluster"
	regionId          = "us-east-1"
	multiAz           = true
	rosaCreatorArn    = "arn:aws:iam::123456789012:dummy/dummy"
	apiUrl            = "https://api.my-cluster.com:6443"
	consoleUrl        = "https://console.my-cluster.com"
	baseDomain        = "alias.p1.openshiftapps.com"
	machineType       = "m5.xlarge"
	availabilityZone1 = "us-east-1a"
	availabilityZone2 = "us-east-1b"
	ccsEnabled        = true
	awsAccountID      = "123456789012"
	privateLink       = false
	oidcEndpointUrl   = "example.com"
	roleArn           = "arn:aws:iam::123456789012:role/role-name"
	httpProxy         = "http://proxy.com"
	httpsProxy        = "https://proxy.com"
	httpTokens        = "required"
***REMOVED***

var (
	mockHttpClient = MockHttpClient{
		response: &http.Response{
			TLS: &tls.ConnectionState{
				PeerCertificates: []*x509.Certificate{
					{
						Raw: []byte("nonce"***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***,
***REMOVED***,
	}
***REMOVED***

func generateBasicRosaClassicClusterJson(***REMOVED*** map[string]interface{} {
	return map[string]interface{}{
		"id":   clusterId,
		"name": clusterName,
		"region": map[string]interface{}{
			"id": regionId,
***REMOVED***,
		"multi_az": multiAz,
		"properties": map[string]interface{}{
			"rosa_creator_arn": rosaCreatorArn,
			"rosa_tf_version":  build.Version,
			"rosa_tf_commit":   build.Commit,
***REMOVED***,
		"api": map[string]interface{}{
			"url": apiUrl,
***REMOVED***,
		"console": map[string]interface{}{
			"url": consoleUrl,
***REMOVED***,
		"dns": map[string]interface{}{
			"base_domain": baseDomain,
***REMOVED***,
		"nodes": map[string]interface{}{
			"compute_machine_type": map[string]interface{}{
				"id": machineType,
	***REMOVED***,
			"availability_zones": []interface{}{
				availabilityZone1,
	***REMOVED***,
***REMOVED***,
		"ccs": map[string]interface{}{
			"enabled": ccsEnabled,
***REMOVED***,
		"aws": map[string]interface{}{
			"account_id":               awsAccountID,
			"private_link":             privateLink,
			"ec2_metadata_http_tokens": httpTokens,
			"sts": map[string]interface{}{
				"oidc_endpoint_url": oidcEndpointUrl,
				"role_arn":          roleArn,
	***REMOVED***,
***REMOVED***,
	}
}

func generateBasicRosaClassicClusterState(***REMOVED*** *ClusterRosaClassicState {
	azs, err := common.StringArrayToList([]string{availabilityZone1}***REMOVED***
	if err != nil {
		return nil
	}
	properties, err := common.ConvertStringMapToMapType(map[string]string{"rosa_creator_arn": rosaCreatorArn}***REMOVED***
	if err != nil {
		return nil
	}
	return &ClusterRosaClassicState{
		Name:              types.StringValue(clusterName***REMOVED***,
		CloudRegion:       types.StringValue(regionId***REMOVED***,
		AWSAccountID:      types.StringValue(awsAccountID***REMOVED***,
		AvailabilityZones: azs,
		Properties:        properties,
		ChannelGroup:      types.StringValue("stable"***REMOVED***,
		Version:           types.StringValue("4.10"***REMOVED***,
		Proxy: &Proxy{
			HttpProxy:  types.StringValue(httpProxy***REMOVED***,
			HttpsProxy: types.StringValue(httpsProxy***REMOVED***,
***REMOVED***,
		Sts:         &Sts{},
		Replicas:    types.Int64Value(2***REMOVED***,
		MinReplicas: types.Int64Null(***REMOVED***,
		MaxReplicas: types.Int64Null(***REMOVED***,
		KMSKeyArn:   types.StringNull(***REMOVED***,
	}
}

func TestResource(t *testing.T***REMOVED*** {
	RegisterFailHandler(Fail***REMOVED***
	RunSpecs(t, "Cluster Rosa Resource Suite"***REMOVED***
}

var _ = Describe("Rosa Classic Sts cluster", func(***REMOVED*** {
	Context("createClassicClusterObject", func(***REMOVED*** {
		It("Creates a cluster with correct field values", func(***REMOVED*** {
			clusterState := generateBasicRosaClassicClusterState(***REMOVED***
			rosaClusterObject, err := createClassicClusterObject(context.Background(***REMOVED***, clusterState, diag.Diagnostics{}***REMOVED***
			Expect(err***REMOVED***.To(BeNil(***REMOVED******REMOVED***

			Expect(rosaClusterObject.Name(***REMOVED******REMOVED***.To(Equal(clusterName***REMOVED******REMOVED***

			id, ok := rosaClusterObject.Region(***REMOVED***.GetID(***REMOVED***
			Expect(ok***REMOVED***.To(BeTrue(***REMOVED******REMOVED***
			Expect(id***REMOVED***.To(Equal(regionId***REMOVED******REMOVED***

			Expect(rosaClusterObject.AWS(***REMOVED***.AccountID(***REMOVED******REMOVED***.To(Equal(awsAccountID***REMOVED******REMOVED***

			availabilityZones := rosaClusterObject.Nodes(***REMOVED***.AvailabilityZones(***REMOVED***
			Expect(availabilityZones***REMOVED***.To(HaveLen(1***REMOVED******REMOVED***
			Expect(availabilityZones[0]***REMOVED***.To(Equal(availabilityZone1***REMOVED******REMOVED***

			Expect(rosaClusterObject.Proxy(***REMOVED***.HTTPProxy(***REMOVED******REMOVED***.To(Equal(httpProxy***REMOVED******REMOVED***
			Expect(rosaClusterObject.Proxy(***REMOVED***.HTTPSProxy(***REMOVED******REMOVED***.To(Equal(httpsProxy***REMOVED******REMOVED***

			arn, ok := rosaClusterObject.Properties(***REMOVED***["rosa_creator_arn"]
			Expect(ok***REMOVED***.To(BeTrue(***REMOVED******REMOVED***
			Expect(arn***REMOVED***.To(Equal(rosaCreatorArn***REMOVED******REMOVED***

			version, ok := rosaClusterObject.Version(***REMOVED***.GetID(***REMOVED***
			Expect(ok***REMOVED***.To(BeTrue(***REMOVED******REMOVED***
			Expect(version***REMOVED***.To(Equal("openshift-v4.10"***REMOVED******REMOVED***
			channel, ok := rosaClusterObject.Version(***REMOVED***.GetChannelGroup(***REMOVED***
			Expect(ok***REMOVED***.To(BeTrue(***REMOVED******REMOVED***
			Expect(channel***REMOVED***.To(Equal("stable"***REMOVED******REMOVED***
***REMOVED******REMOVED***
	}***REMOVED***
	It("Throws an error when version format is invalid", func(***REMOVED*** {
		clusterState := generateBasicRosaClassicClusterState(***REMOVED***
		clusterState.Version = types.StringValue("a.4.1"***REMOVED***
		_, err := createClassicClusterObject(context.Background(***REMOVED***, clusterState, diag.Diagnostics{}***REMOVED***
		Expect(err***REMOVED***.ToNot(BeNil(***REMOVED******REMOVED***
	}***REMOVED***

	It("Throws an error when version is unsupported", func(***REMOVED*** {
		clusterState := generateBasicRosaClassicClusterState(***REMOVED***
		clusterState.Version = types.StringValue("4.1.0"***REMOVED***
		_, err := createClassicClusterObject(context.Background(***REMOVED***, clusterState, diag.Diagnostics{}***REMOVED***
		Expect(err***REMOVED***.ToNot(BeNil(***REMOVED******REMOVED***
	}***REMOVED***

	It("appends the non-default channel name to the requested version", func(***REMOVED*** {
		clusterState := generateBasicRosaClassicClusterState(***REMOVED***
		clusterState.ChannelGroup = types.StringValue("somechannel"***REMOVED***
		rosaClusterObject, err := createClassicClusterObject(context.Background(***REMOVED***, clusterState, diag.Diagnostics{}***REMOVED***
		Expect(err***REMOVED***.To(BeNil(***REMOVED******REMOVED***

		version, ok := rosaClusterObject.Version(***REMOVED***.GetID(***REMOVED***
		Expect(ok***REMOVED***.To(BeTrue(***REMOVED******REMOVED***
		Expect(version***REMOVED***.To(Equal("openshift-v4.10-somechannel"***REMOVED******REMOVED***
		channel, ok := rosaClusterObject.Version(***REMOVED***.GetChannelGroup(***REMOVED***
		Expect(ok***REMOVED***.To(BeTrue(***REMOVED******REMOVED***
		Expect(channel***REMOVED***.To(Equal("somechannel"***REMOVED******REMOVED***
	}***REMOVED***

	Context("populateRosaClassicClusterState", func(***REMOVED*** {
		It("Converts correctly a Cluster object into a ClusterRosaClassicState", func(***REMOVED*** {
			clusterState := &ClusterRosaClassicState{}
			clusterJson := generateBasicRosaClassicClusterJson(***REMOVED***
			clusterJsonString, err := json.Marshal(clusterJson***REMOVED***
			Expect(err***REMOVED***.To(BeNil(***REMOVED******REMOVED***

			clusterObject, err := cmv1.UnmarshalCluster(clusterJsonString***REMOVED***
			Expect(err***REMOVED***.To(BeNil(***REMOVED******REMOVED***

			Expect(populateRosaClassicClusterState(context.Background(***REMOVED***, clusterObject, clusterState, mockHttpClient***REMOVED******REMOVED***.To(Succeed(***REMOVED******REMOVED***

			Expect(clusterState.ID.ValueString(***REMOVED******REMOVED***.To(Equal(clusterId***REMOVED******REMOVED***
			Expect(clusterState.CloudRegion.ValueString(***REMOVED******REMOVED***.To(Equal(regionId***REMOVED******REMOVED***
			Expect(clusterState.MultiAZ.ValueBool(***REMOVED******REMOVED***.To(Equal(multiAz***REMOVED******REMOVED***

			properties, err := common.OptionalMap(context.Background(***REMOVED***, clusterState.Properties***REMOVED***
			Expect(err***REMOVED***.To(BeNil(***REMOVED******REMOVED***
			Expect(properties["rosa_creator_arn"]***REMOVED***.To(Equal(rosaCreatorArn***REMOVED******REMOVED***

			ocmProperties, err := common.OptionalMap(context.Background(***REMOVED***, clusterState.OCMProperties***REMOVED***
			Expect(err***REMOVED***.To(BeNil(***REMOVED******REMOVED***
			Expect(ocmProperties["rosa_tf_version"]***REMOVED***.To(Equal(build.Version***REMOVED******REMOVED***
			Expect(ocmProperties["rosa_tf_commit"]***REMOVED***.To(Equal(build.Commit***REMOVED******REMOVED***

			Expect(clusterState.APIURL.ValueString(***REMOVED******REMOVED***.To(Equal(apiUrl***REMOVED******REMOVED***
			Expect(clusterState.ConsoleURL.ValueString(***REMOVED******REMOVED***.To(Equal(consoleUrl***REMOVED******REMOVED***
			Expect(clusterState.Domain.ValueString(***REMOVED******REMOVED***.To(Equal(fmt.Sprintf("%s.%s", clusterName, baseDomain***REMOVED******REMOVED******REMOVED***
			Expect(clusterState.ComputeMachineType.ValueString(***REMOVED******REMOVED***.To(Equal(machineType***REMOVED******REMOVED***

			Expect(clusterState.AvailabilityZones.Elements(***REMOVED******REMOVED***.To(HaveLen(1***REMOVED******REMOVED***
			azs, err := common.StringListToArray(context.Background(***REMOVED***, clusterState.AvailabilityZones***REMOVED***
			Expect(err***REMOVED***.To(BeNil(***REMOVED******REMOVED***
			Expect(azs[0]***REMOVED***.To(Equal(availabilityZone1***REMOVED******REMOVED***

			Expect(clusterState.CCSEnabled.ValueBool(***REMOVED******REMOVED***.To(Equal(ccsEnabled***REMOVED******REMOVED***
			Expect(clusterState.AWSAccountID.ValueString(***REMOVED******REMOVED***.To(Equal(awsAccountID***REMOVED******REMOVED***
			Expect(clusterState.AWSPrivateLink.ValueBool(***REMOVED******REMOVED***.To(Equal(privateLink***REMOVED******REMOVED***
			Expect(clusterState.Sts.OIDCEndpointURL.ValueString(***REMOVED******REMOVED***.To(Equal(oidcEndpointUrl***REMOVED******REMOVED***
			Expect(clusterState.Sts.RoleARN.ValueString(***REMOVED******REMOVED***.To(Equal(roleArn***REMOVED******REMOVED***
			Expect(clusterState.Ec2MetadataHttpTokens.ValueString(***REMOVED******REMOVED***.To(Equal(httpTokens***REMOVED******REMOVED***
***REMOVED******REMOVED***

		It("Check trimming of oidc url with https perfix", func(***REMOVED*** {
			clusterState := &ClusterRosaClassicState{}
			clusterJson := generateBasicRosaClassicClusterJson(***REMOVED***
			clusterJson["aws"].(map[string]interface{}***REMOVED***["sts"].(map[string]interface{}***REMOVED***["oidc_endpoint_url"] = "https://nonce.com"
			clusterJson["aws"].(map[string]interface{}***REMOVED***["sts"].(map[string]interface{}***REMOVED***["operator_role_prefix"] = "terraform-operator"

			clusterJsonString, err := json.Marshal(clusterJson***REMOVED***
			Expect(err***REMOVED***.To(BeNil(***REMOVED******REMOVED***
			print(string(clusterJsonString***REMOVED******REMOVED***

			clusterObject, err := cmv1.UnmarshalCluster(clusterJsonString***REMOVED***
			Expect(err***REMOVED***.To(BeNil(***REMOVED******REMOVED***

			err = populateRosaClassicClusterState(context.Background(***REMOVED***, clusterObject, clusterState, mockHttpClient***REMOVED***
			Expect(err***REMOVED***.To(BeNil(***REMOVED******REMOVED***
			Expect(clusterState.Sts.OIDCEndpointURL.ValueString(***REMOVED******REMOVED***.To(Equal("nonce.com"***REMOVED******REMOVED***
***REMOVED******REMOVED***

		It("Throws an error when oidc_endpoint_url is an invalid url", func(***REMOVED*** {
			clusterState := &ClusterRosaClassicState{}
			clusterJson := generateBasicRosaClassicClusterJson(***REMOVED***
			clusterJson["aws"].(map[string]interface{}***REMOVED***["sts"].(map[string]interface{}***REMOVED***["oidc_endpoint_url"] = "invalid$url"
			clusterJsonString, err := json.Marshal(clusterJson***REMOVED***
			Expect(err***REMOVED***.To(BeNil(***REMOVED******REMOVED***
			print(string(clusterJsonString***REMOVED******REMOVED***

			clusterObject, err := cmv1.UnmarshalCluster(clusterJsonString***REMOVED***
			Expect(err***REMOVED***.To(BeNil(***REMOVED******REMOVED***

			err = populateRosaClassicClusterState(context.Background(***REMOVED***, clusterObject, clusterState, mockHttpClient***REMOVED***
			Expect(err***REMOVED***.To(BeNil(***REMOVED******REMOVED***
			Expect(clusterState.Sts.Thumbprint.ValueString(***REMOVED******REMOVED***.To(Equal(""***REMOVED******REMOVED***
***REMOVED******REMOVED***
	}***REMOVED***

	Context("http tokens state validation", func(***REMOVED*** {
		It("Fail validation with lower version than allowed", func(***REMOVED*** {
			clusterState := generateBasicRosaClassicClusterState(***REMOVED***
			clusterState.Ec2MetadataHttpTokens = types.StringValue(string(cmv1.Ec2MetadataHttpTokensRequired***REMOVED******REMOVED***
			err := validateHttpTokensVersion(context.Background(***REMOVED***, clusterState, "openshift-v4.10.0"***REMOVED***
			Expect(err***REMOVED***.ToNot(BeNil(***REMOVED******REMOVED***
			Expect(err.Error(***REMOVED******REMOVED***.To(ContainSubstring("is not supported with ec2_metadata_http_tokens"***REMOVED******REMOVED***
***REMOVED******REMOVED***
		It("Pass validation with http_tokens_state and supported version", func(***REMOVED*** {
			clusterState := generateBasicRosaClassicClusterState(***REMOVED***
			err := validateHttpTokensVersion(context.Background(***REMOVED***, clusterState, "openshift-v4.11.0"***REMOVED***
			Expect(err***REMOVED***.To(BeNil(***REMOVED******REMOVED***
***REMOVED******REMOVED***
	}***REMOVED***

}***REMOVED***
