resource "aws_iam_role" "image_registry_role" {
  name = "${var.operator_role_prefix}-openshift-image-registry-installer-cloud-cr"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRoleWithWebIdentity"
        Effect = "Allow"
        Condition = {
            StringEquals = {
                "${var.rh_oidc_provider_url}:sub" = [
                  "system:serviceaccount:openshift-image-registry:cluster-image-registry-operator",
                  "system:serviceaccount:openshift-image-registry:registry"
                  ]
            }
        }
        Principal = {
          Federated = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:oidc-provider/${var.rh_oidc_provider_url}"
        }
      },
    ]
  }***REMOVED***

  tags = {
    red-hat-managed = true
    rosa_cluster_id = var.cluster_id
    operator_namespace = "openshift-image-registry"
    operator_name = "installer-cloud-credentials"
  }
}

resource "aws_iam_role_policy_attachment" "image_registry_role_policy_attachment" {
  role = aws_iam_role.image_registry_role.name
  policy_arn = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/${var.account_role_prefix}-openshift-image-registry-installer-cloud-creden"
}