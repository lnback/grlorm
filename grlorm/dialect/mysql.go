package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct {}

var _ Dialect = (*mysql)(nil)
func init()  {
	RegisterDialect("mysql",&mysql{})
}

func (m mysql) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int16, reflect.Int32:
		return "int"
	case reflect.Uint8:
		return "tinyint unsigned"
	case reflect.Int8:
		return "tinyint"
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "int unsigned"

	case reflect.Int64:
		return "bigint unsigned"
	case reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "double"
	case reflect.String:
		return "varchar(255)"
	case reflect.Array, reflect.Slice:
		return "varbinary(255)"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (m mysql) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT table_name FROM information_schema.TABLES WHERE table_name = ?;",args
}



