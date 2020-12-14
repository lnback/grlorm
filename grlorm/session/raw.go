package session

import (
	"database/sql"
	"grlorm/dialect"
	"grlorm/log"
	"grlorm/schema"
	"strings"
)

type Session struct {
	db * sql.DB
	dialect dialect.Dialect
	refTable * schema.Schema
	sql strings.Builder
	sqlVars []interface{}
}
//初始化
func New(db * sql.DB,dialect dialect.Dialect) * Session{
	return &Session{
		db: db,
		dialect: dialect,
	}
}
//清空
func (s *Session) Clear(){
	s.sql.Reset()
	s.sqlVars = nil
}
func (s * Session)  DB() * sql.DB{
	return s.db
}

func (s * Session) Raw(sql string,values ...interface{})  *Session{
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars,values...)
	return s
}

func (s * Session) Exec() (result sql.Result, err error){
	defer s.Clear()
	log.Info(s.sql.String(),s.sqlVars)
	if result,err = s.DB().Exec(s.sql.String(),s.sqlVars...); err != nil{
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s * Session) QueryRows() (rows * sql.Rows,err error)  {
	defer s.Clear()
	log.Info(s.sql.String(),s.sqlVars)
	if rows , err = s.DB().Query(s.sql.String(),s.sqlVars...); err != nil{
		log.Error(err)
	}
	return
}
