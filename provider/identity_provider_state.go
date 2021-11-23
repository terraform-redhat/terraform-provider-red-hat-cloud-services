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

type IdentityProviderState struct {
	Cluster  types.String              `tfsdk:"cluster"`
	ID       types.String              `tfsdk:"id"`
	Name     types.String              `tfsdk:"name"`
	HTPasswd *HTPasswdIdentityProvider `tfsdk:"htpasswd"`
	LDAP     *LDAPIdentityProvider     `tfsdk:"ldap"`
}

type HTPasswdIdentityProvider struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type LDAPIdentityProvider struct {
	BindDN       types.String                    `tfsdk:"bind_dn"`
	BindPassword types.String                    `tfsdk:"bind_password"`
	CA           types.String                    `tfsdk:"ca"`
	Insecure     types.Bool                      `tfsdk:"insecure"`
	URL          types.String                    `tfsdk:"url"`
	Attributes   *LDAPIdentityProviderAttributes `tfsdk:"attributes"`
}

type LDAPIdentityProviderAttributes struct {
	EMail             []string `tfsdk:"email"`
	ID                []string `tfsdk:"id"`
	Name              []string `tfsdk:"name"`
	PreferredUsername []string `tfsdk:"preferred_username"`
}
