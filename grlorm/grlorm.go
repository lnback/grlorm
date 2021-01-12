package grlorm

import (
	"database/sql"
	"grlorm/dialect"
	"grlorm/log"
	"grlorm/session"
)


type Engine struct {
	db * sql.DB //获得一个数据库连接指针

	dialect dialect.Dialect
}

func NewEngine(driver ,source string) (engine * Engine,err error){
	//开一个数据库连接 后面用这个db去获得session
	db,err := sql.Open(driver,source)


	if err != nil{
		log.Error(err)
		return
	}
	//make sure the  connection is alive
	if err = db.Ping(); err != nil{
		log.Error(err)
		return
	}

	dial , ok := dialect.GetDialect(driver)

	if !ok {
		log.Errorf("dialect %s Not Found",driver)
		return
	}

	engine = &Engine{
		db: db,
		dialect: dial,
	}
	
	log.Info("Connect database success")
	return
}

func (e * Engine) Close()  {
	//是否正常close
	if err := e.db.Close(); err != nil{
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (e * Engine) NewSession() * session.Session{
	//返回一个session
	return session.New(e.db,e.dialect)
}