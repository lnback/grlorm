package session

import (
	"database/sql"
	"grlorm/clause"
	"grlorm/dialect"
	"grlorm/log"
	"grlorm/schema"
	"strings"
)

type Session struct {
	db * sql.DB // 一个数据库的连接
	dialect dialect.Dialect //方言、也就是连接数据库的类型
	refTable * schema.Schema //依赖的表 schema
	sql strings.Builder // sql用一个builder来做，之后连接很方便
	sqlVars []interface{} //传入的值

	clause clause.Clause //生成sql
}
//初始化 构造函数
func New(db * sql.DB,dialect dialect.Dialect) * Session{
	return &Session{
		db: db,
		dialect: dialect,
	}
}
//清空sql
func (s *Session) Clear(){
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}
// 返回数据库的连接DB
func (s * Session)  DB() * sql.DB{
	return s.db
}
// 把sql和values写入sql和sqlvars中
func (s * Session) Raw(sql string,values ...interface{})  *Session{
	s.sql.WriteString(sql) //
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars,values...)
	return s
}
//add update delete
func (s * Session) Exec() (result sql.Result, err error){
	defer s.Clear()
	log.Info(s.sql.String(),s.sqlVars)
	//在这里调用s.db().exec()执行sql语句
	if result,err = s.DB().Exec(s.sql.String(),s.sqlVars...); err != nil{
		log.Error(err)
	}
	return
}
// queryRow 查询一行
func (s *Session) QueryRow() (*sql.Row){
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)

	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}
// 查询多行 rows
func (s * Session) QueryRows() (rows * sql.Rows,err error)  {
	defer s.Clear()
	log.Info(s.sql.String(),s.sqlVars)
	if rows , err = s.DB().Query(s.sql.String(),s.sqlVars...); err != nil{
		log.Error(err)
	}
	return
}

