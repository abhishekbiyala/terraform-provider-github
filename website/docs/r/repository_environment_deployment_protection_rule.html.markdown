---
layout: "github"
page_title: "GitHub: github_repository_environment_deployment_protection_rule"
description: |-
  Creates and manages deployment protection rules
---

# github_repository_environment_deployment_protection_rule

This resource allows you to create and manage deployment protection rules.


## Example Usage

```hcl
resource "github_repository_environment" "env" {
  repository  = "my_repo"
  environment = "my_env"
  deployment_branch_policy {
    protected_branches     = false
    custom_branch_policies = true
  }
}

resource "github_repository_environment_deployment_protection_rule" "foo" {
  depends_on = [github_repository_environment.env]

  repository      = "my_repo"
  environment     = "my_env"
  integration_id  = "5"
}
```


## Argument Reference

The following arguments are supported:

* `repository` - (Required) The repository to create the protection rule in.

* `environment` - (Required) The name of the environment to apply the protection rule.

* `integration_id` - (Required) The ID of the app providing this custom deployment protection rule.

## Attributes Reference

The following additional attributes are exported:

* `id` - The ID of the deployment protection rule.

## Import

```
$ terraform import github_repository_environment_deployment_protection_rule.foo repo:env:id
```
