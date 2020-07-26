package dbo

import (
	"database/sql"
	"reflect"
	"strings"
)

func scanRowsToSimpleSlice(rows *sql.Rows,p reflect.Value)  {
	sliceType := reflect.TypeOf(p.Elem().Interface()).Elem()
	elemSlice := reflect.MakeSlice(reflect.SliceOf(sliceType),0,0)
	for rows.Next(){
		obj := reflect.New(sliceType).Interface()
		rows.Scan(obj)
		elemSlice = reflect.Append(elemSlice, reflect.ValueOf(obj).Elem())
	}
	p.Elem().Set(elemSlice)
	rows.Close()
}

func scanModel(rows *sql.Rows, target reflect.Value) error {
	st := target.Elem().Type()
	stType := target.Elem().Type()
	cols,err := rows.Columns()
	if err != nil{
		return err
	}
	var columnPtr = make([]int,len(cols),len(cols))
	for i,col := range cols{

		found := false
		for j := 0; j < st.NumField(); j++ {
			field := st.Field(j)
			if strings.ToLower(field.Name) ==  strings.ToLower(ToCamel(col)) || field.Tag.Get("db") == col{
				found = true
				columnPtr[i] = j
				break
			}
		}
		if !found{
			columnPtr[i] = -1
		}
	}
	var nullPtr interface{}

	for rows.Next(){
		var ptr = make([]interface{},len(columnPtr),len(columnPtr))
		ref := reflect.New(stType)

		for i,item := range columnPtr{
			if item == -1{
				ptr[i]  = &nullPtr
			}else {
				ptr[i] = ref.Elem().Field(item).Addr().Interface()

			}
		}
		rows.Scan(ptr...)
		target.Elem().Set(ref.Elem())
		break
	}



	return nil
}


func scanSliceModel(rows *sql.Rows, target reflect.Value) error {
	st := target.Elem().Type().Elem()
	stType := target.Elem().Type().Elem()
	cols,err := rows.Columns()
	if err != nil{
		return err
	}
	var columnPtr = make([]int,len(cols),len(cols))
	for i,col := range cols{

		found := false
		for j := 0; j < st.NumField(); j++ {
			field := st.Field(j)
			if strings.ToLower(field.Name) ==  strings.ToLower(ToCamel(col)) || field.Tag.Get("db") == col{
				found = true
				columnPtr[i] = j
				break
			}
		}
		if !found{
			columnPtr[i] = -1
		}
	}
	var nullPtr interface{}
	var slice  = reflect.MakeSlice(reflect.SliceOf(stType),0,0)
	for rows.Next(){
		var ptr = make([]interface{},len(columnPtr),len(columnPtr))
		ref := reflect.New(stType)

		for i,item := range columnPtr{
			if item == -1{
				ptr[i]  = &nullPtr
			}else {
				ptr[i] = ref.Elem().Field(item).Addr().Interface()

			}
		}
		rows.Scan(ptr...)
		slice = reflect.Append(slice,ref.Elem())
	}
	target.Elem().Set(slice)


	return nil
}
