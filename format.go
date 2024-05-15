package ddlfmt

import (
	"fmt"
	"slices"
	"strings"

	"cloud.google.com/go/spanner/spansql"
)

type Item interface {
	Pos() spansql.Position
	SQL(ddl *spansql.DDL) string
}

type DDLItem struct {
	ddlStmt        spansql.DDLStmt
	leadingComment *spansql.Comment
}

func NewDDLItem(ddl *spansql.DDL, ddlStmt spansql.DDLStmt) *DDLItem {
	leadingComment := ddl.LeadingComment(ddlStmt)
	return &DDLItem{
		ddlStmt:        ddlStmt,
		leadingComment: leadingComment,
	}
}

func (d *DDLItem) Pos() spansql.Position {
	return d.ddlStmt.Pos()
}

func (d *DDLItem) SQL(ddl *spansql.DDL) string {
	var sql string

	if d.leadingComment != nil {
		cmt := formatComment(d.leadingComment, "")
		sql = fmt.Sprintf("%s\n", cmt)
		diff := d.ddlStmt.Pos().Line - d.leadingComment.Pos().Line
		if diff != 1 {
			sql += "\n"
		}
	}
	switch stmt := d.ddlStmt.(type) {
	case *spansql.CreateTable:
		str := "CREATE TABLE "
		if stmt.IfNotExists {
			str += "IF NOT EXISTS "
		}
		str += stmt.Name.SQL() + " (\n"
		for _, c := range stmt.Columns {
			if leadingComment := ddl.LeadingComment(c); leadingComment != nil {
				str += formatComment(leadingComment, "  ") + "\n"
			}
			str += "  " + c.SQL() + ","
			if inlineComment := ddl.InlineComment(c); inlineComment != nil {
				str += " " + formatComment(inlineComment, "")
			}
			str += "\n"
		}
		for _, tc := range stmt.Constraints {
			str += "  " + tc.SQL() + ",\n"
		}
		if len(stmt.Synonym) > 0 {
			str += "  SYNONYM(" + stmt.Synonym.SQL() + "),\n"
		}
		str += ") PRIMARY KEY("
		for i, c := range stmt.PrimaryKey {
			if i > 0 {
				str += ", "
			}
			str += c.SQL()
		}
		str += ")"
		if il := stmt.Interleave; il != nil {
			str += ",\n  INTERLEAVE IN PARENT " + il.Parent.SQL() + " ON DELETE " + il.OnDelete.SQL()
		}
		if rdp := stmt.RowDeletionPolicy; rdp != nil {
			str += ",\n  " + rdp.SQL()
		}
		sql += fmt.Sprintf("%s;", str)
	default:
		sql += fmt.Sprintf("%s;", stmt.SQL())
	}
	return sql
}

type CommentItem struct {
	comment *spansql.Comment
}

func (c *CommentItem) Pos() spansql.Position {
	return c.comment.Pos()
}

func (c *CommentItem) SQL(_ *spansql.DDL) string {
	return formatComment(c.comment, "")
}

func formatComment(comment *spansql.Comment, indent string) string {
	comments := make([]string, 0, len(comment.Text))
	for _, text := range comment.Text {
		comments = append(comments, fmt.Sprintf("%s-- %s", indent, text))
	}
	return strings.Join(comments, "\n")
}

func FormatDDL(filename, ddlStr string) (string, error) {
	ddl, err := spansql.ParseDDL(filename, ddlStr)
	if err != nil {
		return "", fmt.Errorf("parse DDL: %v", err)
	}

	commentsMap := make(map[int]*spansql.Comment)
	for _, comment := range ddl.Comments {
		commentsMap[comment.Pos().Line] = comment
	}

	items := make([]Item, 0)

	for _, ddlStmt := range ddl.List {
		leadingComment := ddl.LeadingComment(ddlStmt)
		if leadingComment != nil {
			delete(commentsMap, leadingComment.Pos().Line)
		}
		switch stmt := ddlStmt.(type) {
		case *spansql.CreateTable:
			for _, c := range stmt.Columns {
				if leadingComment := ddl.LeadingComment(c); leadingComment != nil {
					delete(commentsMap, leadingComment.Pos().Line)
				}
				if inlineComment := ddl.InlineComment(c); inlineComment != nil {
					delete(commentsMap, inlineComment.Pos().Line)
				}
			}
		}

		items = append(items, NewDDLItem(ddl, ddlStmt))
	}

	for _, comment := range commentsMap {
		items = append(items, &CommentItem{comment})
	}

	slices.SortFunc(items, func(i, j Item) int {
		return i.Pos().Line - j.Pos().Line
	})

	var sqls []string
	for _, item := range items {
		sqls = append(sqls, item.SQL(ddl))
	}

	return strings.Join(sqls, "\n\n"), nil
}
