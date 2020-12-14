package grlorm

import (
	"database/sql"
	"grlorm/log"
	"grlorm/session"
)

type Engine struct {
	db * sql.DB
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

	engine = &Engine{
		db: db,
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

func (e * Engine)  NewSession() * session.Session{
	return session.New(e.db)
}