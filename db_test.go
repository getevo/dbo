package dbo_test

import (
	"fmt"
	"github.com/getevo/dbo/sqlite"
	"testing"
)

type Contact struct {
	ContactID int
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

func TestModelInsert(t *testing.T) {
	dbo,err := sqlite.Connect("./test.sqlite3")
	if err != nil{
		fmt.Println(err)
		return
	}



}

func TestSelectRawQuery(t *testing.T) {
	dbo,err := sqlite.Connect("./test.sqlite3")
	if err != nil{
		fmt.Println(err)
		return
	}

	x := []string{}
	err = dbo.Query(`
		SELECT phone FROM contacts 
	`).Scan(&x)
	fmt.Println(x)

	var y string
	err = dbo.Query(`
		SELECT phone FROM contacts 
	`).Scan(&y)
	fmt.Println(y)


	var z map[string]interface{}
	err = dbo.Query(`
		SELECT * FROM contacts 
	`).Scan(&z)
	fmt.Println(z,err)

	var f []map[string]interface{}
	err = dbo.Query(`
		SELECT * FROM contacts 
	`).Scan(&f)
	fmt.Println(f,err)

	var n []Contact
	err = dbo.Query(`
		SELECT * FROM contacts 
	`).Scan(&n)
	fmt.Println(n,err)

	var m Contact
	err = dbo.Query(`
		SELECT * FROM contacts 
	`).Scan(&m)
	fmt.Println(m,err)
}
