// GENERATED, DO NOT EDIT THIS FILE
package aws

import "github.com/zclconf/go-cty/cty"

const AwsDefaultRouteTableResourceType = "aws_default_route_table"

type AwsDefaultRouteTable struct {
	DefaultRouteTableId *string   `cty:"default_route_table_id"`
	Id                  string    `cty:"id" computed:"true"`
	OwnerId             *string   `cty:"owner_id" computed:"true"`
	PropagatingVgws     *[]string `cty:"propagating_vgws"` // Could be null in state
	Route               *[]struct {
		CidrBlock              *string `cty:"cidr_block"`
		EgressOnlyGatewayId    *string `cty:"egress_only_gateway_id"`
		GatewayId              *string `cty:"gateway_id"`
		InstanceId             *string `cty:"instance_id"`
		Ipv6CidrBlock          *string `cty:"ipv6_cidr_block"`
		NatGatewayId           *string `cty:"nat_gateway_id"`
		NetworkInterfaceId     *string `cty:"network_interface_id"`
		TransitGatewayId       *string `cty:"transit_gateway_id"`
		VpcEndpointId          *string `cty:"vpc_endpoint_id"`
		VpcPeeringConnectionId *string `cty:"vpc_peering_connection_id"`
	} `cty:"route" computed:"true"`
	Tags   map[string]string `cty:"tags"`
	VpcId  *string           `cty:"vpc_id" computed:"true"`
	CtyVal *cty.Value        `diff:"-"`
}

func (r *AwsDefaultRouteTable) TerraformId() string {
	return r.Id
}

func (r *AwsDefaultRouteTable) TerraformType() string {
	return AwsDefaultRouteTableResourceType
}

func (r *AwsDefaultRouteTable) CtyValue() *cty.Value {
	return r.CtyVal
}
