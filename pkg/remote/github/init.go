package github

import (
	"github.com/cloudskiff/driftctl/pkg/alerter"
	"github.com/cloudskiff/driftctl/pkg/output"
	"github.com/cloudskiff/driftctl/pkg/remote/cache"
	"github.com/cloudskiff/driftctl/pkg/remote/common"
	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/pkg/resource/github"
	"github.com/cloudskiff/driftctl/pkg/terraform"
)

const RemoteGithubTerraform = "github+tf"

/**
 * Initialize remote (configure credentials, launch tf providers and start gRPC clients)
 * Required to use Scanner
 */

func Init(version string, alerter *alerter.Alerter,
	providerLibrary *terraform.ProviderLibrary,
	supplierLibrary *resource.SupplierLibrary,
	remoteLibrary *common.RemoteLibrary,
	progress output.Progress,
	resourceSchemaRepository *resource.SchemaRepository,
	factory resource.ResourceFactory,
	configDir string) error {
	if version == "" {
		version = "4.4.0"
	}

	provider, err := NewGithubTerraformProvider(version, progress, configDir)
	if err != nil {
		return err
	}
	err = provider.Init()
	if err != nil {
		return err
	}

	repositoryCache := cache.New(100)

	repository := NewGithubRepository(provider.GetConfig(), repositoryCache)
	deserializer := resource.NewDeserializer(factory)
	providerLibrary.AddProvider(terraform.GITHUB, provider)

	remoteLibrary.AddEnumerator(NewGithubTeamEnumerator(repository, factory))
	remoteLibrary.AddDetailsFetcher(github.GithubTeamResourceType, common.NewGenericDetailsFetcher(github.GithubTeamResourceType, provider, deserializer))

	remoteLibrary.AddEnumerator(NewGithubRepositoryEnumerator(repository, factory))
	remoteLibrary.AddDetailsFetcher(github.GithubRepositoryResourceType, common.NewGenericDetailsFetcher(github.GithubRepositoryResourceType, provider, deserializer))

	remoteLibrary.AddEnumerator(NewGithubMembershipEnumerator(repository, factory))
	remoteLibrary.AddDetailsFetcher(github.GithubMembershipResourceType, common.NewGenericDetailsFetcher(github.GithubMembershipResourceType, provider, deserializer))

	remoteLibrary.AddEnumerator(NewGithubTeamMembershipEnumerator(repository, factory))
	remoteLibrary.AddDetailsFetcher(github.GithubTeamMembershipResourceType, common.NewGenericDetailsFetcher(github.GithubTeamMembershipResourceType, provider, deserializer))

	remoteLibrary.AddEnumerator(NewGithubBranchProtectionEnumerator(repository, factory))
	remoteLibrary.AddDetailsFetcher(github.GithubBranchProtectionResourceType, common.NewGenericDetailsFetcher(github.GithubBranchProtectionResourceType, provider, deserializer))

	err = resourceSchemaRepository.Init(version, provider.Schema())
	if err != nil {
		return err
	}
	github.InitResourcesMetadata(resourceSchemaRepository)

	return nil
}
