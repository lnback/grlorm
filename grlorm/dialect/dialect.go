package dialect

import "reflect"
//保存一个支持数据的map
var dialectsMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string//数据类型，用于将Go语言的类型转换为该数据库的数据类型
	TableExistSQL(tableName string) (string,[]interface{}) //表是否存在
}
//注册dialect到map中
func RegisterDialect(name string,dialect Dialect)  {
	dialectsMap[name] = dialect
}

//从map中获取dialect
func GetDialect(name string) (dialect Dialect,ok bool){
	dialect,ok = dialectsMap[name]
	return
}
