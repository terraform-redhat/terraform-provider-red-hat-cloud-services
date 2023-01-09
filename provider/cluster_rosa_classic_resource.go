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
package provider

***REMOVED***
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
***REMOVED***
	"os"

***REMOVED***
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"

	semver "github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift-online/ocm-sdk-go/logging"
***REMOVED***

const (
	awsCloudProvider = "aws"
	rosaProduct      = "rosa"
	MinVersion       = "4.10"
***REMOVED***

type ClusterRosaClassicResourceType struct {
}

type ClusterRosaClassicResource struct {
	logger     logging.Logger
	collection *cmv1.ClustersClient
}

func (t *ClusterRosaClassicResourceType***REMOVED*** GetSchema(ctx context.Context***REMOVED*** (result tfsdk.Schema,
	diags diag.Diagnostics***REMOVED*** {
	result = tfsdk.Schema{
		Description: "OpenShift managed cluster using rosa sts.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Unique identifier of the cluster.",
				Type:        types.StringType,
				Computed:    true,
	***REMOVED***,
			"external_id": {
				Description: "Unique external identifier of the cluster.",
				Type:        types.StringType,
				Computed:    true,
	***REMOVED***,
			"name": {
				Description: "Name of the cluster.",
				Type:        types.StringType,
				Required:    true,
	***REMOVED***,
			"cloud_region": {
				Description: "Cloud region identifier, for example 'us-east-1'.",
				Type:        types.StringType,
				Required:    true,
	***REMOVED***,
			"sts": {
				Description: "STS Configuration",
				Attributes:  stsResource(***REMOVED***,
				Optional:    true,
	***REMOVED***,
			"multi_az": {
				Description: "Indicates if the cluster should be deployed to " +
					"multiple availability zones. Default value is 'false'.",
				Type:     types.BoolType,
				Optional: true,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(***REMOVED***,
		***REMOVED***,
	***REMOVED***,
			"properties": {
				Description: "User defined properties.",
				Type: types.MapType{
					ElemType: types.StringType,
		***REMOVED***,
				Optional: true,
				Computed: true,
	***REMOVED***,
			"ccs_enabled": {
				Description: "Enables customer cloud subscription.",
				Type:        types.BoolType,
				Computed:    true,
	***REMOVED***,
			"etcd_encryption": {
				Description: "Encrypt etcd data.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
	***REMOVED***,
			"autoscaling_enabled": {
				Description: "Enables autoscaling.",
				Type:        types.BoolType,
				Optional:    true,
	***REMOVED***,
			"min_replicas": {
				Description: "Min replicas.",
				Type:        types.Int64Type,
				Optional:    true,
	***REMOVED***,
			"max_replicas": {
				Description: "Max replicas.",
				Type:        types.Int64Type,
				Optional:    true,
	***REMOVED***,
			"api_url": {
				Description: "URL of the API server.",
				Type:        types.StringType,
				Computed:    true,
	***REMOVED***,
			"console_url": {
				Description: "URL of the console.",
				Type:        types.StringType,
				Computed:    true,
	***REMOVED***,
			"compute_nodes": {
				Description: "Number of compute nodes of the cluster.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
	***REMOVED***,
			"compute_machine_type": {
				Description: "Identifier of the machine type used by the compute nodes, " +
					"for example `r5.xlarge`. Use the `ocm_machine_types` data " +
					"source to find the possible values.",
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(***REMOVED***,
		***REMOVED***,
	***REMOVED***,
			"aws_account_id": {
				Description: "Identifier of the AWS account.",
				Type:        types.StringType,
				Required:    true,
	***REMOVED***,
			"aws_subnet_ids": {
				Description: "aws subnet ids",
				Type: types.ListType{
					ElemType: types.StringType,
		***REMOVED***,
				Optional: true,
	***REMOVED***,
			"aws_private_link": {
				Description: "aws subnet ids",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(***REMOVED***,
		***REMOVED***,
	***REMOVED***,
			"availability_zones": {
				Description: "availability zones",
				Type: types.ListType{
					ElemType: types.StringType,
		***REMOVED***,
				Optional: true,
	***REMOVED***,
			"machine_cidr": {
				Description: "Block of IP addresses for nodes.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
	***REMOVED***,
			"proxy": {
				Description: "proxy",
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"http_proxy": {
						Description: "http proxy",
						Type:        types.StringType,
						Required:    true,
			***REMOVED***,
					"https_proxy": {
						Description: "https proxy",
						Type:        types.StringType,
						Required:    true,
			***REMOVED***,
					"no_proxy": {
						Description: "no proxy",
						Type:        types.StringType,
						Optional:    true,
			***REMOVED***,
		***REMOVED******REMOVED***,
				Optional: true,
	***REMOVED***,
			"service_cidr": {
				Description: "Block of IP addresses for services.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
	***REMOVED***,
			"pod_cidr": {
				Description: "Block of IP addresses for pods.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
	***REMOVED***,
			"host_prefix": {
				Description: "Length of the prefix of the subnet assigned to each node.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
	***REMOVED***,
			"version": {
				Description: "Identifier of the version of OpenShift, for example 'openshift-v4.1.0'.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
	***REMOVED***,
			"state": {
				Description: "State of the cluster.",
				Type:        types.StringType,
				Computed:    true,
	***REMOVED***,
***REMOVED***,
	}
	return
}

func (t *ClusterRosaClassicResourceType***REMOVED*** NewResource(ctx context.Context,
	p tfsdk.Provider***REMOVED*** (result tfsdk.Resource, diags diag.Diagnostics***REMOVED*** {
	// Cast the provider interface to the specific implementation:
	parent := p.(*Provider***REMOVED***

	// Get the collection:
	collection := parent.connection.ClustersMgmt(***REMOVED***.V1(***REMOVED***.Clusters(***REMOVED***

	// Create the resource:
	result = &ClusterRosaClassicResource{
		logger:     parent.logger,
		collection: collection,
	}

	return
}

func createClassicClusterObject(ctx context.Context,
	state *ClusterRosaClassicState, logger logging.Logger, diags diag.Diagnostics***REMOVED*** (*cmv1.Cluster, error***REMOVED*** {
	builder := cmv1.NewCluster(***REMOVED***
	builder.Name(state.Name.Value***REMOVED***
	builder.CloudProvider(cmv1.NewCloudProvider(***REMOVED***.ID(awsCloudProvider***REMOVED******REMOVED***
	builder.Product(cmv1.NewProduct(***REMOVED***.ID(rosaProduct***REMOVED******REMOVED***
	builder.Region(cmv1.NewCloudRegion(***REMOVED***.ID(state.CloudRegion.Value***REMOVED******REMOVED***
	if !state.MultiAZ.Unknown && !state.MultiAZ.Null {
		builder.MultiAZ(state.MultiAZ.Value***REMOVED***
	}
	if !state.Properties.Unknown && !state.Properties.Null {
		properties := map[string]string{}
		for k, v := range state.Properties.Elems {
			properties[k] = v.(types.String***REMOVED***.Value
***REMOVED***
		builder.Properties(properties***REMOVED***
	}

	if !state.EtcdEncryption.Unknown && !state.EtcdEncryption.Null {
		builder.EtcdEncryption(state.EtcdEncryption.Value***REMOVED***
	}

	nodes := cmv1.NewClusterNodes(***REMOVED***
	if !state.ComputeNodes.Unknown && !state.ComputeNodes.Null {
		nodes.Compute(int(state.ComputeNodes.Value***REMOVED******REMOVED***
	}
	if !state.ComputeMachineType.Unknown && !state.ComputeMachineType.Null {
		nodes.ComputeMachineType(
			cmv1.NewMachineType(***REMOVED***.ID(state.ComputeMachineType.Value***REMOVED***,
		***REMOVED***
	}

	if !state.AvailabilityZones.Unknown && !state.AvailabilityZones.Null {
		azs := make([]string, 0***REMOVED***
		for _, e := range state.AvailabilityZones.Elems {
			azs = append(azs, e.(types.String***REMOVED***.Value***REMOVED***
***REMOVED***
		nodes.AvailabilityZones(azs...***REMOVED***
	}

	if !state.AutoScalingEnabled.Unknown && !state.AutoScalingEnabled.Null && state.AutoScalingEnabled.Value {
		autoscaling := cmv1.NewMachinePoolAutoscaling(***REMOVED***
		if !state.MaxReplicas.Unknown && !state.MaxReplicas.Null {
			autoscaling.MaxReplicas(int(state.MaxReplicas.Value***REMOVED******REMOVED***
***REMOVED***
		if !state.MinReplicas.Unknown && !state.MinReplicas.Null {
			autoscaling.MinReplicas(int(state.MinReplicas.Value***REMOVED******REMOVED***
***REMOVED***
		if !autoscaling.Empty(***REMOVED*** {
			nodes.AutoscaleCompute(autoscaling***REMOVED***
***REMOVED***
	}

	if !nodes.Empty(***REMOVED*** {
		builder.Nodes(nodes***REMOVED***
	}

	// ccs should be enabled in ocm rosa clusters
	ccs := cmv1.NewCCS(***REMOVED***
	ccs.Enabled(true***REMOVED***
	builder.CCS(ccs***REMOVED***

	aws := cmv1.NewAWS(***REMOVED***
	if !state.AWSAccountID.Unknown && !state.AWSAccountID.Null {
		aws.AccountID(state.AWSAccountID.Value***REMOVED***
	}
	if !state.AWSPrivateLink.Unknown && !state.AWSPrivateLink.Null {
		aws.PrivateLink((state.AWSPrivateLink.Value***REMOVED******REMOVED***
		api := cmv1.NewClusterAPI(***REMOVED***
		if state.AWSPrivateLink.Value {
			api.Listening(cmv1.ListeningMethodInternal***REMOVED***
***REMOVED***
		builder.API(api***REMOVED***
	}

	sts := cmv1.NewSTS(***REMOVED***
	if state.Sts != nil {
		sts.RoleARN(state.Sts.RoleARN.Value***REMOVED***
		sts.SupportRoleARN(state.Sts.SupportRoleArn.Value***REMOVED***
		instanceIamRoles := cmv1.NewInstanceIAMRoles(***REMOVED***
		instanceIamRoles.MasterRoleARN(state.Sts.InstanceIAMRoles.MasterRoleARN.Value***REMOVED***
		instanceIamRoles.WorkerRoleARN(state.Sts.InstanceIAMRoles.WorkerRoleARN.Value***REMOVED***
		sts.InstanceIAMRoles(instanceIamRoles***REMOVED***

		sts.OperatorRolePrefix(state.Sts.OperatorRolePrefix.Value***REMOVED***
		aws.STS(sts***REMOVED***
	}

	if !state.AWSSubnetIDs.Unknown && !state.AWSSubnetIDs.Null {
		subnetIds := make([]string, 0***REMOVED***
		for _, e := range state.AWSSubnetIDs.Elems {
			subnetIds = append(subnetIds, e.(types.String***REMOVED***.Value***REMOVED***
***REMOVED***
		aws.SubnetIDs(subnetIds...***REMOVED***
	}

	if !aws.Empty(***REMOVED*** {
		builder.AWS(aws***REMOVED***
	}
	network := cmv1.NewNetwork(***REMOVED***
	if !state.MachineCIDR.Unknown && !state.MachineCIDR.Null {
		network.MachineCIDR(state.MachineCIDR.Value***REMOVED***
	}
	if !state.ServiceCIDR.Unknown && !state.ServiceCIDR.Null {
		network.ServiceCIDR(state.ServiceCIDR.Value***REMOVED***
	}
	if !state.PodCIDR.Unknown && !state.PodCIDR.Null {
		network.PodCIDR(state.PodCIDR.Value***REMOVED***
	}
	if !state.HostPrefix.Unknown && !state.HostPrefix.Null {
		network.HostPrefix(int(state.HostPrefix.Value***REMOVED******REMOVED***
	}
	if !network.Empty(***REMOVED*** {
		builder.Network(network***REMOVED***
	}
	if !state.Version.Unknown && !state.Version.Null {
		// TODO: update it to support all cluster versions
		isSupported, err := checkSupportedVersion(state.Version.Value***REMOVED***
		if err != nil {
			logger.Error(ctx, "Error validating required cluster version %s\", err***REMOVED***"***REMOVED***
			errHeadline := "Can't build cluster"
			errDecription := fmt.Sprintf(
				"Can't check if cluster version is supported '%s': %v",
				state.Version.Value, err,
			***REMOVED***
			diags.AddError(
				errHeadline,
				errDecription,
			***REMOVED***
			return nil, errors.New(errHeadline + "\n" + errDecription***REMOVED***
***REMOVED***
		if isSupported {
			builder.Version(cmv1.NewVersion(***REMOVED***.ID(state.Version.Value***REMOVED******REMOVED***
***REMOVED*** else {
			logger.Error(ctx, "Cluster version %s is not supported", state.Version.Value***REMOVED***
			errHeadline := "Can't build cluster"
			errDecription := fmt.Sprintf(
				"Can't check if cluster version is supported '%s': %v",
				state.Version.Value, err,
			***REMOVED***
			diags.AddError(
				errHeadline,
				errDecription,
			***REMOVED***
			return nil, errors.New(errHeadline + "\n" + errDecription***REMOVED***
***REMOVED***
	}

	proxy := cmv1.NewProxy(***REMOVED***
	if state.Proxy != nil {
		proxy.HTTPProxy(state.Proxy.HttpProxy.Value***REMOVED***
		proxy.HTTPSProxy(state.Proxy.HttpsProxy.Value***REMOVED***
		builder.Proxy(proxy***REMOVED***
	}

	object, err := builder.Build(***REMOVED***
	return object, err
}

func (r *ClusterRosaClassicResource***REMOVED*** Create(ctx context.Context,
	request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse***REMOVED*** {
	// Get the plan:
	state := &ClusterRosaClassicState{}
	diags := request.Plan.Get(ctx, state***REMOVED***
	response.Diagnostics.Append(diags...***REMOVED***
	if response.Diagnostics.HasError(***REMOVED*** {
		return
	}

	object, err := createClassicClusterObject(ctx, state, r.logger, diags***REMOVED***
	if err != nil {
		response.Diagnostics.AddError(
			"Can't build cluster",
			fmt.Sprintf(
				"Can't build cluster with name '%s': %v",
				state.Name.Value, err,
			***REMOVED***,
		***REMOVED***
		return
	}
	add, err := r.collection.Add(***REMOVED***.Body(object***REMOVED***.SendContext(ctx***REMOVED***
	if err != nil {
		response.Diagnostics.AddError(
			"Can't create cluster",
			fmt.Sprintf(
				"Can't create cluster with name '%s': %v",
				state.Name.Value, err,
			***REMOVED***,
		***REMOVED***
		return
	}
	object = add.Body(***REMOVED***

	// Save the state:
	populateRosaClassicClusterState(ctx, object, state, r.logger, DefaultHttpClient{}***REMOVED***
	diags = response.State.Set(ctx, state***REMOVED***
	response.Diagnostics.Append(diags...***REMOVED***
}

func (r *ClusterRosaClassicResource***REMOVED*** Read(ctx context.Context, request tfsdk.ReadResourceRequest,
	response *tfsdk.ReadResourceResponse***REMOVED*** {
	// Get the current state:
	state := &ClusterRosaClassicState{}
	diags := request.State.Get(ctx, state***REMOVED***
	response.Diagnostics.Append(diags...***REMOVED***
	if response.Diagnostics.HasError(***REMOVED*** {
		return
	}

	// Find the cluster:
	get, err := r.collection.Cluster(state.ID.Value***REMOVED***.Get(***REMOVED***.SendContext(ctx***REMOVED***
	if err != nil {
		response.Diagnostics.AddError(
			"Can't find cluster",
			fmt.Sprintf(
				"Can't find cluster with identifier '%s': %v",
				state.ID.Value, err,
			***REMOVED***,
		***REMOVED***
		return
	}
	object := get.Body(***REMOVED***

	// Save the state:
	populateRosaClassicClusterState(ctx, object, state, r.logger, DefaultHttpClient{}***REMOVED***
	diags = response.State.Set(ctx, state***REMOVED***
	response.Diagnostics.Append(diags...***REMOVED***
}

func (r *ClusterRosaClassicResource***REMOVED*** Update(ctx context.Context, request tfsdk.UpdateResourceRequest,
	response *tfsdk.UpdateResourceResponse***REMOVED*** {
	var diags diag.Diagnostics

	// Get the state:
	state := &ClusterRosaClassicState{}
	diags = request.State.Get(ctx, state***REMOVED***
	response.Diagnostics.Append(diags...***REMOVED***
	if response.Diagnostics.HasError(***REMOVED*** {
		return
	}

	// Get the plan:
	plan := &ClusterRosaClassicState{}
	diags = request.Plan.Get(ctx, plan***REMOVED***
	response.Diagnostics.Append(diags...***REMOVED***
	if response.Diagnostics.HasError(***REMOVED*** {
		return
	}

	// Send request to update the cluster:
	updateNodes := false
	clusterBuilder := cmv1.NewCluster(***REMOVED***
	clusterNodesBuilder := cmv1.NewClusterNodes(***REMOVED***
	compute, ok := shouldPatchInt(state.ComputeNodes, plan.ComputeNodes***REMOVED***
	if ok {
		clusterNodesBuilder = clusterNodesBuilder.Compute(int(compute***REMOVED******REMOVED***
		updateNodes = true
	}

	if !plan.AutoScalingEnabled.Unknown && !plan.AutoScalingEnabled.Null && plan.AutoScalingEnabled.Value {
		// autoscaling enabled
		autoscaling := cmv1.NewMachinePoolAutoscaling(***REMOVED***

		if !plan.MaxReplicas.Unknown && !plan.MaxReplicas.Null {
			autoscaling = autoscaling.MaxReplicas(int(plan.MaxReplicas.Value***REMOVED******REMOVED***
***REMOVED***
		if !plan.MinReplicas.Unknown && !plan.MinReplicas.Null {
			autoscaling = autoscaling.MinReplicas(int(plan.MinReplicas.Value***REMOVED******REMOVED***
***REMOVED***

		clusterNodesBuilder = clusterNodesBuilder.AutoscaleCompute(autoscaling***REMOVED***
		updateNodes = true

	} else {
		if (!plan.MaxReplicas.Unknown && !plan.MaxReplicas.Null***REMOVED*** || (!plan.MinReplicas.Unknown && !plan.MinReplicas.Null***REMOVED*** {
			response.Diagnostics.AddError(
				"Can't update cluster",
				fmt.Sprintf(
					"Can't update MaxReplica and/or MinReplica of cluster when autoscaling is not enabled",
				***REMOVED***,
			***REMOVED***
			return
***REMOVED***
	}

	if updateNodes {
		clusterBuilder = clusterBuilder.Nodes(clusterNodesBuilder***REMOVED***
	}
	clusterSpec, err := clusterBuilder.Build(***REMOVED***
	if err != nil {
		response.Diagnostics.AddError(
			"Can't build cluster patch",
			fmt.Sprintf(
				"Can't build patch for cluster with identifier '%s': %v",
				state.ID.Value, err,
			***REMOVED***,
		***REMOVED***
		return
	}
	update, err := r.collection.Cluster(state.ID.Value***REMOVED***.Update(***REMOVED***.
		Body(clusterSpec***REMOVED***.
		SendContext(ctx***REMOVED***
	if err != nil {
		response.Diagnostics.AddError(
			"Can't update cluster",
			fmt.Sprintf(
				"Can't update cluster with identifier '%s': %v",
				state.ID.Value, err,
			***REMOVED***,
		***REMOVED***
		return
	}

	// update the autoscaling enabled with the plan value (important for nil and false cases***REMOVED***
	state.AutoScalingEnabled = plan.AutoScalingEnabled
	// update the ComputeNodes with the plan value (important for nil and zero value cases***REMOVED***
	state.ComputeNodes = plan.ComputeNodes

	object := update.Body(***REMOVED***

	// Update the state:
	populateRosaClassicClusterState(ctx, object, state, r.logger, DefaultHttpClient{}***REMOVED***
	diags = response.State.Set(ctx, state***REMOVED***
	response.Diagnostics.Append(diags...***REMOVED***
}

func (r *ClusterRosaClassicResource***REMOVED*** Delete(ctx context.Context, request tfsdk.DeleteResourceRequest,
	response *tfsdk.DeleteResourceResponse***REMOVED*** {
	// Get the state:
	state := &ClusterRosaClassicState{}
	diags := request.State.Get(ctx, state***REMOVED***
	response.Diagnostics.Append(diags...***REMOVED***
	if response.Diagnostics.HasError(***REMOVED*** {
		return
	}

	// Send the request to delete the cluster:
	resource := r.collection.Cluster(state.ID.Value***REMOVED***
	_, err := resource.Delete(***REMOVED***.SendContext(ctx***REMOVED***
	if err != nil {
		response.Diagnostics.AddError(
			"Can't delete cluster",
			fmt.Sprintf(
				"Can't delete cluster with identifier '%s': %v",
				state.ID.Value, err,
			***REMOVED***,
		***REMOVED***
		return
	}

	// Remove the state:
	response.State.RemoveResource(ctx***REMOVED***
}

func (r *ClusterRosaClassicResource***REMOVED*** ImportState(ctx context.Context, request tfsdk.ImportResourceStateRequest,
	response *tfsdk.ImportResourceStateResponse***REMOVED*** {
	// Try to retrieve the object:
	get, err := r.collection.Cluster(request.ID***REMOVED***.Get(***REMOVED***.SendContext(ctx***REMOVED***
	if err != nil {
		response.Diagnostics.AddError(
			"Can't find cluster",
			fmt.Sprintf(
				"Can't find cluster with identifier '%s': %v",
				request.ID, err,
			***REMOVED***,
		***REMOVED***
		return
	}
	object := get.Body(***REMOVED***

	// Save the state:
	state := &ClusterRosaClassicState{}
	err = populateRosaClassicClusterState(ctx, object, state, r.logger, DefaultHttpClient{}***REMOVED***
	if err != nil {
		response.Diagnostics.AddError(
			"Can't populate cluster state",
			fmt.Sprintf(
				"Received error %v", err,
			***REMOVED***,
		***REMOVED***
		return
	}

	diags := response.State.Set(ctx, state***REMOVED***
	response.Diagnostics.Append(diags...***REMOVED***
}

// populateRosaClassicClusterState copies the data from the API object to the Terraform state.
func populateRosaClassicClusterState(ctx context.Context, object *cmv1.Cluster, state *ClusterRosaClassicState, logger logging.Logger, httpClient HttpClient***REMOVED*** error {
	state.ID = types.String{
		Value: object.ID(***REMOVED***,
	}
	state.ExternalID = types.String{
		Value: object.ExternalID(***REMOVED***,
	}
	object.API(***REMOVED***
	state.Name = types.String{
		Value: object.Name(***REMOVED***,
	}
	state.CloudRegion = types.String{
		Value: object.Region(***REMOVED***.ID(***REMOVED***,
	}
	state.MultiAZ = types.Bool{
		Value: object.MultiAZ(***REMOVED***,
	}
	state.Properties = types.Map{
		ElemType: types.StringType,
		Elems:    map[string]attr.Value{},
	}
	for k, v := range object.Properties(***REMOVED*** {
		state.Properties.Elems[k] = types.String{
			Value: v,
***REMOVED***
	}
	state.APIURL = types.String{
		Value: object.API(***REMOVED***.URL(***REMOVED***,
	}
	state.ConsoleURL = types.String{
		Value: object.Console(***REMOVED***.URL(***REMOVED***,
	}
	state.ComputeNodes = types.Int64{
		Value: int64(object.Nodes(***REMOVED***.Compute(***REMOVED******REMOVED***,
	}
	state.ComputeMachineType = types.String{
		Value: object.Nodes(***REMOVED***.ComputeMachineType(***REMOVED***.ID(***REMOVED***,
	}

	autoScaleCompute, ok := object.Nodes(***REMOVED***.GetAutoscaleCompute(***REMOVED***
	if ok {
		var maxReplicas, minReplicas int
		state.AutoScalingEnabled = types.Bool{
			Value: true,
***REMOVED***

		maxReplicas, ok = autoScaleCompute.GetMaxReplicas(***REMOVED***
		if ok {
			state.MaxReplicas = types.Int64{
				Value: int64(maxReplicas***REMOVED***,
	***REMOVED***
***REMOVED***

		minReplicas, ok = autoScaleCompute.GetMinReplicas(***REMOVED***
		if ok {
			state.MinReplicas = types.Int64{
				Value: int64(minReplicas***REMOVED***,
	***REMOVED***
***REMOVED***
	} else {
		// autoscaling not enabled - initialize the MaxReplica and MinReplica
		state.MaxReplicas.Null = true
		state.MinReplicas.Null = true
	}

	azs, ok := object.Nodes(***REMOVED***.GetAvailabilityZones(***REMOVED***
	if ok {
		state.AvailabilityZones.Elems = make([]attr.Value, 0***REMOVED***
		for _, az := range azs {
			state.AvailabilityZones.Elems = append(state.AvailabilityZones.Elems, types.String{
				Value: az,
	***REMOVED******REMOVED***
***REMOVED***
	}

	state.CCSEnabled = types.Bool{
		Value: object.CCS(***REMOVED***.Enabled(***REMOVED***,
	}

	state.EtcdEncryption = types.Bool{
		Value: object.EtcdEncryption(***REMOVED***,
	}

	//The API does not return account id
	awsAccountID, ok := object.AWS(***REMOVED***.GetAccountID(***REMOVED***
	if ok {
		state.AWSAccountID = types.String{
			Value: awsAccountID,
***REMOVED***
	}

	awsPrivateLink, ok := object.AWS(***REMOVED***.GetPrivateLink(***REMOVED***
	if ok {
		state.AWSPrivateLink = types.Bool{
			Value: awsPrivateLink,
***REMOVED***
	} else {
		state.AWSPrivateLink = types.Bool{
			Null: true,
***REMOVED***
	}

	sts, ok := object.AWS(***REMOVED***.GetSTS(***REMOVED***
	if ok {
		if state.Sts == nil {
			state.Sts = &Sts{}
***REMOVED***
		oidc_endpoint_url := sts.OIDCEndpointURL(***REMOVED***
		if strings.HasPrefix(oidc_endpoint_url, "https://"***REMOVED*** {
			oidc_endpoint_url = strings.TrimPrefix(oidc_endpoint_url, "https://"***REMOVED***
***REMOVED***

		state.Sts.OIDCEndpointURL = types.String{
			Value: oidc_endpoint_url,
***REMOVED***
		state.Sts.RoleARN = types.String{
			Value: sts.RoleARN(***REMOVED***,
***REMOVED***
		state.Sts.SupportRoleArn = types.String{
			Value: sts.SupportRoleARN(***REMOVED***,
***REMOVED***
		instanceIAMRoles := sts.InstanceIAMRoles(***REMOVED***
		if instanceIAMRoles != nil {
			state.Sts.InstanceIAMRoles.MasterRoleARN = types.String{
				Value: instanceIAMRoles.MasterRoleARN(***REMOVED***,
	***REMOVED***
			state.Sts.InstanceIAMRoles.WorkerRoleARN = types.String{
				Value: instanceIAMRoles.WorkerRoleARN(***REMOVED***,
	***REMOVED***

***REMOVED***
		// TODO: fix a bug in uhc-cluster-services
		if state.Sts.OperatorRolePrefix.Unknown || state.Sts.OperatorRolePrefix.Null {
			operatorRolePrefix, ok := sts.GetOperatorRolePrefix(***REMOVED***
			if ok {
				state.Sts.OperatorRolePrefix = types.String{
					Value: operatorRolePrefix,
		***REMOVED***
	***REMOVED***
***REMOVED***
		thumbprint, err := getThumbprint(sts.OIDCEndpointURL(***REMOVED***, httpClient***REMOVED***
		if err != nil {
			logger.Error(ctx, "cannot get thumbprint", err***REMOVED***
			state.Sts.Thumbprint = types.String{
				Value: "",
	***REMOVED***
***REMOVED*** else {
			state.Sts.Thumbprint = types.String{
				Value: thumbprint,
	***REMOVED***
***REMOVED***
	}

	subnetIds, ok := object.AWS(***REMOVED***.GetSubnetIDs(***REMOVED***
	if ok {
		state.AWSSubnetIDs.Elems = make([]attr.Value, 0***REMOVED***
		for _, subnetId := range subnetIds {
			state.AWSSubnetIDs.Elems = append(state.AWSSubnetIDs.Elems, types.String{
				Value: subnetId,
	***REMOVED******REMOVED***
***REMOVED***
	}

	proxy, ok := object.GetProxy(***REMOVED***
	if ok {
		state.Proxy.HttpProxy = types.String{
			Value: proxy.HTTPProxy(***REMOVED***,
***REMOVED***
		state.Proxy.HttpsProxy = types.String{
			Value: proxy.HTTPSProxy(***REMOVED***,
***REMOVED***
	}
	machineCIDR, ok := object.Network(***REMOVED***.GetMachineCIDR(***REMOVED***
	if ok {
		state.MachineCIDR = types.String{
			Value: machineCIDR,
***REMOVED***
	} else {
		state.MachineCIDR = types.String{
			Null: true,
***REMOVED***
	}
	serviceCIDR, ok := object.Network(***REMOVED***.GetServiceCIDR(***REMOVED***
	if ok {
		state.ServiceCIDR = types.String{
			Value: serviceCIDR,
***REMOVED***
	} else {
		state.ServiceCIDR = types.String{
			Null: true,
***REMOVED***
	}
	podCIDR, ok := object.Network(***REMOVED***.GetPodCIDR(***REMOVED***
	if ok {
		state.PodCIDR = types.String{
			Value: podCIDR,
***REMOVED***
	} else {
		state.PodCIDR = types.String{
			Null: true,
***REMOVED***
	}
	hostPrefix, ok := object.Network(***REMOVED***.GetHostPrefix(***REMOVED***
	if ok {
		state.HostPrefix = types.Int64{
			Value: int64(hostPrefix***REMOVED***,
***REMOVED***
	} else {
		state.HostPrefix = types.Int64{
			Null: true,
***REMOVED***
	}
	version, ok := object.Version(***REMOVED***.GetID(***REMOVED***
	if ok {
		state.Version = types.String{
			Value: version,
***REMOVED***
	} else {
		state.Version = types.String{
			Null: true,
***REMOVED***
	}
	state.State = types.String{
		Value: string(object.State(***REMOVED******REMOVED***,
	}

	return nil
}

type HttpClient interface {
	Get(url string***REMOVED*** (resp *http.Response, err error***REMOVED***
}

type DefaultHttpClient struct {
}

func (c DefaultHttpClient***REMOVED*** Get(url string***REMOVED*** (resp *http.Response, err error***REMOVED*** {
	return http.Get(url***REMOVED***
}

func getThumbprint(oidcEndpointURL string, httpClient HttpClient***REMOVED*** (thumbprint string, err error***REMOVED*** {
	defer func(***REMOVED*** {
		if panicErr := recover(***REMOVED***; panicErr != nil {
			fmt.Fprintf(os.Stderr, "recovering from: %q\n", panicErr***REMOVED***
			thumbprint = ""
			err = fmt.Errorf("recovering from: %q", panicErr***REMOVED***
***REMOVED***
	}(***REMOVED***

	connect, err := url.ParseRequestURI(oidcEndpointURL***REMOVED***
	if err != nil {
		return "", err
	}

	response, err := httpClient.Get(fmt.Sprintf("https://%s:443", connect.Host***REMOVED******REMOVED***
	if err != nil {
		return "", err
	}

	certChain := response.TLS.PeerCertificates

	// Grab the CA in the chain
	for _, cert := range certChain {
		if cert.IsCA {
			if bytes.Equal(cert.RawIssuer, cert.RawSubject***REMOVED*** {
				hash, err := sha1Hash(cert.Raw***REMOVED***
				if err != nil {
					return "", err
		***REMOVED***
				return hash, nil
	***REMOVED***
***REMOVED***
	}

	// Fall back to using the last certficiate in the chain
	cert := certChain[len(certChain***REMOVED***-1]
	return sha1Hash(cert.Raw***REMOVED***
}

// sha1Hash computes the SHA1 of the byte array and returns the hex encoding as a string.
func sha1Hash(data []byte***REMOVED*** (string, error***REMOVED*** {
	// nolint:gosec
	hasher := sha1.New(***REMOVED***
	_, err := hasher.Write(data***REMOVED***
	if err != nil {
		return "", fmt.Errorf("Couldn't calculate hash:\n %v", err***REMOVED***
	}
	hashed := hasher.Sum(nil***REMOVED***
	return hex.EncodeToString(hashed***REMOVED***, nil
}

func checkSupportedVersion(clusterVersion string***REMOVED*** (bool, error***REMOVED*** {
	rawID := strings.Replace(clusterVersion, "openshift-v", "", 1***REMOVED***
	v1, err := semver.NewVersion(rawID***REMOVED***
	if err != nil {
		return false, err
	}
	v2, err := semver.NewVersion(MinVersion***REMOVED***
	if err != nil {
		return false, err
	}
	//Cluster version is greater than or equal to MinVersion
	return v1.GreaterThanOrEqual(v2***REMOVED***, nil
}
