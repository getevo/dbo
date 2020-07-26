package dbo

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
)

func SQLQueryDebugString(query string, args ...interface{}) string {
	var buffer bytes.Buffer
	nArgs := len(args)
	// Break the string by question marks, iterate over its parts and for each
	// question mark - append an argument and format the argument according to
	// it's type, taking into consideration NULL values and quoting strings.
	for i, part := range strings.Split(query, "?") {
		buffer.WriteString(part)

		if i < nArgs {
			switch a := args[i].(type) {
			case int,int64,int32,int16,int8,uint32,uint,uint16,uint64:
				buffer.WriteString(fmt.Sprintf("'%d'", a))
			case bool:
				buffer.WriteString(fmt.Sprintf("%t", a))
			case sql.NullBool:
				if a.Valid {
					buffer.WriteString(fmt.Sprintf("%t", a.Bool))
				} else {
					buffer.WriteString("NULL")
				}
			case sql.NullInt64:
				if a.Valid {
					buffer.WriteString(fmt.Sprintf("'%d'", a.Int64))
				} else {
					buffer.WriteString("NULL")
				}
			case sql.NullString:
				if a.Valid {
					buffer.WriteString(fmt.Sprintf("'%q'", a.String))
				} else {
					buffer.WriteString("NULL")
				}
			case sql.NullFloat64:
				if a.Valid {
					buffer.WriteString(fmt.Sprintf("'%f'", a.Float64))
				} else {
					buffer.WriteString("NULL")
				}
			default:
				buffer.WriteString(fmt.Sprintf("'%q'", a))
			}
		}
	}
	return buffer.String()
}
