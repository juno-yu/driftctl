package resource

import (
	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/test/schemas"
	"github.com/hashicorp/terraform/providers"
)

type FakeResource struct {
	Id    string
	Type  string
	Attrs *resource.Attributes
}

func (d *FakeResource) Schema() *resource.Schema {
	return nil
}

func (d *FakeResource) TerraformId() string {
	return d.Id
}

func (d *FakeResource) TerraformType() string {
	if d.Type != "" {
		return d.Type
	}
	return "FakeResource"
}

func (d *FakeResource) Attributes() *resource.Attributes {
	return d.Attrs
}

type FakeResourceStringer struct {
	Id    string
	Attrs *resource.Attributes
}

func (d *FakeResourceStringer) Schema() *resource.Schema {
	return nil
}

func (d *FakeResourceStringer) HumanReadableAttributes() map[string]string {
	attrs := make(map[string]string)
	if name := d.Attributes().GetString("name"); name != nil && *name != "" {
		attrs["Name"] = *name
	}
	return attrs
}

func (d *FakeResourceStringer) TerraformId() string {
	return d.Id
}

func (d *FakeResourceStringer) TerraformType() string {
	return "FakeResourceStringer"
}

func (d *FakeResourceStringer) Attributes() *resource.Attributes {
	return d.Attrs
}

func InitFakeSchemaRepository(provider, version string) resource.SchemaRepositoryInterface {
	repo := resource.NewSchemaRepository()
	schema := make(map[string]providers.Schema)
	if provider != "" {
		s, err := schemas.ReadTestSchema(provider, version)
		if err != nil {
			// TODO HANDLER ERROR PROPERLY
			panic(err)
		}
		schema = s
	}
	repo.Init(schema)
	return repo
}
