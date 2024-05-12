package main

import (
	"fmt"
	"strings"

	"cloud.google.com/go/spanner/spansql"
)

func FormatDDL(ddlStr string) (string, error) {
	parsedDDL, err := spansql.ParseDDL("f", ddlStr)
	if err != nil {
		return "", fmt.Errorf("parse DDL: %v", err)
	}

	var sqls []string

	for _, ddl := range parsedDDL.List {
		sql := fmt.Sprintf("%s;", ddl.SQL())
		sqls = append(sqls, sql)
	}

	return strings.Join(sqls, "\n\n"), nil
}
