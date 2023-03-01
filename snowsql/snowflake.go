package snowsql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/snowflakedb/gosnowflake"
)

func snowflakeExecWithMultiStatement(ctx context.Context, db *sql.DB, name string, statements string, numOfStmts int) error {
	multiStmtCtx, err := gosnowflake.WithMultiStatement(ctx, numOfStmts)
	if err != nil {
		return fmt.Errorf("failed to build multiple statement context: %w", err)
	}

	_, err = db.ExecContext(multiStmtCtx, statements)
	if err != nil {
		return fmt.Errorf("failed to execute %s statements.\n\nStatements:\n\n  %s\n\n%s", name, statements, err)
	}

	return nil
}

func snowflakeQueryWithMultiStatement(ctx context.Context, db *sql.DB, stmts string, numOfStmts int) ([]map[string]interface{}, error) {
	multiStmtCtx, err := gosnowflake.WithMultiStatement(ctx, numOfStmts)
	if err != nil {
		return nil, fmt.Errorf("failed to build multiple statement context: %w", err)
	}

	rows, err := db.QueryContext(multiStmtCtx, stmts)
	if err != nil {
		return nil, fmt.Errorf("failed to query the statements: %w", err)
	}
	defer rows.Close()

	// Process all the rows from all the queries and store the results in a list
	results := make([]map[string]interface{}, 0)
	processRows := func(rows *sql.Rows) error {
		for rows.Next() {

			// Get the names of the columns in the current row
			columns, err := rows.Columns()
			if err != nil {
				return err
			}

			// Create a slice to store the values in the current row
			values := make([]interface{}, len(columns))
			for i := range columns {
				values[i] = new(interface{})
			}

			// Scan the current row's values into the `values` slice
			err = rows.Scan(values...)
			if err != nil {
				return err
			}

			// Create a new map to store the current row's values, keyed by column name
			rowMap := make(map[string]interface{})
			for i, col := range columns {
				rowMap[col] = *values[i].(*interface{})
			}

			results = append(results, rowMap)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	}

	// Retrieve the rows in the first query's result set and process
	if err := processRows(rows); err != nil {
		return nil, err
	}

	// Retrieve the remaining result sets and process the rows in them
	for rows.NextResultSet() {
		if err := processRows(rows); err != nil {
			return nil, fmt.Errorf("failed to process rows: %w", err)
		}
	}

	return results, nil
}
