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
	"context"
***REMOVED***
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift-online/ocm-sdk-go/logging"
***REMOVED***

type MachinePoolResourceType struct {
}

type MachinePoolResource struct {
	logger     logging.Logger
	collection *cmv1.ClustersClient
}

func (t *MachinePoolResourceType***REMOVED*** GetSchema(ctx context.Context***REMOVED*** (result tfsdk.Schema,
	diags diag.Diagnostics***REMOVED*** {
	result = tfsdk.Schema{
		Description: "Machine pool.",
		Attributes: map[string]tfsdk.Attribute{
			"cluster_id": {
				Description: "Identifier of the cluster.",
				Type:        types.StringType,
				Required:    true,
	***REMOVED***,
			"id": {
				Description: "Unique identifier of the machine pool.",
				Type:        types.StringType,
				Computed:    true,
	***REMOVED***,
			"name": {
				Description: "Name of the machine pool.",
				Type:        types.StringType,
				Required:    true,
	***REMOVED***,
			"machine_type": {
				Description: "Identifier of the machine type used by the nodes, " +
					"for example `r5.xlarge`. Use the `ocm_machine_types` data " +
					"source to find the possible values.",
				Type:     types.StringType,
				Required: true,
	***REMOVED***,
			"replicas": {
				Description: "The number of machines of the pool",
				Type:        types.Int64Type,
				Required:    true,
	***REMOVED***,
***REMOVED***,
	}
	return
}

func (t *MachinePoolResourceType***REMOVED*** NewResource(ctx context.Context,
	p tfsdk.Provider***REMOVED*** (result tfsdk.Resource, diags diag.Diagnostics***REMOVED*** {
	// Cast the provider interface to the specific implementation: use it directly when needed.
	parent := p.(*Provider***REMOVED***

	// Get the collection of clusters:
	collection := parent.connection.ClustersMgmt(***REMOVED***.V1(***REMOVED***.Clusters(***REMOVED***

	// Create the resource:
	result = &MachinePoolResource{
		logger:     parent.logger,
		collection: collection,
	}

	return
}

func (r *MachinePoolResource***REMOVED*** Create(ctx context.Context,
	request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse***REMOVED*** {
	// Get the plan:
	state := &MachinePoolState{}
	diags := request.Plan.Get(ctx, state***REMOVED***
	response.Diagnostics.Append(diags...***REMOVED***
	if response.Diagnostics.HasError(***REMOVED*** {
		return
	}

	// Wait till the cluster is ready:
	resource := r.collection.Cluster(state.ClusterID.Value***REMOVED***
	pollCtx, cancel := context.WithTimeout(ctx, 1*time.Hour***REMOVED***
	defer cancel(***REMOVED***
	_, err := resource.Poll(***REMOVED***.
		Interval(30 * time.Second***REMOVED***.
		Predicate(func(get *cmv1.ClusterGetResponse***REMOVED*** bool {
			return get.Body(***REMOVED***.State(***REMOVED*** == cmv1.ClusterStateReady
***REMOVED******REMOVED***.
		StartContext(pollCtx***REMOVED***
	if err != nil {
		response.Diagnostics.AddError(
			"Can't poll cluster state",
			fmt.Sprintf(
				"Can't poll state of cluster with identifier '%s': %v",
				state.ClusterID.Value, err,
			***REMOVED***,
		***REMOVED***
		return
	}

	// Create the machine pool:
	builder := cmv1.NewMachinePool(***REMOVED***
	builder.ID(state.Name.Value***REMOVED***
	builder.InstanceType(state.MachineType.Value***REMOVED***
	builder.Replicas(int(state.Replicas.Value***REMOVED******REMOVED***
	object, err := builder.Build(***REMOVED***
	if err != nil {
		response.Diagnostics.AddError(
			"Can't build machine pool",
			fmt.Sprintf(
				"Can't build machine pool for cluster '%s': %v",
				state.ClusterID.Value, err,
			***REMOVED***,
		***REMOVED***
		return
	}
	collection := resource.MachinePools(***REMOVED***
	add, err := collection.Add(***REMOVED***.Body(object***REMOVED***.SendContext(ctx***REMOVED***
	if err != nil {
		response.Diagnostics.AddError(
			"Can't create machine pool",
			fmt.Sprintf(
				"Can't create machine pool for cluster '%s': %v",
				state.ClusterID.Value, err,
			***REMOVED***,
		***REMOVED***
		return
	}
	object = add.Body(***REMOVED***

	// Save the state:
	r.populateState(object, state***REMOVED***
	diags = response.State.Set(ctx, state***REMOVED***
	response.Diagnostics.Append(diags...***REMOVED***
}

func (r *MachinePoolResource***REMOVED*** Read(ctx context.Context, request tfsdk.ReadResourceRequest,
	response *tfsdk.ReadResourceResponse***REMOVED*** {
	// Get the current state:
	state := &MachinePoolState{}
	diags := request.State.Get(ctx, state***REMOVED***
	response.Diagnostics.Append(diags...***REMOVED***
	if response.Diagnostics.HasError(***REMOVED*** {
		return
	}

	// Find the identity provider:
	resource := r.collection.Cluster(state.ClusterID.Value***REMOVED***.
		MachinePools(***REMOVED***.
		MachinePool(state.ID.Value***REMOVED***
	get, err := resource.Get(***REMOVED***.SendContext(ctx***REMOVED***
	if err != nil {
		response.Diagnostics.AddError(
			"Can't find machine pool",
			fmt.Sprintf(
				"Can't find machine pool with identifier '%s' for "+
					"cluster '%s': %v",
				state.ID.Value, state.ClusterID.Value, err,
			***REMOVED***,
		***REMOVED***
		return
	}
	object := get.Body(***REMOVED***

	// Save the state:
	r.populateState(object, state***REMOVED***
	diags = response.State.Set(ctx, state***REMOVED***
	response.Diagnostics.Append(diags...***REMOVED***
}

func (r *MachinePoolResource***REMOVED*** Update(ctx context.Context, request tfsdk.UpdateResourceRequest,
	response *tfsdk.UpdateResourceResponse***REMOVED*** {
}

func (r *MachinePoolResource***REMOVED*** Delete(ctx context.Context, request tfsdk.DeleteResourceRequest,
	response *tfsdk.DeleteResourceResponse***REMOVED*** {
	// Get the state:
	state := &MachinePoolState{}
	diags := request.State.Get(ctx, state***REMOVED***
	response.Diagnostics.Append(diags...***REMOVED***
	if response.Diagnostics.HasError(***REMOVED*** {
		return
	}

	// Send the request to delete the machine pool:
	resource := r.collection.Cluster(state.ClusterID.Value***REMOVED***.
		MachinePools(***REMOVED***.
		MachinePool(state.ID.Value***REMOVED***
	_, err := resource.Delete(***REMOVED***.SendContext(ctx***REMOVED***
	if err != nil {
		response.Diagnostics.AddError(
			"Can't delete machine pool",
			fmt.Sprintf(
				"Can't delete machine pool with identifier '%s' for "+
					"cluster '%s': %v",
				state.ID.Value, state.ClusterID.Value, err,
			***REMOVED***,
		***REMOVED***
		return
	}

	// Remove the state:
	response.State.RemoveResource(ctx***REMOVED***
}

func (r *MachinePoolResource***REMOVED*** ImportState(ctx context.Context, request tfsdk.ImportResourceStateRequest,
	response *tfsdk.ImportResourceStateResponse***REMOVED*** {
	tfsdk.ResourceImportStatePassthroughID(
		ctx,
		tftypes.NewAttributePath(***REMOVED***.WithAttributeName("id"***REMOVED***,
		request,
		response,
	***REMOVED***
}

// populateState copies the data from the API object to the Terraform state.
func (r *MachinePoolResource***REMOVED*** populateState(object *cmv1.MachinePool, state *MachinePoolState***REMOVED*** {
	state.ID = types.String{
		Value: object.ID(***REMOVED***,
	}
	state.Name = types.String{
		Value: object.ID(***REMOVED***,
	}
	state.MachineType = types.String{
		Value: object.InstanceType(***REMOVED***,
	}
	state.Replicas = types.Int64{
		Value: int64(object.Replicas(***REMOVED******REMOVED***,
	}
}
