package structprovider

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

type TestDataSource struct {
	Name       string `tf:",optional,computed"`
	Region     string `tf:",computed"`
	NameSuffix string `tf:",computed"`
	State      string `tf:",optional,computed"`
}

func (ds *TestDataSource) ID() string {
	return ""
}

func (ds *TestDataSource) Read(meta interface{}) error {
	return nil
}

func TestNewDataSource(t *testing.T) {

	expected := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name_suffix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}

	actual := NewDataSource((*TestDataSource)(nil))
	if actual.Read == nil {
		t.Fail()
	}

	actual.Read = nil

	if !reflect.DeepEqual(expected, actual) {
		t.Fail()
	}
}
