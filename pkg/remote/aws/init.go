package aws

import (
	"github.com/cloudskiff/driftctl/pkg/alerter"
	"github.com/cloudskiff/driftctl/pkg/output"
	"github.com/cloudskiff/driftctl/pkg/remote/aws/client"
	"github.com/cloudskiff/driftctl/pkg/remote/aws/repository"
	"github.com/cloudskiff/driftctl/pkg/remote/cache"
	"github.com/cloudskiff/driftctl/pkg/remote/common"
	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/pkg/resource/aws"
	"github.com/cloudskiff/driftctl/pkg/terraform"
)

const RemoteAWSTerraform = "aws+tf"

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
		version = "3.19.0"
	}
	provider, err := NewAWSTerraformProvider(version, progress, configDir)
	if err != nil {
		return err
	}
	err = provider.Init()
	if err != nil {
		return err
	}

	repositoryCache := cache.New(100)

	s3Repository := repository.NewS3Repository(client.NewAWSClientFactory(provider.session), repositoryCache)
	ec2repository := repository.NewEC2Repository(provider.session, repositoryCache)
	route53repository := repository.NewRoute53Repository(provider.session, repositoryCache)
	lambdaRepository := repository.NewLambdaRepository(provider.session, repositoryCache)
	rdsRepository := repository.NewRDSRepository(provider.session, repositoryCache)
	sqsRepository := repository.NewSQSRepository(provider.session, repositoryCache)
	snsRepository := repository.NewSNSRepository(provider.session, repositoryCache)
	cloudfrontRepository := repository.NewCloudfrontRepository(provider.session, repositoryCache)
	dynamoDBRepository := repository.NewDynamoDBRepository(provider.session, repositoryCache)
	ecrRepository := repository.NewECRRepository(provider.session, repositoryCache)
	kmsRepository := repository.NewKMSRepository(provider.session, repositoryCache)
	iamRepository := repository.NewIAMRepository(provider.session, repositoryCache)

	deserializer := resource.NewDeserializer(factory)
	providerLibrary.AddProvider(terraform.AWS, provider)

	remoteLibrary.AddEnumerator(NewS3BucketEnumerator(s3Repository, factory, provider.Config))
	remoteLibrary.AddDetailsFetcher(aws.AwsS3BucketResourceType, NewS3BucketDetailsFetcher(provider, deserializer))
	remoteLibrary.AddEnumerator(NewS3BucketInventoryEnumerator(s3Repository, factory, provider.Config))
	remoteLibrary.AddDetailsFetcher(aws.AwsS3BucketInventoryResourceType, NewS3BucketInventoryDetailsFetcher(provider, deserializer))
	remoteLibrary.AddEnumerator(NewS3BucketNotificationEnumerator(s3Repository, factory, provider.Config))
	remoteLibrary.AddDetailsFetcher(aws.AwsS3BucketNotificationResourceType, NewS3BucketNotificationDetailsFetcher(provider, deserializer))
	remoteLibrary.AddEnumerator(NewS3BucketMetricsEnumerator(s3Repository, factory, provider.Config))
	remoteLibrary.AddDetailsFetcher(aws.AwsS3BucketMetricResourceType, NewS3BucketMetricsDetailsFetcher(provider, deserializer))
	remoteLibrary.AddEnumerator(NewS3BucketPolicyEnumerator(s3Repository, factory, provider.Config))
	remoteLibrary.AddDetailsFetcher(aws.AwsS3BucketPolicyResourceType, NewS3BucketPolicyDetailsFetcher(provider, deserializer))
	remoteLibrary.AddEnumerator(NewS3BucketAnalyticEnumerator(s3Repository, factory, provider.Config))
	remoteLibrary.AddDetailsFetcher(aws.AwsS3BucketAnalyticsConfigurationResourceType, NewS3BucketAnalyticDetailsFetcher(provider, deserializer))

	remoteLibrary.AddEnumerator(NewEC2EbsVolumeEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsEbsVolumeResourceType, common.NewGenericDetailsFetcher(aws.AwsEbsVolumeResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewEC2EbsSnapshotEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsEbsSnapshotResourceType, common.NewGenericDetailsFetcher(aws.AwsEbsSnapshotResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewEC2EipEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsEipResourceType, common.NewGenericDetailsFetcher(aws.AwsEipResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewEC2AmiEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsAmiResourceType, common.NewGenericDetailsFetcher(aws.AwsAmiResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewEC2KeyPairEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsKeyPairResourceType, common.NewGenericDetailsFetcher(aws.AwsKeyPairResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewEC2EipAssociationEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsEipAssociationResourceType, common.NewGenericDetailsFetcher(aws.AwsEipAssociationResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewEC2InstanceEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsInstanceResourceType, common.NewGenericDetailsFetcher(aws.AwsInstanceResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewEC2InternetGatewayEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsInternetGatewayResourceType, common.NewGenericDetailsFetcher(aws.AwsInternetGatewayResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewVPCEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsVpcResourceType, common.NewGenericDetailsFetcher(aws.AwsVpcResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewDefaultVPCEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsDefaultVpcResourceType, common.NewGenericDetailsFetcher(aws.AwsDefaultVpcResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewEC2RouteTableEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsRouteTableResourceType, common.NewGenericDetailsFetcher(aws.AwsRouteTableResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewEC2DefaultRouteTableEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsDefaultRouteTableResourceType, NewEC2DefaultRouteTableDetailsFetcher(provider, deserializer))
	remoteLibrary.AddEnumerator(NewEC2RouteTableAssociationEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsRouteTableAssociationResourceType, NewEC2RouteTableAssociationDetailsFetcher(provider, deserializer))
	remoteLibrary.AddEnumerator(NewEC2SubnetEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsSubnetResourceType, common.NewGenericDetailsFetcher(aws.AwsSubnetResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewEC2DefaultSubnetEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsDefaultSubnetResourceType, common.NewGenericDetailsFetcher(aws.AwsDefaultSubnetResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewVPCSecurityGroupEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsSecurityGroupResourceType, common.NewGenericDetailsFetcher(aws.AwsSecurityGroupResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewVPCDefaultSecurityGroupEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsDefaultSecurityGroupResourceType, common.NewGenericDetailsFetcher(aws.AwsDefaultSecurityGroupResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewEC2NatGatewayEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsNatGatewayResourceType, common.NewGenericDetailsFetcher(aws.AwsNatGatewayResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewEC2RouteEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsRouteResourceType, NewEC2RouteDetailsFetcher(provider, deserializer))
	remoteLibrary.AddEnumerator(NewVPCSecurityGroupRuleEnumerator(ec2repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsSecurityGroupRuleResourceType, NewVPCSecurityGroupRuleDetailsFetcher(provider, deserializer))

	remoteLibrary.AddEnumerator(NewKMSKeyEnumerator(kmsRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsKmsKeyResourceType, common.NewGenericDetailsFetcher(aws.AwsKmsKeyResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewKMSAliasEnumerator(kmsRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsKmsAliasResourceType, common.NewGenericDetailsFetcher(aws.AwsKmsAliasResourceType, provider, deserializer))

	remoteLibrary.AddEnumerator(NewRoute53HealthCheckEnumerator(route53repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsRoute53HealthCheckResourceType, common.NewGenericDetailsFetcher(aws.AwsRoute53HealthCheckResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewRoute53ZoneEnumerator(route53repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsRoute53ZoneResourceType, common.NewGenericDetailsFetcher(aws.AwsRoute53ZoneResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewRoute53RecordEnumerator(route53repository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsRoute53RecordResourceType, common.NewGenericDetailsFetcher(aws.AwsRoute53RecordResourceType, provider, deserializer))

	remoteLibrary.AddEnumerator(NewCloudfrontDistributionEnumerator(cloudfrontRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsCloudfrontDistributionResourceType, common.NewGenericDetailsFetcher(aws.AwsCloudfrontDistributionResourceType, provider, deserializer))

	remoteLibrary.AddEnumerator(NewRDSDBInstanceEnumerator(rdsRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsDbInstanceResourceType, common.NewGenericDetailsFetcher(aws.AwsDbInstanceResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewRDSDBSubnetGroupEnumerator(rdsRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsDbSubnetGroupResourceType, common.NewGenericDetailsFetcher(aws.AwsDbSubnetGroupResourceType, provider, deserializer))

	remoteLibrary.AddEnumerator(NewSQSQueueEnumerator(sqsRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsSqsQueueResourceType, common.NewGenericDetailsFetcher(aws.AwsSqsQueueResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewSQSQueuePolicyEnumerator(sqsRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsSqsQueuePolicyResourceType, common.NewGenericDetailsFetcher(aws.AwsSqsQueuePolicyResourceType, provider, deserializer))

	remoteLibrary.AddEnumerator(NewSNSTopicEnumerator(snsRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsSnsTopicResourceType, NewSNSTopicDetailsFetcher(provider, deserializer))
	remoteLibrary.AddEnumerator(NewSNSTopicPolicyEnumerator(snsRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsSnsTopicPolicyResourceType, NewSNSTopicPolicyDetailsFetcher(provider, deserializer))
	remoteLibrary.AddEnumerator(NewSNSTopicSubscriptionEnumerator(snsRepository, factory, alerter))
	remoteLibrary.AddDetailsFetcher(aws.AwsSnsTopicSubscriptionResourceType, NewSNSTopicSubscriptionDetailsFetcher(provider, deserializer))

	remoteLibrary.AddEnumerator(NewDynamoDBTableEnumerator(dynamoDBRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsDynamodbTableResourceType, NewDynamoDBTableDetailsFetcher(provider, deserializer))

	remoteLibrary.AddEnumerator(NewIamPolicyEnumerator(iamRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsIamPolicyResourceType, common.NewGenericDetailsFetcher(aws.AwsIamPolicyResourceType, provider, deserializer))

	remoteLibrary.AddEnumerator(NewLambdaFunctionEnumerator(lambdaRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsLambdaFunctionResourceType, NewLambdaFunctionDetailsFetcher(provider, deserializer))
	remoteLibrary.AddEnumerator(NewLambdaEventSourceMappingEnumerator(lambdaRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsLambdaEventSourceMappingResourceType, common.NewGenericDetailsFetcher(aws.AwsLambdaEventSourceMappingResourceType, provider, deserializer))

	remoteLibrary.AddEnumerator(NewIamUserEnumerator(iamRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsIamUserResourceType, common.NewGenericDetailsFetcher(aws.AwsIamUserResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewIamUserPolicyEnumerator(iamRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsIamUserPolicyResourceType, common.NewGenericDetailsFetcher(aws.AwsIamUserPolicyResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewIamRoleEnumerator(iamRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsIamRoleResourceType, common.NewGenericDetailsFetcher(aws.AwsIamRoleResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewIamAccessKeyEnumerator(iamRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsIamAccessKeyResourceType, NewIamAccessKeyDetailsFetcher(provider, deserializer))
	remoteLibrary.AddEnumerator(NewIamRolePolicyAttachmentEnumerator(iamRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsIamRolePolicyAttachmentResourceType, NewIamRolePolicyAttachmentDetailsFetcher(provider, deserializer))
	remoteLibrary.AddEnumerator(NewIamRolePolicyEnumerator(iamRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsIamRolePolicyResourceType, common.NewGenericDetailsFetcher(aws.AwsIamRolePolicyResourceType, provider, deserializer))
	remoteLibrary.AddEnumerator(NewIamUserPolicyAttachmentEnumerator(iamRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsIamUserPolicyAttachmentResourceType, NewIamUserPolicyAttachmentDetailsFetcher(provider, deserializer))

	remoteLibrary.AddEnumerator(NewECRRepositoryEnumerator(ecrRepository, factory))
	remoteLibrary.AddDetailsFetcher(aws.AwsEcrRepositoryResourceType, common.NewGenericDetailsFetcher(aws.AwsEcrRepositoryResourceType, provider, deserializer))

	err = resourceSchemaRepository.Init(version, provider.Schema())
	if err != nil {
		return err
	}
	aws.InitResourcesMetadata(resourceSchemaRepository)

	return nil
}
