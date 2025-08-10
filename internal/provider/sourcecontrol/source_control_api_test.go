package sourcecontrol

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tryretool/terraform-provider-retool/internal/sdk/api"
)

func testAppGitHubConfig(t *testing.T, client *api.APIClient) {
	apiRequestConfigPart := api.GitHubConfigAnyOf{
		Type:           "App",
		AppId:          "app_id",
		InstallationId: "installation_id",
		PrivateKey:     "private_key",
	}
	apiRequest := api.SourceControlConfigPutRequest{
		Config: api.SourceControlConfigPutRequestConfig{
			GitHub: &api.GitHub{
				Provider:      "GitHub",
				Org:           "org",
				Repo:          "repo",
				DefaultBranch: "default_branch",
				Config: api.GitHubConfig{
					GitHubConfigAnyOf: &apiRequestConfigPart,
				},
			},
		},
	}

	response, httpResponse, err := client.SourceControlAPI.SourceControlConfigPut(context.Background()).SourceControlConfigPutRequest(apiRequest).Execute()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, response.Data.GitHub.Provider, "GitHub")
	assert.Equal(t, response.Data.GitHub.Config.GitHubConfigAnyOf.Type, "App")
	assert.NotNil(t, httpResponse)
}

func testPersonalGitHubConfig(t *testing.T, client *api.APIClient) {
	apiRequestConfigPart := api.GitHubConfigAnyOf1{
		Type:                "Personal",
		PersonalAccessToken: "personal_access_token",
	}
	apiRequest := api.SourceControlConfigPutRequest{
		Config: api.SourceControlConfigPutRequestConfig{
			GitHub: &api.GitHub{
				Provider:      "GitHub",
				Org:           "org",
				Repo:          "repo",
				DefaultBranch: "default_branch",
				Config: api.GitHubConfig{
					GitHubConfigAnyOf1: &apiRequestConfigPart,
				},
			},
		},
	}

	response, httpResponse, err := client.SourceControlAPI.SourceControlConfigPut(context.Background()).SourceControlConfigPutRequest(apiRequest).Execute()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, response.Data.GitHub.Provider, "GitHub")
	assert.Equal(t, response.Data.GitHub.Config.GitHubConfigAnyOf1.Type, "Personal")
	assert.NotNil(t, httpResponse)
}

func testGitLabConfig(t *testing.T, client *api.APIClient) {
	apiRequest := api.SourceControlConfigPutRequest{
		Config: api.SourceControlConfigPutRequestConfig{
			GitLab: &api.GitLab{
				Provider:      "GitLab",
				Org:           "org",
				Repo:          "repo",
				DefaultBranch: "default_branch",
				Config: api.GitLabConfig{
					ProjectId:          1234,
					Url:                "https://gitlab.com",
					ProjectAccessToken: "project_access_token",
				},
			},
		},
	}

	response, httpResponse, err := client.SourceControlAPI.SourceControlConfigPut(context.Background()).SourceControlConfigPutRequest(apiRequest).Execute()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, response.Data.GitLab.Provider, "GitLab")
	assert.Equal(t, response.Data.GitLab.Config.ProjectId, float32(1234))
	assert.NotNil(t, httpResponse)
}

func testCodeCommitConfig(t *testing.T, client *api.APIClient) {
	apiRequest := api.SourceControlConfigPutRequest{
		Config: api.SourceControlConfigPutRequestConfig{
			AWSCodeCommit: &api.AWSCodeCommit{
				Provider:      "AWS CodeCommit",
				Org:           "org",
				Repo:          "repo",
				DefaultBranch: "default_branch",
				Config: api.AWSCodeCommitConfig{
					Url:             "https://git-codecommit.us-west-2.amazonaws.com",
					Region:          "us-west-2",
					AccessKeyId:     "access_key_id",
					SecretAccessKey: "secret_access_key",
					HttpsUsername:   "https_username",
					HttpsPassword:   "https_password",
				},
			},
		},
	}

	response, httpResponse, err := client.SourceControlAPI.SourceControlConfigPut(context.Background()).SourceControlConfigPutRequest(apiRequest).Execute()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, response.Data.AWSCodeCommit.Provider, "AWS CodeCommit")
	assert.Equal(t, response.Data.AWSCodeCommit.Config.Region, "us-west-2")
	assert.NotNil(t, httpResponse)
}

