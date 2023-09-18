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
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/terraform-redhat/terraform-provider-rhcs/provider"
***REMOVED***

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

const rhcsProviderAddress = "registry.terraform.io/terraform-redhat/rhcs"

func main(***REMOVED*** {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve"***REMOVED***
	flag.Parse(***REMOVED***

	opts := providerserver.ServeOpts{
		Address: rhcsProviderAddress,
		Debug:   debug,
	}

	if err := providerserver.Serve(context.Background(***REMOVED***, provider.New, opts***REMOVED***; err != nil {
		log.Fatal(err.Error(***REMOVED******REMOVED***
	}
}
