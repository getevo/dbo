package dbo

import (
	"database/sql"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"
)

var Logger = log.New(os.Stdout, "[sql]", log.LstdFlags)

type Query struct {
	debug  bool
	query  string
	params []interface{}
	db     *sql.DB
}
type Result struct {
	LastInsertID int64
	AffectedRows int64
	Error        error
}

func (conn *Query)Debug() *Query {
	conn.debug = true
	return conn
}

func (conn *Query)Scan(v interface{},params ...interface{})  error {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr{
		return fmt.Errorf("input interface is not a pointer")
	}
	target := reflect.ValueOf(v)


	switch target.Type().Elem().Kind() {
	case reflect.Struct:
		var rows *sql.Rows
		var err error

		if len(params) > 0{
			rows,err = conn.Rows(params...)
		}else{
			rows,err = conn.Rows()
		}
		defer rows.Close()
		if err != nil{
			return err
		}
		return scanModel(rows,target)
	case reflect.Map:
		var rows *sql.Rows
		var err error
		if len(params) > 0{
			rows,err = conn.Rows(params...)
		}else{
			rows,err = conn.Rows()
		}
		defer rows.Close()
		if err != nil{
			return err
		}
		cols,err := rows.Columns()
		if err != nil{
			return err
		}
		data := make(map[string]interface{})
		if _,ok := v.(*map[string]interface{}); ok{
			columns := make([]string, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i, _ := range columns {
				columnPointers[i] = &columns[i]
			}
			for rows.Next() {
				rows.Scan(columnPointers...)
				break
			}
			for i, colName := range cols {
				data[colName] = columns[i]
			}

			target.Elem().Set(reflect.ValueOf(data))

			break

		}else{
			return fmt.Errorf("input map is not *[]map[string]interface{}")
		}

	case reflect.Slice:
		var rows *sql.Rows
		var err error

		if len(params) > 0{
			rows,err = conn.Rows(params...)
		}else{
			rows,err = conn.Rows()
		}
		defer rows.Close()
		if err != nil{
			return err
		}
		cols,err := rows.Columns()
		if err != nil{
			return err
		}


		if _,ok := v.(*[]map[string]interface{}); ok {
			var data = []map[string]interface{}{}
			for rows.Next() {
				item := make(map[string]interface{})
				columns := make([]string, len(cols))
				columnPointers := make([]interface{}, len(cols))
				for i, _ := range columns {
					columnPointers[i] = &columns[i]
				}
				rows.Scan(columnPointers...)
				for i, colName := range cols {
					item[colName] = columns[i]
				}
				data = append(data,item)
			}

			target.Elem().Set(reflect.ValueOf(data))

		}else if target.Elem().Type().Elem().Kind() == reflect.Struct {
			return scanSliceModel(rows,target)
		}else{
			sliceType := reflect.TypeOf(target.Elem().Interface()).Elem()
			elemSlice := reflect.MakeSlice(reflect.SliceOf(sliceType),0,0)
			for rows.Next(){
				obj := reflect.New(sliceType).Interface()
				rows.Scan(obj)
				elemSlice = reflect.Append(elemSlice, reflect.ValueOf(obj).Elem())
			}
			target.Elem().Set(elemSlice)
		}


	default:
		var row *sql.Row
		if len(params) > 0{
			row = conn.Row(params...)
		}else{
			row = conn.Row()
		}
		return row.Scan(v)
	}

	return nil
}



func (conn *Query)Row(params... interface{})  *sql.Row {
	var start = time.Now().UnixNano()
	var result *sql.Row
	if len(params) > 0{
		if conn.debug{
			fmt.Println(getCaller())
			fmt.Println( SQLQueryDebugString(conn.query,params...) )
		}
		result = conn.db.QueryRow(conn.query,params...)
		if conn.debug{
			color.Red("Time: %f ms", passedTime(start) )
		}
	}else{
		if conn.debug{
			fmt.Println(getCaller())
			fmt.Println( SQLQueryDebugString(conn.query,conn.params...))
		}
		result = conn.db.QueryRow(conn.query,conn.params...)
		if conn.debug{
			color.Red("Time: %f ms", passedTime(start) )
		}
	}
	return result
}

func (conn *Query)Rows(params... interface{})  (*sql.Rows,error) {
	var start = time.Now().UnixNano()
	var result *sql.Rows
	var err error
	if len(params) > 0{
		if conn.debug{
			fmt.Println(getCaller())
			fmt.Println( SQLQueryDebugString(conn.query,params...) )
		}
		result,err = conn.db.Query(conn.query,params...)
		if conn.debug{
			color.Red("Time: %f ms", passedTime(start) )
			if err != nil {
				color.Red(err.Error())
			}
		}
	}else{
		if conn.debug{
			fmt.Println(getCaller())
			fmt.Println( SQLQueryDebugString(conn.query,conn.params...))
		}
		result,err = conn.db.Query(conn.query,conn.params...)
		if conn.debug{
			color.Red("Time: %f ms", passedTime(start) )
			if err != nil {
				color.Red(err.Error())
			}

		}
	}
	return result,err
}

func (conn *Query)Exec(params... interface{})  Result {
	var res sql.Result
	var err error
	var start = time.Now().UnixNano()
	var result Result
	if len(params) > 0{
		if conn.debug{
		 	fmt.Println(getCaller())
			fmt.Println( SQLQueryDebugString(conn.query,params...) )
		}
		res,err = conn.db.Exec(conn.query,params...)
		result = parseResult(res,err)
		if conn.debug && err != nil{
		 	color.Red(err.Error())
		}
		if conn.debug{
			color.Red("Time: %f ms", passedTime(start) )
			color.Red("Affected Rows: %d",result.AffectedRows)
		}
	}else{
		if conn.debug{
			fmt.Println(getCaller())
			fmt.Println( SQLQueryDebugString(conn.query,conn.params...))
		}
		res,err = conn.db.Exec(conn.query,conn.params...)
		result = parseResult(res,err)
		if conn.debug && err != nil{
			color.Red(err.Error())
		}
		if conn.debug{
			color.Red("Time: %f ms", passedTime(start) )
			color.Red("Affected Rows: %d\n----------",result.AffectedRows)
		}
	}

	return result
}

func passedTime(start int64) float64 {
	return float64(time.Now().UnixNano() - start)/1000000
}

func parseResult(s sql.Result, e error) Result {
	res := Result{}
	if s != nil {
		res.AffectedRows, _ = s.RowsAffected()
		res.LastInsertID, _ = s.LastInsertId()
	}
	res.Error = e
	return res
}

func getCaller() string  {
	_, filename, line, _ := runtime.Caller(3)
	return "\n----------\n"+time.Now().Format("2006-02-01 15:04:05 ")+ filepath.Base(filename)+":"+fmt.Sprint(line)
}