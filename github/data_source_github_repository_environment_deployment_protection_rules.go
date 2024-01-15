package github

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGithubRepositoryEnvironmentDeploymentProtectionRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGithubRepositoryEnvironmentDeploymentProtectionRulesRead,

		Schema: map[string]*schema.Schema{
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The GitHub repository name.",
			},
			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target environment name.",
			},
			"deployment_protection_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"integration_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slug": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGithubRepositoryEnvironmentDeploymentProtectionRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Owner).v3client
	owner := meta.(*Owner).name
	repoName := d.Get("repository").(string)
	environmentName := d.Get("environment").(string)

	rules, _, err := client.Repositories.GetAllDeploymentProtectionRules(context.Background(), owner, repoName, environmentName)
	if err != nil {
		return nil
	}

	results := make([]map[string]interface{}, 0)

	for _, rule := range rules.ProtectionRules {
		ruleMap := make(map[string]interface{})
		ruleMap["integration_id"] = strconv.FormatInt(*rule.App.ID, 10)
		ruleMap["slug"] = rule.App.Slug
		results = append(results, ruleMap)
	}

	d.SetId(repoName + ":" + environmentName)
	d.Set("deployment_protection_rules", results)

	return nil
}
