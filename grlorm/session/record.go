package session

import (
	"errors"
	"fmt"
	"grlorm/clause"
	"reflect"
)

func (s * Session) Insert(values ...interface{}) (int64,error) {
	//值
	recordValues := make([]interface{},0)
	// 循环值
	for _,value := range values {
		//得到表
		table := s.Model(values).RefTable()
		//设定clause的表、列、sql语句的方式
		s.clause.Set(clause.INSERT,table.Name,table.FieldNames)
		//值追加到values中
		recordValues = append(recordValues,table.RecordValues(value))
	}
	//值放到clause
	s.clause.Set(clause.VALUES,recordValues...)
	//构建最终的sql语句
	sql,vars := s.clause.Build(clause.INSERT,clause.VALUES)
	//执行
	result,err := s.Raw(sql,vars...).Exec()

	if err != nil{
		return 0, err
	}
	//返回增加的行数
	return result.RowsAffected()
}

func (s * Session) Find(values interface{}) error  {
	//values -> slice
	destSlice := reflect.Indirect(reflect.ValueOf(values))

	fmt.Println(destSlice)

	//结构体类型
	destType := destSlice.Type().Elem()

	//得到表 s.model返回一个session的指针，用这个指针获取reftable
	// new根据类型生成一个零值的指针
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	//设置表、列、sql语句
	s.clause.Set(clause.SELECT,table.Name,table.FieldNames)

	//build一个sql 并把之前set的值返回过来
	sql,vars := s.clause.Build(clause.SELECT,clause.WHERE,clause.ORDERBY,clause.LIMIT)
	//执行查询
	rows , err := s.Raw(sql,vars...).QueryRows()
	if err != nil{
		return err
	}
	//处理返回的数据
	for rows.Next(){
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _,name := range table.FieldNames{
			//把
			values = append(values,dest.FieldByName(name).Addr().Interface())
		}
		fmt.Println(values)

		//把数据写入values
		if err := rows.Scan(values...); err != nil{
			return err
		}
		destSlice.Set(reflect.Append(destSlice,dest))
	}
	return rows.Close()
}

func (s * Session) Update(kv ...interface{}) (int64,error){
	m,ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})

		for i := 0; i < len(kv);i+=2{
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE,s.RefTable().Name,m)

	sql,vars := s.clause.Build(clause.UPDATE,clause.WHERE)

	result,err := s.Raw(sql,vars...).Exec()

	if err != nil{
		return 0, err
	}
	return result.RowsAffected()

}

func (s * Session) Delete()(int64,error){
	s.clause.Set(clause.DELETE,s.RefTable().Name)
	sql,vars := s.clause.Build(clause.DELETE,clause.WHERE)
	result,err := s.Raw(sql,vars...).Exec()
	if err != nil{
		return 0, err
	}
	return result.RowsAffected()
}

func (s * Session) Count() (int64,error) {
	s.clause.Set(clause.COUNT,s.RefTable().Name)
	sql,vars := s.clause.Build(clause.COUNT,clause.WHERE)
	row := s.Raw(sql,vars...).QueryRow()

	var tmp int64

	if err := row.Scan(&tmp);err != nil{

		return 0, err
	}

	return tmp,nil
}
func (s * Session) Limit(num int) *Session  {
	s.clause.Set(clause.LIMIT,num)
	return s
}

func (s * Session) Where(desc string,args ...interface{})  *Session{
	var vars []interface{}

	s.clause.Set(clause.WHERE,append(append(vars,desc),args...)...)

	return s
}

func (s * Session) OrderBy(desc string)  *Session{
	s.clause.Set(clause.ORDERBY,desc)
	return s
}

func (s * Session) First(value interface{}) error  {
	dest := reflect.Indirect(reflect.ValueOf(value))

	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()

	if err := s.Limit(1).Find(destSlice.Addr().Interface());err != nil{
		return err
	}

	if destSlice.Len() == 0{
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))

	return nil
}