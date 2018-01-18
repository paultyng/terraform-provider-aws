package structprovider

import (
	"reflect"

	"github.com/hashicorp/terraform/helper/schema"
)

type DataSource interface {
	Identifier
	Reader
}

func NewDataSource(s interface{}) *schema.Resource {
	if _, ok := s.(DataSource); !ok {
		panic("type is not a DataSource")
	}

	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("a struct type is required to register a data source")
	}

	sr := newStructResource(t)

	r := &schema.Resource{
		Read:   sr.Read,
		Schema: sr.Schema,
	}

	return r
}
