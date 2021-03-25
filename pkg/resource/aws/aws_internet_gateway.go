// GENERATED, DO NOT EDIT THIS FILE
package aws

import "github.com/zclconf/go-cty/cty"

const AwsInternetGatewayResourceType = "aws_internet_gateway"

type AwsInternetGateway struct {
	Arn     *string           `cty:"arn" computed:"true"`
	Id      string            `cty:"id" computed:"true"`
	OwnerId *string           `cty:"owner_id" computed:"true"`
	Tags    map[string]string `cty:"tags"`
	VpcId   *string           `cty:"vpc_id"`
	CtyVal  *cty.Value        `diff:"-"`
}

func (r *AwsInternetGateway) TerraformId() string {
	return r.Id
}

func (r *AwsInternetGateway) TerraformType() string {
	return AwsInternetGatewayResourceType
}

func (r *AwsInternetGateway) CtyValue() *cty.Value {
	return r.CtyVal
}
