package clause

import (
	"strings"
)

type Clause struct {
	//存放一个sql语句的map
	sql map[Type]string
	//存放一个vars的map 每一个value都是二维数组
	sqlVars map[Type][] interface{}
}

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

func (c * Clause) Set(name Type,vars ...interface{})  {
	if c.sql == nil{
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql,vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

func (c * Clause) Build(orders ...Type)(string,[]interface{})  {
	var sqls []string
	var vars []interface{}

	//根据type的顺序来生成sql
	for _,order := range orders{
		if sql,ok := c.sql[order];ok{
			sqls = append(sqls,sql)
			vars = append(vars,c.sqlVars[order]...)
		}
	}
	//生成一个 sql sql sql sql的sql语句
	return strings.Join(sqls," "),vars
}