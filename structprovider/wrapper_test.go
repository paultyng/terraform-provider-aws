package structprovider

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

func TestStructResourceSchema(t *testing.T) {
	type TestResource struct {
		Name         string `tf:",optional,computed"`
		Region       string `tf:",computed"`
		Bucket       string `tf:",required,forcenew"`
		ForceDestroy bool   `tf:",optional"`
	}

	sr := newStructResource(reflect.TypeOf((*TestResource)(nil)).Elem())
	actual := sr.Schema
	expected := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},

		"region": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"bucket": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"force_destroy": {
			Type:     schema.TypeBool,
			Optional: true,
			//Default:  false,
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fail()
	}
}
