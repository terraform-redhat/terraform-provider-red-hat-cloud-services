/*
Copyright (c) 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
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

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

type RosaOperatorRolesDataSourceType struct {
}

type RosaOperatorRolesDataSource struct {
	awsInquiries *cmv1.AWSInquiriesClient
}

const (
	DefaultAccountRolePrefix = "ManagedOpenShift"
	serviceAccountFmt        = "system:serviceaccount:%s:%s"
)

func (t *RosaOperatorRolesDataSourceType) GetSchema(ctx context.Context) (result tfsdk.Schema,
	diags diag.Diagnostics) {
	result = tfsdk.Schema{
		Description: "List of rosa operator role for a specific cluster.",
		Attributes: map[string]tfsdk.Attribute{
			"operator_role_prefix": {
				Description: "Operator role prefix.",
				Type:        types.StringType,
				Required:    true,
			},
			"account_role_prefix": {
				Description: "Account role prefix.",
				Type:        types.StringType,
				Optional:    true,
			},
			"operator_iam_roles": {
				Description: "Operator IAM Roles.",
				Attributes: tfsdk.ListNestedAttributes(
					t.itemAttributes(),
					tfsdk.ListNestedAttributesOptions{},
				),
				Computed: true,
			},
			"hypershift_enabled": {
				Description: "Enables hypershift.",
				Type:        types.BoolType,
				Optional:    true,
			},
		},
	}
	return
}

func (t *RosaOperatorRolesDataSourceType) itemAttributes() map[string]tfsdk.Attribute {
	return map[string]tfsdk.Attribute{
		"operator_name": {
			Description: "Operator Name",
			Type:        types.StringType,
			Computed:    true,
		},
		"operator_namespace": {
			Description: "Kubernetes Namespace",
			Type:        types.StringType,
			Computed:    true,
		},
		"role_name": {
			Description: "policy name",
			Type:        types.StringType,
			Computed:    true,
		},
		"policy_name": {
			Description: "policy name",
			Type:        types.StringType,
			Computed:    true,
		},
		"service_accounts": {
			Description: "service accounts",
			Type: types.ListType{
				ElemType: types.StringType,
			},
			Computed: true,
		},
	}
}

func (t *RosaOperatorRolesDataSourceType) NewDataSource(ctx context.Context,
	p tfsdk.Provider) (result tfsdk.DataSource, diags diag.Diagnostics) {
	// Cast the provider interface to the specific implementation:
	parent := p.(*Provider)

	// Get the collection of clusters:
	awsInquiries := parent.connection.ClustersMgmt().V1().AWSInquiries()

	// Create the resource:
	result = &RosaOperatorRolesDataSource{
		awsInquiries: awsInquiries,
	}
	return
}

func (t *RosaOperatorRolesDataSource) Read(ctx context.Context, request tfsdk.ReadDataSourceRequest,
	response *tfsdk.ReadDataSourceResponse) {
	// Get the state:
	state := &RosaOperatorRolesState{}
	diags := request.Config.Get(ctx, state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var is_hypershift bool
	if !state.HypershiftEnabled.Unknown && !state.HypershiftEnabled.Null {
		is_hypershift = state.HypershiftEnabled.Value
	}

	stsOperatorRolesList, err := t.awsInquiries.STSCredentialRequests().List().Parameter("is_hypershift", is_hypershift).Send()
	if err != nil {
		description := fmt.Sprintf("Failed to get STS Operator Roles list with error: %v", err)
		tflog.Error(ctx, description)
		response.Diagnostics.AddError(
			description,
			"hint: validate the credetials (token) used to run this provider",
		)
		return
	}

	stsOperatorMap := make(map[string]*cmv1.STSOperator)
	roleServiceAccounts := make([]string, 0)
	stsOperatorRolesList.Items().Each(func(stsCredentialRequest *cmv1.STSCredentialRequest) bool {
		// TODO: check the MinVersion of the operator role
		tflog.Debug(ctx, fmt.Sprintf("Operator name: %s, namespace %s, service account %s",
			stsCredentialRequest.Operator().Name(),
			stsCredentialRequest.Operator().Namespace(),
			stsCredentialRequest.Operator().ServiceAccounts(),
		))
		// The key can't be stsCredentialRequest.Operator().Name() because of constants between
		// the name of `ingress_operator_cloud_credentials` and `cloud_network_config_controller_cloud_credentials`
		// both of them includes the same Name `cloud-credentials` and it cannot be fixed
		roleServiceAccounts = append(roleServiceAccounts, stsCredentialRequest.Operator().ServiceAccounts()[0])
		stsOperatorMap[stsCredentialRequest.Operator().ServiceAccounts()[0]] = stsCredentialRequest.Operator()
		return true
	})

	accountRolePrefix := DefaultAccountRolePrefix
	if !state.AccountRolePrefix.Unknown && !state.AccountRolePrefix.Null && state.AccountRolePrefix.Value != "" {
		accountRolePrefix = state.AccountRolePrefix.Value
	}

	// TODO: use the sts.OperatorRolePrefix() if not empty
	// There is a bug in the return value of sts.OperatorRolePrefix() - it's always empty string
	sort.Strings(roleServiceAccounts)
	for _, key := range roleServiceAccounts {
		r := OperatorIAMRole{
			Name: types.String{
				Value: stsOperatorMap[key].Name(),
			},
			Namespace: types.String{
				Value: stsOperatorMap[key].Namespace(),
			},
			RoleName: types.String{
				Value: getRoleName(state.OperatorRolePrefix.Value, stsOperatorMap[key]),
			},
			PolicyName: types.String{
				Value: getPolicyName(accountRolePrefix, stsOperatorMap[key].Namespace(), stsOperatorMap[key].Name()),
			},
			ServiceAccounts: buildServiceAccountsArray(stsOperatorMap[key].ServiceAccounts(), stsOperatorMap[key].Namespace()),
		}
		state.OperatorIAMRoles = append(state.OperatorIAMRoles, &r)
	}

	// Save the state:
	diags = response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)
}

// TODO: should be in a separate repo
func getRoleName(rolePrefix string, operatorRole *cmv1.STSOperator) string {
	role := fmt.Sprintf("%s-%s-%s", rolePrefix, operatorRole.Namespace(), operatorRole.Name())
	if len(role) > 64 {
		role = role[0:64]
	}
	return role
}

// TODO: should be in a separate repo
func getPolicyName(prefix string, namespace string, name string) string {
	policy := fmt.Sprintf("%s-%s-%s", prefix, namespace, name)
	if len(policy) > 64 {
		policy = policy[0:64]
	}
	return policy
}

func buildServiceAccountsArray(serviceAccountArr []string, operatorNamespace string) types.List {
	serviceAccounts := types.List{
		ElemType: types.StringType,
		Elems:    []attr.Value{},
	}

	for _, v := range serviceAccountArr {
		serviceAccount := fmt.Sprintf(serviceAccountFmt, operatorNamespace, v)
		serviceAccounts.Elems = append(serviceAccounts.Elems, types.String{Value: serviceAccount})
	}

	return serviceAccounts
}
