package structprovider

import (
	"reflect"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	tagComputed = "computed"
	tagOptional = "optional"
	tagRequired = "required"
	tagForceNew = "forcenew"
)

type Resource interface {
	Identifier
	Reader
	Create(interface{}) error
	Delete(interface{}) error
}

type ResourceUpdater interface {
	Resource
	Update(interface{}) error
}

func NewResource(s interface{}) *schema.Resource {
	if _, ok := s.(Resource); !ok {
		panic("type is not a Resource")
	}

	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic("a struct type is required to register a resource")
	}

	sr := newStructResource(t)

	r := &schema.Resource{
		Create: sr.Create,
		Read:   sr.Read,
		Delete: sr.Delete,
		Schema: sr.Schema,
	}

	if _, ok := s.(ResourceUpdater); ok {
		r.Update = sr.Update
	}

	return r

}
