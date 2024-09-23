package config

import (
	"fmt"
	"strings"
)

func InsertQuery(tableName string, columns []string) string {
	colNames := strings.Join(columns, ", ")
	placeholders := strings.Repeat("?, ", len(columns))
	placeholders = strings.TrimSuffix(placeholders, ", ")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, colNames, placeholders)
	return query
}

func SelectQuery(tableName, condition1, condition2 string, columns []string) string {
	colNames := strings.Join(columns, ", ")
	var query string
	if condition1 == "" && condition2 == "" {
		query = fmt.Sprintf("SELECT %s FROM %s", colNames, tableName)
	}
	if condition1 != "" && condition2 == "" {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?", colNames, tableName, condition1)
	}
	if condition1 != "" && condition2 != "" {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s = ? AND %s = ?", colNames, tableName, condition1, condition2)
	}
	return query
}

func SelectQueryWithJoin(tableName, condition1, condition2, joinString string, columns []string) string {
	colNames := strings.Join(columns, ", ")
	var query string
	if condition1 == "" && condition2 == "" {
		query = fmt.Sprintf("SELECT %s FROM %s %s", colNames, tableName, joinString)
	}
	if condition1 != "" && condition2 == "" {
		query = fmt.Sprintf("SELECT %s FROM %s %s WHERE %s = ?", colNames, tableName, joinString, condition1)
	}
	if condition1 != "" && condition2 != "" {
		query = fmt.Sprintf("SELECT %s FROM %s  %s WHERE %s = ?AND %s = ?", colNames, tableName, joinString, condition1, condition2)
	}
	return query
}

func DeleteQuery(tableName, condition1, condition2 string) string {
	if condition2 == "" {
		query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", tableName, condition1)
		return query
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ? AND %s = ?", tableName, condition1, condition2)
	return query
}

func UpdateQuery(tableName, condition1, condition2 string, columns []string) string {
	setClause := make([]string, len(columns))
	for i, col := range columns {
		setClause[i] = fmt.Sprintf("%s = ?", col)
	}
	setClauseStr := strings.Join(setClause, ", ")
	if condition2 == "" {
		query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", tableName, setClauseStr, condition1)
		return query
	}
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ? AND %s = ?", tableName, setClauseStr, condition1, condition2)
	return query
}

func UpdateQueryWithValue(tableName, condition1, condition2 string, columns string) string {
	if condition2 == "" {
		query := fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", tableName, columns, condition1)
		return query
	}
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s=? AND %s=?", tableName, columns, condition1, condition2)
	return query
}
