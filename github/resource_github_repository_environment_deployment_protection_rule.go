package github

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/go-github/v57/github"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGithubRepositoryEnvironmentDeploymentProtectionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceGithubRepositoryEnvironmentDeploymentProtectionRuleCreate,
		Read:   resourceGithubRepositoryEnvironmentDeploymentProtectionRuleRead,
		Delete: resourceGithubRepositoryEnvironmentDeploymentProtectionRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the repository. The name is not case sensitive.",
			},
			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the environment.",
			},
			"integration_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the app providing this custom deployment rule.",
			},
		},
	}

}

func resourceGithubRepositoryEnvironmentDeploymentProtectionRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Owner).v3client
	ctx := context.Background()

	owner := meta.(*Owner).name
	repoName := d.Get("repository").(string)
	envName := d.Get("environment").(string)
	integrationId, err := strconv.ParseInt(d.Get("integration_id").(string), 10, 64)
	if err != nil {
		return err
	}
	escapedEnvName := url.PathEscape(envName)

	createData := github.CustomDeploymentProtectionRuleRequest{
		IntegrationID: github.Int64(integrationId),
	}

	resultKey, _, err := client.Repositories.CreateCustomDeploymentProtectionRule(ctx, owner, repoName, escapedEnvName, &createData)
	if err != nil {
		return err
	}

	d.SetId(buildThreePartID(repoName, escapedEnvName, strconv.FormatInt(resultKey.GetID(), 10)))
	return resourceGithubRepositoryEnvironmentDeploymentProtectionRuleRead(d, meta)
}

func resourceGithubRepositoryEnvironmentDeploymentProtectionRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Owner).v3client
	ctx := context.WithValue(context.Background(), ctxId, d.Id())

	owner := meta.(*Owner).name
	repoName, envName, protectionRuleIDString, err := parseThreePartID(d.Id(), "repository", "environment", "protectionRuleID")
	if err != nil {
		return err
	}

	protectionRuleID, err := strconv.ParseInt(protectionRuleIDString, 10, 64)
	if err != nil {
		return err
	}

	protectionRule, _, err := client.Repositories.GetCustomDeploymentProtectionRule(ctx, owner, repoName, envName, protectionRuleID)
	if err != nil {
		if ghErr, ok := err.(*github.ErrorResponse); ok {
			if ghErr.Response.StatusCode == http.StatusNotModified {
				return nil
			}
			if ghErr.Response.StatusCode == http.StatusNotFound {
				log.Printf("[INFO] Removing custom deployment rule for %s/%s/%s from state because it no longer exists in GitHub",
					owner, repoName, envName)
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("repository", repoName)
	d.Set("environment", envName)
	d.Set("integration_id", protectionRule.App.ID)
	return nil
}

func resourceGithubRepositoryEnvironmentDeploymentProtectionRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Owner).v3client
	ctx := context.Background()

	owner := meta.(*Owner).name
	repoName, envName, protectionRuleIdString, err := parseThreePartID(d.Id(), "repository", "environment", "protectionRuleId")
	if err != nil {
		return err
	}

	protectionRuleId, err := strconv.ParseInt(protectionRuleIdString, 10, 64)
	if err != nil {
		return err
	}

	_, err = client.Repositories.DisableCustomDeploymentProtectionRule(ctx, owner, repoName, envName, protectionRuleId)
	if err != nil {
		return err
	}

	return nil
}
