---
layout: "github"
page_title: "GitHub: github_repository_environment_deployment_protection_rules"
description: |-
  Get the list of deployment protection rules for a given repo / env.
---

# github_repository_environment_deployment_protection_rules

Use this data source to retrieve deployment protection rules for a repository / environment.

## Example Usage

```hcl
data "github_repository_environment_deployment_protection_rules" "example" {
    repository  = "example-repository"
    environment = "env_name"
}
```

## Argument Reference

* `repository` - (Required) Name of the repository to retrieve the deployment protection rules from.

* `environment` - (Required) Name of the environment to retrieve the deployment protection rules from.

## Attributes Reference

* `deployment_protection_rules` - The list of this repository / environment deployment protection rules. Each element of `deployment_protection_rules` has the following attributes:
    * `integration_id` - Id of the app providing custom protection rule.
    * `slug` - GitHub App slug.
