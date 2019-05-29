package stmt

import "github.com/getbud/bud/lab/sql/builder"

// Select ...
type Select struct {
	distinct          bool
	selectExpressions []SelectExpression
	tables            []Table
	whereConditions   []Condition
	alias             string
}

// NewSelect ...
func NewSelect(selectExpressions ...SelectExpression) Select {
	return Select{
		selectExpressions: selectExpressions,
	}
}

// Distinct ...
func (s Select) Distinct() Select {
	s.distinct = true
	return s
}

// Select ...
func (s Select) Select(expressions ...SelectExpression) Select {
	s.selectExpressions = append(s.selectExpressions, expressions...)
	return s
}

// From ...
func (s Select) From(tables ...Table) Select {
	s.tables = append(s.tables, tables...)
	return s
}

// Where ...
func (s Select) Where(conditions ...Condition) Select {
	s.whereConditions = append(s.whereConditions, conditions...)
	return s
}

// As ...
func (s Select) As(alias string) Select {
	s.alias = alias
	return s
}

// Build ...
func (s Select) Build() (string, []interface{}) {
	ctx := builder.NewContext()
	ctx.Write("SELECT ")

	if len(s.selectExpressions) > 0 {
		for i, expr := range s.selectExpressions {
			expr.BuildExpression(ctx)

			if alias := expr.Alias(); alias != "" {
				ctx.Write(" AS ")
				ctx.Write(alias)
			}

			if i < len(s.selectExpressions)-1 {
				ctx.Write(", ")
			}
		}
	} else {
		ctx.Write("*")
	}

	ctx.Write(" FROM ")

	// TODO: Change me...
	for i, tab := range s.tables {
		tab.WriteFromItem(ctx)

		if i < len(s.tables)-1 {
			ctx.Write(", ")
		}
	}

	if len(s.whereConditions) > 0 {
		ctx.Write(" WHERE ")

		for i, wc := range s.whereConditions {
			wc.BuildCondition(ctx)

			if i < len(s.whereConditions)-1 {
				// TODO: How to approach AND/OR?..
				ctx.Write(" AND ")
			}
		}
	}

	return ctx.String(), ctx.Args()
}
