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

package main

***REMOVED***
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/terraform-redhat/terraform-provider-ocm/provider"
***REMOVED***

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main(***REMOVED*** {
	tfsdk.Serve(
		context.Background(***REMOVED***,
		provider.New,
		tfsdk.ServeOpts{
			Name: "ocm",
***REMOVED***,
	***REMOVED***
}
