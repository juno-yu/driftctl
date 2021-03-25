// GENERATED, DO NOT EDIT THIS FILE
package aws

import "github.com/zclconf/go-cty/cty"

const AwsSecurityGroupRuleResourceType = "aws_security_group_rule"

type AwsSecurityGroupRule struct {
	CidrBlocks            *[]string  `cty:"cidr_blocks"`
	Description           *string    `cty:"description"`
	FromPort              *int       `cty:"from_port"`
	Id                    string     `cty:"id" computed:"true"`
	Ipv6CidrBlocks        *[]string  `cty:"ipv6_cidr_blocks"`
	PrefixListIds         *[]string  `cty:"prefix_list_ids"`
	Protocol              *string    `cty:"protocol"`
	SecurityGroupId       *string    `cty:"security_group_id"`
	Self                  *bool      `cty:"self" diff:"-"`
	SourceSecurityGroupId *string    `cty:"source_security_group_id" computed:"true"`
	ToPort                *int       `cty:"to_port"`
	Type                  *string    `cty:"type"`
	CtyVal                *cty.Value `diff:"-"`
}

func (r *AwsSecurityGroupRule) TerraformId() string {
	return r.Id
}

func (r *AwsSecurityGroupRule) TerraformType() string {
	return AwsSecurityGroupRuleResourceType
}

func (r *AwsSecurityGroupRule) CtyValue() *cty.Value {
	return r.CtyVal
}
