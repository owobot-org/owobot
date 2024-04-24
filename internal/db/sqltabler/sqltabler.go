package sqltabler

import (
	"errors"
	"io"
	"strings"

	sqlparser "github.com/rqlite/sql"
)

// Modify adds a prefix and suffix to every table name found in stmt.
func Modify(stmt, prefix, suffix string) (string, error) {
	parser := sqlparser.NewParser(strings.NewReader(stmt))
	sb := strings.Builder{}
	for {
		s, err := parser.ParseStatement()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return "", err
		}
		modify(s, prefix, suffix)
		sb.WriteString(s.String())
		sb.WriteByte(';')
	}
	return sb.String(), nil
}

// modify changes all the table, viee, trigger, and index names in a single statement
func modify(stmt any, prefix, suffix string) {
	switch stmt := stmt.(type) {
	case *sqlparser.SelectStatement:
		modifySource(stmt.Source, prefix, suffix)
		modify(stmt.WhereExpr, prefix, suffix)
	case *sqlparser.InsertStatement:
		stmt.Table.Name = prefix + stmt.Table.Name + suffix
		if stmt.Select != nil {
			modify(stmt.Select, prefix, suffix)
		}
	case *sqlparser.UpdateStatement:
		stmt.Table.Name.Name = prefix + stmt.Table.Name.Name + suffix
		for _, assignment := range stmt.Assignments {
			modify(assignment, prefix, suffix)
		}
	case *sqlparser.CreateTableStatement:
		stmt.Name.Name = prefix + stmt.Name.Name + suffix
		if stmt.Select != nil {
			modify(stmt.Select, prefix, suffix)
		}
		for _, col := range stmt.Columns {
			modify(col, prefix, suffix)
		}
		for _, constraint := range stmt.Constraints {
			modify(constraint, prefix, suffix)
		}
	case *sqlparser.CreateViewStatement:
		stmt.Name.Name = prefix + stmt.Name.Name + suffix
		if stmt.Select != nil {
			modify(stmt.Select, prefix, suffix)
		}
	case *sqlparser.AlterTableStatement:
		stmt.Name.Name = prefix + stmt.Name.Name + suffix
		if stmt.NewName != nil {
			stmt.NewName.Name = prefix + stmt.NewName.Name + suffix
		}
		if stmt.ColumnDef != nil {
			modify(stmt.ColumnDef, prefix, suffix)
		}
	case *sqlparser.Call:
		for _, arg := range stmt.Args {
			modify(arg, prefix, suffix)
		}
	case *sqlparser.FilterClause:
		modify(stmt.X, prefix, suffix)
	case *sqlparser.DeleteStatement:
		stmt.Table.Name.Name = prefix + stmt.Table.Name.Name + suffix
		if stmt.WhereExpr != nil {
			modify(stmt.WhereExpr, prefix, suffix)
		}
	case *sqlparser.AnalyzeStatement:
		stmt.Name.Name = prefix + stmt.Name.Name + suffix
	case *sqlparser.ExplainStatement:
		modify(stmt.Stmt, prefix, suffix)
	case *sqlparser.CreateIndexStatement:
		stmt.Name.Name = prefix + stmt.Name.Name + suffix
		stmt.Table.Name = prefix + stmt.Table.Name + suffix
	case *sqlparser.CreateTriggerStatement:
		stmt.Name.Name = prefix + stmt.Name.Name + suffix
		stmt.Table.Name = prefix + stmt.Table.Name + suffix
		if stmt.WhenExpr != nil {
			modify(stmt.WhenExpr, prefix, suffix)
		}
		for _, istmt := range stmt.Body {
			modify(istmt, prefix, suffix)
		}
	case *sqlparser.CTE:
		stmt.TableName.Name = prefix + stmt.TableName.Name + suffix
		if stmt.Select != nil {
			modify(stmt.Select, prefix, suffix)
		}
	case *sqlparser.DropTableStatement:
		stmt.Name.Name = prefix + stmt.Name.Name + suffix
	case *sqlparser.DropViewStatement:
		stmt.Name.Name = prefix + stmt.Name.Name + suffix
	case *sqlparser.DropIndexStatement:
		stmt.Name.Name = prefix + stmt.Name.Name + suffix
	case *sqlparser.DropTriggerStatement:
		stmt.Name.Name = prefix + stmt.Name.Name + suffix
	case *sqlparser.ForeignKeyConstraint:
		stmt.ForeignTable.Name = prefix + stmt.ForeignTable.Name + suffix
	case *sqlparser.OnConstraint:
		modify(stmt.X, prefix, suffix)
	case *sqlparser.ExprList:
		for _, expr := range stmt.Exprs {
			modify(expr, prefix, suffix)
		}
	case *sqlparser.UnaryExpr:
		modify(stmt.X, prefix, suffix)
	case *sqlparser.BinaryExpr:
		modify(stmt.X, prefix, suffix)
		modify(stmt.Y, prefix, suffix)
	case *sqlparser.ParenExpr:
		modify(stmt.X, prefix, suffix)
	case *sqlparser.CastExpr:
		modify(stmt.X, prefix, suffix)
	case *sqlparser.OrderingTerm:
		modify(stmt.X, prefix, suffix)
	case *sqlparser.Assignment:
		modify(stmt.Expr, prefix, suffix)
	case *sqlparser.ColumnDefinition:
		for _, constraint := range stmt.Constraints {
			modify(constraint, prefix, suffix)
		}
	case *sqlparser.QualifiedRef:
		if stmt.Table != nil {
			stmt.Table.Name = prefix + stmt.Table.Name + suffix
		}
	case *sqlparser.CaseExpr:
		modify(stmt.ElseExpr, prefix, suffix)
		for _, block := range stmt.Blocks {
			modify(block.Condition, prefix, suffix)
			modify(block.Body, prefix, suffix)
		}
	}
}

func modifySource(source sqlparser.Source, prefix, suffix string) {
	switch source := source.(type) {
	case *sqlparser.QualifiedTableName:
		source.Name.Name = prefix + source.Name.Name + suffix
	case *sqlparser.JoinClause:
		modifySource(source.X, prefix, suffix)
		modifySource(source.Y, prefix, suffix)
		modify(source.Constraint, prefix, suffix)
	}
}
