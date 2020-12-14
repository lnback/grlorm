package grlorm

import (
	"database/sql"
	"grlorm/dialect"
	"grlorm/log"
	"grlorm/session"
)

type Engine struct {
	db * sql.DB

	dialect dialect.Dialect
}

func NewEngine(driver ,source string) (engine * Engine,err error){
	db,err := sql.Open(driver,source)

	if err != nil{
		log.Error(err)
		return
	}

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
	if err := e.db.Close(); err != nil{
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (e * Engine) NewSession() * session.Session{
	return session.New(e.db,e.dialect)
}