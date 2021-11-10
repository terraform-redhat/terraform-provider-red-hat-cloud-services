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

package ocm

***REMOVED***
	"context"
	"crypto/x509"
***REMOVED***
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/openshift-online/ocm-sdk-go/logging"
***REMOVED***

// Provider creates the schema for the provider.
func Provider(***REMOVED*** *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			urlKey: {
				Description: "URL of the API server.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     sdk.DefaultURL,
	***REMOVED***,
			tokenURLKey: {
				Description: "OpenID token URL.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     sdk.DefaultTokenURL,
	***REMOVED***,
			userKey: {
				Description: "User name.",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					clientIDKey,
					clientSecretKey,
					tokenKey,
		***REMOVED***,
	***REMOVED***,
			passwordKey: {
				Description: "User password.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ConflictsWith: []string{
					clientIDKey,
					clientSecretKey,
					tokenKey,
		***REMOVED***,
	***REMOVED***,
			tokenKey: {
				Description: "Access or refresh token.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OCM_TOKEN", nil***REMOVED***,
				ConflictsWith: []string{
					clientIDKey,
					clientSecretKey,
					passwordKey,
					userKey,
		***REMOVED***,
	***REMOVED***,
			clientIDKey: {
				Description: "OpenID client identifier.",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					passwordKey,
					tokenKey,
					userKey,
		***REMOVED***,
	***REMOVED***,
			clientSecretKey: {
				Description: "OpenID client secret.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ConflictsWith: []string{
					passwordKey,
					tokenKey,
					userKey,
		***REMOVED***,
	***REMOVED***,
			trustedCAsKey: {
				Description: "PEM encoded certificates of authorities that will " +
					"be trusted. If this isn't explicitly specified then " +
					"the provider will trust the certificate authorities " +
					"trusted by default by the system.",
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
	***REMOVED***,
			insecureKey: {
				Description: "When set to 'true' enables insecure communication " +
					"with the server. This disables verification of TLS " +
					"certificates and host names and it isn't recommended " +
					"for production environments.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
	***REMOVED***,
***REMOVED***,
		ResourcesMap: map[string]*schema.Resource{
			"ocm_cluster":           resourceCluster(***REMOVED***,
			"ocm_identity_provider": resourceIdentityProvider(***REMOVED***,
***REMOVED***,
		DataSourcesMap: map[string]*schema.Resource{
			"ocm_cloud_providers": dataSourceCloudProviders(***REMOVED***,
***REMOVED***,
		ConfigureContextFunc: configure,
	}
}

// configure is the configuration function of the provider. It is responsible for checking the
// connection parameters and creating the connection that will be used by the resources.
func configure(ctx context.Context, data *schema.ResourceData***REMOVED*** (config interface{},
	result diag.Diagnostics***REMOVED*** {
	// Determine the log level used by the SDK from the environment variables used by Terraform:
	logLevel := os.Getenv("TF_LOG_PROVIDER"***REMOVED***
	if logLevel == "" {
		logLevel = os.Getenv("TF_LOG"***REMOVED***
	}
	if logLevel == "" {
		logLevel = logLevelInfo
	}

	// The plugin infrastructure redirects the log package output so that it is sent to the main
	// Terraform process, so if we want to have the logs of the SDK redirected we need to use
	// the log package as well.
	logger, err := logging.NewGoLoggerBuilder(***REMOVED***.
		Debug(logLevel == logLevelDebug***REMOVED***.
		Info(logLevel == logLevelInfo***REMOVED***.
		Warn(logLevel == logLevelWarn***REMOVED***.
		Error(logLevel == logLevelError***REMOVED***.
		Build(***REMOVED***
	if err != nil {
		result = diag.FromErr(err***REMOVED***
		return
	}

	// Create the builder:
	builder := sdk.NewConnectionBuilder(***REMOVED***
	builder.Logger(logger***REMOVED***

	// Copy the settings:
	urlValue, ok := data.GetOk(urlKey***REMOVED***
	if ok {
		builder.URL(urlValue.(string***REMOVED******REMOVED***
	}
	tokenURLValue, ok := data.GetOk(tokenURLKey***REMOVED***
	if ok {
		builder.TokenURL(tokenURLValue.(string***REMOVED******REMOVED***
	}
	userValue, userOk := data.GetOk(userKey***REMOVED***
	passwordValue, passwordOk := data.GetOk(passwordKey***REMOVED***
	if userOk || passwordOk {
		builder.User(userValue.(string***REMOVED***, passwordValue.(string***REMOVED******REMOVED***
	}
	tokenValue, ok := data.GetOk(tokenKey***REMOVED***
	if ok {
		builder.Tokens(tokenValue.(string***REMOVED******REMOVED***
	}
	clientIDValue, clientIDOk := data.GetOk(clientIDKey***REMOVED***
	clientSecretValue, clientSecretOk := data.GetOk(clientSecretKey***REMOVED***
	if clientIDOk || clientSecretOk {
		builder.Client(clientIDValue.(string***REMOVED***, clientSecretValue.(string***REMOVED******REMOVED***
	}
	insecureValue, ok := data.GetOk(insecureKey***REMOVED***
	if ok {
		builder.Insecure(insecureValue.(bool***REMOVED******REMOVED***
	}
	trustedCAs, ok := data.GetOk(trustedCAsKey***REMOVED***
	if ok {
		pool := x509.NewCertPool(***REMOVED***
		if !pool.AppendCertsFromPEM([]byte(trustedCAs.(string***REMOVED******REMOVED******REMOVED*** {
			result = append(result, diag.Diagnostic{
				Severity: diag.Error,
				Detail: fmt.Sprintf(
					"the value of '%s' doesn't contain any certificate",
					trustedCAsKey,
				***REMOVED***,
	***REMOVED******REMOVED***
***REMOVED***
		builder.TrustedCAs(pool***REMOVED***
	}

	// Create the connection:
	connection, err := builder.BuildContext(ctx***REMOVED***
	if err != nil {
		result = diag.FromErr(err***REMOVED***
		return
	}
	config = connection

	return
}

// Log levels:
const (
	logLevelDebug = "DEBUG"
	logLevelInfo  = "INFO"
	logLevelWarn  = "WARN"
	logLevelError = "ERROR"
***REMOVED***
