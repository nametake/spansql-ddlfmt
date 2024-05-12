package main

import (
	"fmt"

	"cloud.google.com/go/spanner/spansql"
)

func FormatDDL(ddlStr string) (string, error) {
	parsedDDL, err := spansql.ParseDDL("f", ddlStr)
	if err != nil {
		return "", fmt.Errorf("parse DDL: %v", err)
	}
	var strSQL string

	for _, ddl := range parsedDDL.List {
		strSQL += ddl.SQL() + ";\n"
	}

	return strSQL, nil
}
