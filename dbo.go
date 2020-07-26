package dbo

import (
	"database/sql"
	"time"
)

type Conn struct {
	db 		*sql.DB
	debug 	bool
}



func (conn *Conn)SetConn(db *sql.DB)  {
	conn.db = db
}

func (conn *Conn)Setup()  {

}

func (conn *Conn)Ping() error {
	return conn.db.Ping()
}

func (conn *Conn)SetConnMaxLifetime(d time.Duration) {
	conn.db.SetConnMaxLifetime(d)
}

func (conn *Conn)SetMaxIdleConns(n int) {
	conn.db.SetMaxIdleConns(n)
}

func (conn *Conn)SetMaxOpenConns(n int) {
	conn.db.SetMaxOpenConns(n)
}

func (conn *Conn)SetDebug(v bool) {
	conn.debug = v
}

func (conn *Conn)Stats() sql.DBStats {
	return conn.db.Stats()
}

func (conn *Conn)Connection() *sql.DB {
	return conn.db
}

func (conn *Conn)Debug() *Query {
	return &Query{
		debug:true,
		db:conn.db,
	}
}

func (conn *Conn)Query(query string,params... interface{}) *Query {
	return &Query{
		debug:true,
		query:query,
		params:params,
		db:conn.db,
	}
}
