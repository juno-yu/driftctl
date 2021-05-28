package aws

import (
	"github.com/cloudskiff/driftctl/pkg/remote/aws/repository"
	remoteerror "github.com/cloudskiff/driftctl/pkg/remote/error"

	"github.com/cloudskiff/driftctl/pkg/resource"
	resourceaws "github.com/cloudskiff/driftctl/pkg/resource/aws"

	"github.com/cloudskiff/driftctl/pkg/terraform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
)

type EC2EbsSnapshotSupplier struct {
	reader       terraform.ResourceReader
	deserializer *resource.Deserializer
	client       repository.EC2Repository
	runner       *terraform.ParallelResourceReader
}

func NewEC2EbsSnapshotSupplier(provider *AWSTerraformProvider, deserializer *resource.Deserializer) *EC2EbsSnapshotSupplier {
	return &EC2EbsSnapshotSupplier{
		provider,
		deserializer,
		repository.NewEC2Repository(provider.session),
		terraform.NewParallelResourceReader(provider.Runner().SubRunner()),
	}
}

func (s *EC2EbsSnapshotSupplier) Resources() ([]resource.Resource, error) {
	snapshots, err := s.client.ListAllSnapshots()
	if err != nil {
		return nil, remoteerror.NewResourceEnumerationError(err, resourceaws.AwsEbsSnapshotResourceType)
	}
	results := make([]cty.Value, 0)
	if len(snapshots) > 0 {
		for _, snapshot := range snapshots {
			snap := *snapshot
			s.runner.Run(func() (cty.Value, error) {
				return s.readEbsSnapshot(snap)
			})
		}
		results, err = s.runner.Wait()
		if err != nil {
			return nil, err
		}
	}
	return s.deserializer.Deserialize(resourceaws.AwsEbsSnapshotResourceType, results)
}

func (s *EC2EbsSnapshotSupplier) readEbsSnapshot(snapshot ec2.Snapshot) (cty.Value, error) {
	id := aws.StringValue(snapshot.SnapshotId)
	resSnapshot, err := s.reader.ReadResource(terraform.ReadResourceArgs{
		Ty: resourceaws.AwsEbsSnapshotResourceType,
		ID: id,
	})
	if err != nil {
		logrus.Warnf("Error reading snapshot %s[%s]: %+v", id, resourceaws.AwsEbsSnapshotResourceType, err)
		return cty.NilVal, err
	}
	return *resSnapshot, nil
}
