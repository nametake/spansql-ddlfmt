package main

import (
	"fmt"
	"strings"

	"cloud.google.com/go/spanner/spansql"
)

type DDLItem struct {
	ddl spansql.DDLStmt
}

func (d *DDLItem) Pos() spansql.Position {
	return d.ddl.Pos()
}

type CommentItem struct {
	comment *spansql.Comment
}

func (c *CommentItem) Pos() spansql.Position {
	return c.comment.Pos()
}

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
