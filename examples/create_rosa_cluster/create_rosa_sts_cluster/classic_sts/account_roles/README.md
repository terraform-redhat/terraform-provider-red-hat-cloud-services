# Account-Wide IAM Roles

Prior to creating a ROSA STS cluster, you must create the required account-wide roles and policies. For more information, see [Account-wide IAM role and policy reference](https://access.redhat.com/documentation/en-us/red_hat_openshift_service_on_aws/4/html/introduction_to_rosa/rosa-sts-about-iam-resources#rosa-sts-account-wide-roles-and-policies_rosa-sts-about-iam-resources***REMOVED*** in the Red Hat Customer Portal.

## Prerequisites

* You have an AWS account.
* You have installed the latest version ROSA CLI (`rosa`***REMOVED***
* You must have an offline OCM token. This token can be generated in the [Red Hat Hybrid Cloud Console](https://console.redhat.com/openshift/token***REMOVED***.
* You have installed Terraform. See the [Hashicorp Terraform page](https://developer.hashicorp.com/terraform/downloads***REMOVED*** for the latest version.

## Account-wide IAM role creation

1. To run the `terraform apply` you need to set up some variables. This guide uses environmental variables. For more on Terraform variables, see [Managing Variables](https://developer.hashicorp.com/terraform/enterprise/workspaces/variables/managing-variables***REMOVED*** in the Terraform documentation.

   Run the following commands to export your variables. Provide your values in lieu of the brackets. Note that any values declared in the `variables.tf` are the default values if you do not export a superseding value.
        
    1. Your account role prefix prepends to all of your created roles.  This value cannot end with a hyphen (-***REMOVED***.
        ```
        export TF_VAR_account_role_prefix=<account_role_prefix>
        ```
    2.  This variable should be your full [OCM offline token](https://console.redhat.com/openshift/token***REMOVED*** that you generated in the prerequisites.  
        ```
        export TF_VAR_token=<ocm_offline_token>
        ```
    3.  This value should point to your OpenShift instance.  
        ```
        export TF_VAR_url=<ocm_url>
        ```
    4.  **Optional**: You can set the desired OpenShift version with this variable. The default is 4.13.
        ```    
        export TF_VAR_openshift_version=<choose_openshift_version>
        ```
    5.  **Optional**: If you want to set any specific AWS tags for your account roles, you can use this variable to declare those tags.   
        ```    
        export TF_VAR_tags=<aws_resource_tags> (Optional***REMOVED*** 
        ```   
1. In your local copy of this `account-roles` folder, run the following command:
   ````
   terraform init
   ````
   Running this command accesses all the necessary provider information to apply your Terraform plan.
1. **Optional**: Run the `plan` command to ensure that your Terraform files build correctly without errors. This is not required to apply your Terraform plans.
   ````
   terraform plan -out account-roles.tfplan
   ````
1. Run the apply command to create your account roles. 

   > **Note**: If you did not run the `plan` command, you can simply just `apply` without specifying a file.

    ````
    terraform apply <"account-roles.tfplan">
    ````
1. The Terraform applies the plan and creates the account roles. You will see a prompt to confirm you want to create these resources. Enter `yes` then the process will complete with your resources.

## Verification

1. In your command-line interface, run the following command to verify that the account roles are created:
    ````
    rosa list account-roles
    ````
1. You see your roles when the command finishes. 
    ````
    I: Fetching account roles
    ROLE NAME                           ROLE TYPE      ROLE ARN                                                    OPENSHIFT VERSION  AWS Managed
    ManagedOpenShift-ControlPlane-Role  Control plane  arn:aws:iam::XXXXX:role/ManagedOpenShift-ControlPlane-Role  4.13               No
    ManagedOpenShift-Installer-Role     Installer      arn:aws:iam::XXXXX:role/ManagedOpenShift-Installer-Role     4.13               No
    ManagedOpenShift-Support-Role       Support        arn:aws:iam::XXXXX:role/ManagedOpenShift-Support-Role       4.13               No
    ManagedOpenShift-Worker-Role        Worker         arn:aws:iam::XXXXX:role/ManagedOpenShift-Worker-Role        4.13               No

## Resource clean up

After you are done with the resources you created, you should not delete them manually, but instead, use the `destroy` command. Run the following to delete all of your created resources:
  
    terraform destroy

After the command is complete, your resources are deleted.

> **NOTE**: If you manually delete a resource, you create unresolvable issues within your environment.