package structprovider

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/paultyng/snake"
)

type structResource struct {
	reflect.Type
	Mapping map[string]string
	Schema  map[string]*schema.Schema
}

type Identifier interface {
	ID() string
}

type Reader interface {
	Read(interface{}) error
}

func newStructResource(t reflect.Type) *structResource {
	sr := &structResource{
		Type:    t,
		Mapping: map[string]string{},
		Schema:  map[string]*schema.Schema{},
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			//TODO: flatten embed?
			panic("not implemented")
		}
		name, opts := parseTag(f.Tag.Get("tf"))
		switch name {
		case "-":
			continue
		case "":
			name = fieldNameToSchemaName(f.Name)
		}
		sr.Mapping[name] = f.Name
		s := &schema.Schema{
			Type:     schemaType(f.Type),
			Computed: opts.Contains(tagComputed),
			Optional: opts.Contains(tagOptional),
			Required: opts.Contains(tagRequired),
			ForceNew: opts.Contains(tagForceNew),
		}
		sr.Schema[name] = s
	}

	return sr
}

func (sr *structResource) Create(d *schema.ResourceData, meta interface{}) error {
	v, err := sr.Get(d)
	if err != nil {
		return err
	}
	r, ok := v.(Resource)
	if !ok {
		panic(fmt.Sprintf("type %s does not implement Resource", sr.Name()))
	}

	err = r.Create(meta)
	if err != nil {
		return err
	}

	err = sr.Set(v, d)
	if err != nil {
		return err
	}

	err = r.Read(meta)
	if err != nil {
		return err
	}

	err = sr.Set(v, d)
	if err != nil {
		return err
	}

	return nil
}

func (sr *structResource) Read(d *schema.ResourceData, meta interface{}) error {
	v, err := sr.Get(d)
	if err != nil {
		return err
	}
	if r, ok := v.(Reader); ok {
		err = r.Read(meta)
		if err != nil {
			return err
		}

		err = sr.Set(v, d)
		if err != nil {
			return err
		}
	}
	// is this an error?
	return nil
}

func (sr *structResource) Update(d *schema.ResourceData, meta interface{}) error {
	v, err := sr.Get(d)
	if err != nil {
		return err
	}

	r, ok := v.(ResourceUpdater)
	if !ok {
		panic(fmt.Sprintf("type %s does not implement ResourceUpdater", sr.Name()))
	}

	err = r.Update(meta)
	if err != nil {
		return err
	}

	err = sr.Set(v, d)
	if err != nil {
		return err
	}

	err = r.Read(meta)
	if err != nil {
		return err
	}

	err = sr.Set(v, d)
	if err != nil {
		return err
	}

	return nil
}

func (sr *structResource) Delete(d *schema.ResourceData, meta interface{}) error {
	v, err := sr.Get(d)
	if err != nil {
		return err
	}

	r, ok := v.(Resource)
	if !ok {
		panic(fmt.Sprintf("type %s does not implement Resource", sr.Name()))
	}

	err = r.Delete(meta)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func (sr *structResource) Get(d *schema.ResourceData) (interface{}, error) {
	ri := reflect.New(sr.Type)

	for dn, fn := range sr.Mapping {
		if raw, ok := d.GetOk(dn); ok {
			fv := ri.Elem().FieldByName(fn)
			fv.Set(reflect.ValueOf(raw))
		}
	}

	return ri.Interface(), nil
}

func (sr *structResource) Set(r interface{}, d *schema.ResourceData) error {
	if ider, ok := r.(Identifier); ok {
		id := ider.ID()
		d.SetId(id)
	}

	rv := reflect.ValueOf(r).Elem()

	for dn, fn := range sr.Mapping {
		fv := rv.FieldByName(fn)
		raw := fv.Interface()
		err := d.Set(dn, raw)
		if err != nil {
			return err
		}
	}

	return nil
}

func fieldNameToSchemaName(f string) string {
	return snake.ToSnake(f)
}

func schemaType(t reflect.Type) schema.ValueType {
	switch t.Kind() {
	case reflect.Bool:
		return schema.TypeBool
	case reflect.String:
		return schema.TypeString
	default:
		panic(fmt.Sprintf("unrecognized kind %v", t.Kind()))
	}
}