func testBitbucketAppConfig(t *testing.T, client *api.APIClient) {
	apiRequestConfigPart := api.BitbucketConfigAnyOf{
		Type:        "App",
		Username:    "username",
		AppPassword: "app_password",
	}
	apiRequest := api.SourceControlConfigPutRequest{
		Config: api.SourceControlConfigPutRequestConfig{
			Bitbucket: &api.Bitbucket{
				Provider:      "Bitbucket",
				Org:           "org",
				Repo:          "repo",
				DefaultBranch: "default_branch",
				Config: api.BitbucketConfig{
					BitbucketConfigAnyOf: &apiRequestConfigPart,
				},
			},
		},
	}

	response, httpResponse, err := client.SourceControlAPI.SourceControlConfigPut(context.Background()).SourceControlConfigPutRequest(apiRequest).Execute()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, response.Data.Bitbucket.Provider, "Bitbucket")
	assert.Equal(t, response.Data.Bitbucket.Config.BitbucketConfigAnyOf.Type, "App")
	assert.NotNil(t, httpResponse)
}

func testBitbucketTokenConfig(t *testing.T, client *api.APIClient) {
	apiRequestConfigPart := api.BitbucketConfigAnyOf1{
		Type:  "Token",
		Token: "personal_access_token",
	}
	apiRequest := api.SourceControlConfigPutRequest{
		Config: api.SourceControlConfigPutRequestConfig{
			Bitbucket: &api.Bitbucket{
				Provider:      "Bitbucket",
				Org:           "org",
				Repo:          "repo",
				DefaultBranch: "default_branch",
				Config: api.BitbucketConfig{
					BitbucketConfigAnyOf1: &apiRequestConfigPart,
				},
			},
		},
	}

	response, httpResponse, err := client.SourceControlAPI.SourceControlConfigPut(context.Background()).SourceControlConfigPutRequest(apiRequest).Execute()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, response.Data.Bitbucket.Provider, "Bitbucket")
	assert.Equal(t, response.Data.Bitbucket.Config.BitbucketConfigAnyOf1.Type, "Token")
	assert.NotNil(t, httpResponse)
}

func testAzureReposConfig(t *testing.T, client *api.APIClient) {
	apiRequest := api.SourceControlConfigPutRequest{
		Config: api.SourceControlConfigPutRequestConfig{
			AzureRepos: &api.AzureRepos{
				Provider:      "Azure Repos",
				Org:           "org",
				Repo:          "repo",
				DefaultBranch: "default_branch",
				Config: api.AzureReposConfig{
					Url:                 "https://dev.azure.com/organization",
					Project:             "project",
					User:                "user",
					PersonalAccessToken: "personal_access_token",
					UseBasicAuth:        true,
				},
			},
		},
	}

	response, httpResponse, err := client.SourceControlAPI.SourceControlConfigPut(context.Background()).SourceControlConfigPutRequest(apiRequest).Execute()
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, response.Data.AzureRepos.Provider, "Azure Repos")
	assert.Equal(t, response.Data.AzureRepos.Config.Url, "https://dev.azure.com/organization")
	assert.NotNil(t, httpResponse)
}

func TestSourceControlAPI(t *testing.T) {
	t.Skip("Skipping optional integration test")
	host := os.Getenv("RETOOL_HOST")
	scheme := os.Getenv("RETOOL_SCHEME")
	if scheme == "" {
		scheme = "https"
	}
	accessToken := os.Getenv("RETOOL_ACCESS_TOKEN")

	clientConfig := api.NewConfiguration()
	clientConfig.Host = host
	clientConfig.Scheme = scheme
	clientConfig.Servers = api.ServerConfigurations{
		api.ServerConfiguration{
			URL: "/api/v2",
		},
	}
	clientConfig.AddDefaultHeader("Authorization", "Bearer "+accessToken)
	client := api.NewAPIClient(clientConfig)

	testAppGitHubConfig(t, client)
	testPersonalGitHubConfig(t, client)
	testGitLabConfig(t, client)
	testCodeCommitConfig(t, client)
	testBitbucketAppConfig(t, client)
	testBitbucketTokenConfig(t, client)
	testAzureReposConfig(t, client)
}
