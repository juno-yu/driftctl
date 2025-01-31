package aws_test

import (
	"testing"

	"github.com/cloudskiff/driftctl/test"
	"github.com/cloudskiff/driftctl/test/acceptance"
)

func TestAcc_AwsS3Bucket_BucketInUsEast1(t *testing.T) {
	acceptance.Run(t, acceptance.AccTestCase{
		TerraformVersion: "0.14.9",
		Paths:            []string{"./testdata/acc/aws_s3_bucket"},
		Args:             []string{"scan", "--filter", "Type=='aws_s3_bucket' || Type=='aws_s3_bucket_policy'"},
		Checks: []acceptance.AccCheck{
			{
				Env: map[string]string{
					"AWS_REGION": "us-east-1",
				},
				Check: func(result *test.ScanResult, stdout string, err error) {
					if err != nil {
						t.Fatal(err)
					}
					result.AssertManagedCount(5)
					// Actually we have a false positive drift on policies due to AWS behavior
					// AWS rework policy and mutate single element array into a string
					// result.AssertDriftCountTotal(0)
				},
			},
		},
	})
}
