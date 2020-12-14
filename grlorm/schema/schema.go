package schema

import (
	"go/ast"
	"grlorm/dialect"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag string
}

type Schema struct {
	Model interface{}
	Name string
	Fields []*Field
	FieldName []string
	fieldMap map[string]*Field
}

func (s * Schema) GetField(name string)  *Field{
	return s.fieldMap[name]
}

func Parse(dest interface{},d dialect.Dialect)  *Schema{
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model: dest,
		Name: modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0;i < modelType.NumField();i++{
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name){
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v,ok := p.Tag.Lookup("grlorm");ok{
				field.Tag = v
			}
			schema.Fields = append(schema.Fields,field)
			schema.FieldName = append(schema.FieldName,p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}