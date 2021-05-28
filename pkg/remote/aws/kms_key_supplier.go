package aws

import (
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/cloudskiff/driftctl/pkg/remote/aws/repository"
	remoteerror "github.com/cloudskiff/driftctl/pkg/remote/error"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"

	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/pkg/resource/aws"

	"github.com/cloudskiff/driftctl/pkg/terraform"
)

type KMSKeySupplier struct {
	reader       terraform.ResourceReader
	deserializer *resource.Deserializer
	client       repository.KMSRepository
	runner       *terraform.ParallelResourceReader
}

func NewKMSKeySupplier(provider *AWSTerraformProvider, deserializer *resource.Deserializer) *KMSKeySupplier {
	return &KMSKeySupplier{
		provider,
		deserializer,
		repository.NewKMSRepository(provider.session),
		terraform.NewParallelResourceReader(provider.Runner().SubRunner()),
	}
}

func (s *KMSKeySupplier) Resources() ([]resource.Resource, error) {
	keys, err := s.client.ListAllKeys()
	if err != nil {
		return nil, remoteerror.NewResourceEnumerationError(err, aws.AwsKmsKeyResourceType)
	}

	for _, key := range keys {
		key := key
		s.runner.Run(func() (cty.Value, error) {
			return s.readKey(key)
		})
	}

	retrieve, err := s.runner.Wait()
	if err != nil {
		return nil, err
	}

	return s.deserializer.Deserialize(aws.AwsKmsKeyResourceType, retrieve)
}

func (s *KMSKeySupplier) readKey(key *kms.KeyListEntry) (cty.Value, error) {
	val, err := s.reader.ReadResource(terraform.ReadResourceArgs{
		ID: *key.KeyId,
		Ty: aws.AwsKmsKeyResourceType,
	})
	if err != nil {
		logrus.Error(err)
		return cty.NilVal, err
	}
	return *val, nil
}
