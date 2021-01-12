package schema

import (
	"go/ast"
	"grlorm/dialect"
	"reflect"
)

type Field struct {
	Name string //变量名
	Type string //类型
	Tag string //标签
}

type Schema struct {
	Model interface{}//被映射的对象
	Name string //表名
	Fields []*Field //字段fields
	FieldNames []string //包含所有的字段名（列名）
	fieldMap map[string]*Field //记录字段名和field的映射关系
}

func (s * Schema) GetField(name string)  *Field{
	return s.fieldMap[name]
}

func Parse(dest interface{},d dialect.Dialect)  *Schema{
	//获取对象类型 最后输出结构体的名字 用此作为表名
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model: dest,
		Name: modelType.Name(),
		fieldMap: make(map[string]*Field),
	}
	//将结构体中的字段名和数据库中的列名进行对应 numfield()可以计算有多少个字段名
	for i := 0;i < modelType.NumField();i++{
		p := modelType.Field(i)
		//是否以大写开头
		if !p.Anonymous && ast.IsExported(p.Name){
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			//是否有tag 有lookup在grlorm这个前缀中寻找
			if v,ok := p.Tag.Lookup("grlorm");ok{
				field.Tag = v
			}
			//字段和列对应起来
			schema.Fields = append(schema.Fields,field)
			schema.FieldNames = append(schema.FieldNames,p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

func (s * Schema) RecordValues(dest interface{})[]interface{}  {
	destValue := reflect.Indirect(reflect.ValueOf(dest))

	var fieldValues []interface{}

	for _,field := range s.Fields{
		fieldValues = append(fieldValues,destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}