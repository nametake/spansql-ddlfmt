package main

import (
	"fmt"
	"slices"
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
	return fmt.Sprintf("%s;", d.ddl.SQL())
}

type CommentItem struct {
	comment *spansql.Comment
}

func (c *CommentItem) Pos() spansql.Position {
	return c.comment.Pos()
}

func (c *CommentItem) SQL() string {
	comments := make([]string, 0, len(c.comment.Text))
	for _, text := range c.comment.Text {
		comments = append(comments, fmt.Sprintf("-- %s", text))
	}
	return strings.Join(comments, "\n")
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

	for _, comment := range parsedDDL.Comments {
		items = append(items, &CommentItem{comment})
	}

	slices.SortFunc(items, func(a, b Item) int {
		return a.Pos().Line - b.Pos().Line
	})

	var sqls []string
	for _, item := range items {
		sql := fmt.Sprintf("%s", item.SQL())
		sqls = append(sqls, sql)
	}

	return strings.Join(sqls, "\n\n"), nil
}
