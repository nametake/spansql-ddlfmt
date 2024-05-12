package main

import (
	"fmt"
	"strings"

	"cloud.google.com/go/spanner/spansql"
)

type Item interface {
	Pos() spansql.Position
	SQL() string
}

type DDLItem struct {
	ddl spansql.DDLStmt
}

func (d *DDLItem) Pos() spansql.Position {
	return d.ddl.Pos()
}

func (d *DDLItem) SQL() string {
	return d.ddl.SQL()
}

type CommentItem struct {
	comment *spansql.Comment
}

func (c *CommentItem) Pos() spansql.Position {
	return c.comment.Pos()
}

func (c *CommentItem) SQL() string {
	// return c.comment.SQL()
	return ""
}

func FormatDDL(ddlStr string) (string, error) {
	parsedDDL, err := spansql.ParseDDL("f", ddlStr)
	if err != nil {
		return "", fmt.Errorf("parse DDL: %v", err)
	}

	items := make([]Item, 0)

	for _, ddl := range parsedDDL.List {
		items = append(items, &DDLItem{ddl})
	}

	var sqls []string
	for _, item := range items {
		sql := fmt.Sprintf("%s;", item.SQL())
		sqls = append(sqls, sql)
	}

	return strings.Join(sqls, "\n\n"), nil
}
