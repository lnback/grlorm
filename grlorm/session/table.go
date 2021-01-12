package session

import (
	"fmt"
	"grlorm/log"
	"grlorm/schema"
	"reflect"
	"strings"
)
//用于给reftable赋值，解析操作很耗时，所以将解析的结果放入reftable中
func (s *Session) Model(value interface{})  *Session{
	//只有满足reftable为空或者原reftable与现在的value不同才改变
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}
//返回reftable
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}
//数据库表的创建
func (s * Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _,field := range table.Fields{
		columns = append(columns,fmt.Sprintf("%s %s %s",field.Name,field.Type,field.Tag))
	}
	desc := strings.Join(columns,",")

	fmt.Println(desc)
	//调用raw生成sql语句并执行
	_,err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);",table.Name,desc)).Exec()

	return err
}

func (s * Session) DropTable() error  {
	_,err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s",s.RefTable().Name)).Exec()
	return err
}

func (s * Session) HasTable() bool{
	sql,values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql,values...).QueryRow()
	var tmp string

	_ = row.Scan(&tmp)

	return tmp == s.RefTable().Name
}