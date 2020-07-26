package sqlite

import (
	"database/sql"
	"github.com/getevo/dbo"
	_ "github.com/mattn/go-sqlite3"
)

func Connect(dsn string) (*dbo.Conn,error)  {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil{
		return nil,err
	}
	conn := dbo.Conn{}
	conn.SetConn(db)
	conn.Setup()
	return &conn,nil
}