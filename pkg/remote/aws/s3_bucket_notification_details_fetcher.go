package aws

import (
	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/pkg/resource/aws"
	"github.com/cloudskiff/driftctl/pkg/terraform"
)

type S3BucketNotificationDetailsFetcher struct {
	reader       terraform.ResourceReader
	deserializer *resource.Deserializer
}

func NewS3BucketNotificationDetailsFetcher(provider terraform.ResourceReader, deserializer *resource.Deserializer) *S3BucketNotificationDetailsFetcher {
	return &S3BucketNotificationDetailsFetcher{
		reader:       provider,
		deserializer: deserializer,
	}
}

func (r *S3BucketNotificationDetailsFetcher) ReadDetails(res resource.Resource) (resource.Resource, error) {
	ctyVal, err := r.reader.ReadResource(terraform.ReadResourceArgs{
		Ty: aws.AwsS3BucketNotificationResourceType,
		ID: res.TerraformId(),
		Attributes: map[string]string{
			"alias": *res.Attributes().GetString("region"),
		},
	})
	if err != nil {
		return nil, err
	}
	deserializedRes, err := r.deserializer.DeserializeOne(aws.AwsS3BucketNotificationResourceType, *ctyVal)
	if err != nil {
		return nil, err
	}

	return deserializedRes, nil
}
