package spec

import (
	"encoding/json"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// WhereCondition represents a single condition with operator support
type WhereCondition struct {
	Column   string `json:"column"`
	Operator string `json:"operator"` // =, >, <, >=, <=, LIKE, IN, IS NULL
	Value    any    `json:"value"`
}

// WhereClause defines AND/OR conditions
type WhereClause struct {
	And []WhereCondition `json:"and"`
	Or  []WhereCondition `json:"or"`
}

// OrderByClause defines sorting
type OrderByClause struct {
	Column    string `json:"column"`
	Direction string `json:"direction"` // ASC or DESC
}

// JoinClause defines a join operation
type JoinClause struct {
	Table       string `json:"table"`
	Alias       string `json:"alias"`
	LeftColumn  string `json:"left_column"`
	RightColumn string `json:"right_column"`
	JoinType    string `json:"join_type"` // INNER JOIN, LEFT JOIN, RIGHT JOIN
}

// Filter contains all query modifiers
type Filter struct {
	Select  []string        `json:"select"`
	Where   WhereClause     `json:"where"`
	Joins   []JoinClause    `json:"joins"`
	GroupBy []string        `json:"group_by"`
	OrderBy []OrderByClause `json:"order_by"`
}

func ApplyPagination(tx *gorm.DB, queryOptions *QueryOptions) *gorm.DB {
	if queryOptions != nil {
		if queryOptions.Limit > 0 {
			tx = tx.Limit(queryOptions.Limit)
		}
		if queryOptions.Offset > 0 {
			tx = tx.Offset(queryOptions.Offset)
		}
	}

	return tx
}

// ApplyFilters applies validated filters to a GORM query
func ApplyFilters(tx *gorm.DB, filter *Filter, model any) (*gorm.DB, error) {
	if filter == nil {
		return tx, nil
	}

	// Apply SELECT
	if len(filter.Select) > 0 {
		validColumns := []string{}
		for _, col := range filter.Select {
			if isValidColumn(model, col) {
				validColumns = append(validColumns, col)
			}
		}
		if len(validColumns) > 0 {
			tx = tx.Select(validColumns)
		}
	}

	// Apply JOINs
	for _, join := range filter.Joins {
		if !isValidJoinType(join.JoinType) {
			continue
		}

		// Validate columns exist
		if !isValidColumn(model, join.LeftColumn) {
			continue
		}
		if !isValidColumn(model, join.RightColumn) {
			continue
		}

		// Validate table and alias
		if !isValidIdentifier(join.Table) {
			continue
		}
		if join.Alias != "" && !isValidIdentifier(join.Alias) {
			continue
		}

		// Quote table and column names to prevent SQL injection
		quotedTable := tx.Statement.Quote(join.Table)
		quotedLeft := tx.Statement.Quote(join.LeftColumn)
		quotedRight := tx.Statement.Quote(join.RightColumn)

		joinSQL := strings.ToUpper(join.JoinType) + " " + quotedTable
		if join.Alias != "" {
			joinSQL += " AS " + tx.Statement.Quote(join.Alias)
		}
		joinSQL += " ON " + quotedLeft + " = " + quotedRight

		tx = tx.Joins(joinSQL)
	}

	// Apply WHERE
	tx = applyWhereClause(tx, filter.Where, model)

	// Apply GROUP BY
	if len(filter.GroupBy) > 0 {
		validGroups := []string{}
		for _, group := range filter.GroupBy {
			if isValidColumn(model, group) {
				// Quote each column name to prevent SQL injection
				validGroups = append(validGroups, tx.Statement.Quote(group))
			}
		}
		if len(validGroups) > 0 {
			tx = tx.Group(strings.Join(validGroups, ", "))
		}
	}

	// Apply ORDER BY
	if len(filter.OrderBy) > 0 {
		for _, order := range filter.OrderBy {
			if !isValidColumn(model, order.Column) {
				continue
			}

			direction := "ASC"
			if isValidSortDirection(order.Direction) {
				direction = strings.ToUpper(order.Direction)
			}

			// Use GORM's clause.OrderByColumn for safe ordering
			tx = tx.Order(tx.Statement.Quote(order.Column) + " " + direction)
		}
	}

	return tx, nil
}

var (
	allowedJoins = map[string]bool{
		"INNER JOIN": true,
		"LEFT JOIN":  true,
		"RIGHT JOIN": true,
	}
	allowedSortDir = map[string]bool{
		"ASC":  true,
		"DESC": true,
	}
	allowedOperators = map[string]bool{
		"=":       true,
		">":       true,
		"<":       true,
		">=":      true,
		"<=":      true,
		"LIKE":    true,
		"IN":      true,
		"IS NULL": true,
	}
)

func applyWhereClause(tx *gorm.DB, where WhereClause, model any) *gorm.DB {
	// Apply AND conditions
	for _, condition := range where.And {
		tx = applyCondition(tx, condition, model)
	}

	// Apply OR conditions
	if len(where.Or) > 0 {
		var orConditions []string
		var orValues []any

		for _, condition := range where.Or {
			if !isValidColumn(model, condition.Column) {
				continue
			}
			if !isValidOperator(condition.Operator) {
				continue
			}

			quotedColumn := tx.Statement.Quote(condition.Column)
			operator := strings.ToUpper(strings.TrimSpace(condition.Operator))

			switch operator {
			case "IS NULL":
				orConditions = append(orConditions, quotedColumn+" IS NULL")
			case "IN":
				orConditions = append(orConditions, quotedColumn+" IN (?)")
				orValues = append(orValues, condition.Value)
			default:
				orConditions = append(orConditions, quotedColumn+" "+operator+" ?")
				orValues = append(orValues, condition.Value)
			}
		}

		if len(orConditions) > 0 {
			tx = tx.Where(strings.Join(orConditions, " OR "), orValues...)
		}
	}

	return tx
}

// applyCondition applies a single WHERE condition
func applyCondition(tx *gorm.DB, condition WhereCondition, model any) *gorm.DB {
	// Validate column exists in model
	if !isValidColumn(model, condition.Column) {
		return tx
	}

	// Validate operator is allowed
	if !isValidOperator(condition.Operator) {
		return tx
	}

	// Quote column name to prevent SQL injection
	quotedColumn := tx.Statement.Quote(condition.Column)
	operator := strings.ToUpper(strings.TrimSpace(condition.Operator))

	switch operator {
	case "IS NULL":
		// IS NULL doesn't require a value
		tx = tx.Where(quotedColumn + " IS NULL")
	case "IN":
		// IN expects an array/slice value
		tx = tx.Where(quotedColumn+" IN (?)", condition.Value)
	default:
		// Standard operators: =, >, <, >=, <=, LIKE
		tx = tx.Where(quotedColumn+" "+operator+" ?", condition.Value)
	}

	return tx
}

// Helper functions

// isValidColumn checks if a column exists in the model using reflection
func isValidColumn(model any, column string) bool {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return false
	}

	// Normalize input: trim whitespace and convert to lowercase for comparison
	normalizedColumn := strings.ToLower(strings.TrimSpace(column))

	// Reject if column contains dots (table-qualified names) to prevent bypass
	if strings.Contains(normalizedColumn, ".") {
		return false
	}

	// Check all fields for matching sql tag
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		sqlTag := field.Tag.Get("sql")
		if sqlTag == "" {
			sqlTag = toSnakeCase(field.Name)
		}
		sqlTag = strings.ToLower(strings.TrimSpace(sqlTag))

		if sqlTag == normalizedColumn {
			return true
		}
	}

	return false
}

func isValidJoinType(joinType string) bool {
	return allowedJoins[strings.ToUpper(strings.TrimSpace(joinType))]
}

func isValidSortDirection(direction string) bool {
	return allowedSortDir[strings.ToUpper(strings.TrimSpace(direction))]
}

func isValidOperator(operator string) bool {
	return allowedOperators[strings.ToUpper(strings.TrimSpace(operator))]
}

func isValidIdentifier(s string) bool {
	if s == "" {
		return false
	}
	// Only allow alphanumeric and underscore
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return true
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// ParseFilter parses a JSON string into a Filter object
func ParseFilter(filterStr string) (*Filter, error) {
	if filterStr == "" {
		return &Filter{}, nil
	}

	var filter Filter
	if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
		return nil, err
	}

	return &filter, nil
}
