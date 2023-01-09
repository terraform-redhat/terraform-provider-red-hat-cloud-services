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
	"github.com/hashicorp/terraform-plugin-framework/types"
***REMOVED***

type MachinePoolState struct {
	Cluster            types.String `tfsdk:"cluster"`
	ID                 types.String `tfsdk:"id"`
	MachineType        types.String `tfsdk:"machine_type"`
	Name               types.String `tfsdk:"name"`
	Replicas           types.Int64  `tfsdk:"replicas"`
	AutoScalingEnabled types.Bool   `tfsdk:"autoscaling_enabled"`
	MinReplicas        types.Int64  `tfsdk:"min_replicas"`
	MaxReplicas        types.Int64  `tfsdk:"max_replicas"`
}
