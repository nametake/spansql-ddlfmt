package ddlfmt

import (
	"fmt"
	"strings"

	"cloud.google.com/go/spanner/spansql"
)

type Item interface {
	Pos() spansql.Position
	SQL(ddl *spansql.DDL) string
}

type DDLItem struct {
	ddl spansql.DDLStmt
}

func (d *DDLItem) Pos() spansql.Position {
	return d.ddl.Pos()
}

func (d *DDLItem) SQL(ddl *spansql.DDL) string {
	var sql string
	leadingComment := ddl.LeadingComment(d.ddl)

	if leadingComment != nil {
		fmt.Println(leadingComment)
		cmt := formatComment(leadingComment, "")
		sql = fmt.Sprintf("%s\n", cmt)
		diff := leadingComment.Pos().Line - d.ddl.Pos().Line
		if diff != 1 {
			sql += "\n"
		}
	}
	sql += fmt.Sprintf("%s;", d.ddl.SQL())
	return sql
}

func formatComment(comment *spansql.Comment, indent string) string {
	comments := make([]string, 0, len(comment.Text))
	for _, text := range comment.Text {
		comments = append(comments, fmt.Sprintf("%s-- %s", indent, text))
	}
	return strings.Join(comments, "\n")
}

func FormatDDL(filename, ddlStr string) (string, error) {
	parsedDDL, err := spansql.ParseDDL(filename, ddlStr)
	if err != nil {
		return "", fmt.Errorf("parse DDL: %v", err)
	}

	items := make([]Item, 0)

	for _, ddl := range parsedDDL.List {
		items = append(items, &DDLItem{ddl})
	}

	var sqls []string
	for _, item := range items {
		sqls = append(sqls, item.SQL(parsedDDL))
	}

	return strings.Join(sqls, "\n\n"), nil
}
