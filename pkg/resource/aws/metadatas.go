package aws

import "github.com/cloudskiff/driftctl/pkg/resource"

func InitResourcesMetadata(resourceSchemaRepository resource.SchemaRepositoryInterface) {
	initAwsCloudfrontDistributionMetaData(resourceSchemaRepository)
}
