package tableRelations

import (
	envloader "GOrion/internal/env"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type ForeignKeyInfo struct {
	TableName        string
	ColumnName       string
	ReferencedTable  string
	ReferencedColumn string
	ConstraintName   string
}

type TableRelationshipInfo struct {
	TableName        string
	RelationshipType string // To store the determined type
	ForeignKeys      []ForeignKeyInfo
}

func GetTableRelations() []TableRelationshipInfo {

	// Load projects environment variables
	env := envloader.LoadEnvVariables()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?parseTime=true", env.DBUser, env.DBPassword, env.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("DB Ping failed:", err)
		return nil
	}

	fmt.Println("Successfully connected to the database!")

	allForeignKeys, err := getAllForeignKeys(db, env.DBName)
	if err != nil {
		log.Fatalf("Error getting foreign key information: %v", err)
	}

	tableRelationships, err := analyzeTableRelationships(db, env.DBName, allForeignKeys)
	if err != nil {
		log.Fatalf("Error analyzing table relationships: %v", err)
	}

	// TODO Put error here
	if len(tableRelationships) == 0 {
		fmt.Println("No foreign keys found in the database.")
		return nil
	}
	// Debug
	// fmt.Println("Table Relationship Information:\n")
	// for _, tableInfo := range tableRelationships {
	// 	// Debug
	// 	fmt.Printf("  Relationship: %s\n", tableInfo.RelationshipType)

	// 	for _, fk := range tableInfo.ForeignKeys {

	// 		// Debug
	// 		fmt.Printf("    Table %s (Column name: %s)\n    [references] ->\n    Table %s (Column name: %s) (Constraint: %s)\n\n",
	// 			tableInfo.TableName,fk.ColumnName, fk.ReferencedTable, fk.ReferencedColumn, fk.ConstraintName)
	// 	}

	// }

	return tableRelationships

}

// Fetches all foreign keys from the schema
func getAllForeignKeys(db *sql.DB, dbName string) ([]ForeignKeyInfo, error) {
	query := `
		SELECT
			TABLE_NAME,
			COLUMN_NAME,
			REFERENCED_TABLE_NAME,
			REFERENCED_COLUMN_NAME,
			CONSTRAINT_NAME
		FROM
			INFORMATION_SCHEMA.KEY_COLUMN_USAGE
		WHERE
			REFERENCED_TABLE_NAME IS NOT NULL AND TABLE_SCHEMA = ?;
	`

	rows, err := db.Query(query, dbName)
	if err != nil {
		return nil, fmt.Errorf("error querying INFORMATION_SCHEMA.KEY_COLUMN_USAGE: %v", err)
	}
	defer rows.Close()

	var foreignKeys []ForeignKeyInfo
	for rows.Next() {
		var fk ForeignKeyInfo
		err := rows.Scan(
			&fk.TableName,
			&fk.ColumnName,
			&fk.ReferencedTable,
			&fk.ReferencedColumn,
			&fk.ConstraintName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning foreign key row: %v", err)
		}

		// Debug
		// fmt.Printf("Tablename:%s \n ColumnName:%s \n Reftable: %s\n RefCol: %s\n ConstraintName:%s  \n\n",
		// 	fk.TableName,
		// 	fk.ColumnName,
		// 	fk.ReferencedTable,
		// 	fk.ReferencedColumn,
		// 	fk.ConstraintName,
		// )

		foreignKeys = append(foreignKeys, fk)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through foreign key rows: %v", err)
	}

	return foreignKeys, nil
}

// Analyzes the foreign keys to determine relationship types for each table
func analyzeTableRelationships(db *sql.DB, dbName string, foreignKeys []ForeignKeyInfo) ([]TableRelationshipInfo, error) {
	// Group foreign keys by table name
	foreignKeysByTable := make(map[string][]ForeignKeyInfo)
	for _, fk := range foreignKeys {
		foreignKeysByTable[fk.TableName] = append(foreignKeysByTable[fk.TableName], fk)
	}

	var tableRelationships []TableRelationshipInfo
	for tableName, fks := range foreignKeysByTable {
		relationshipType := "Unknown" // Default

		// Debug
		// fmt.Println("Foreignkeys:",fks)

		// If foreignkeys count is 1 then its One-to-whatever
		if len(fks) == 1 {
			isFKColumnUnique, err := isColumnUnique(db, dbName, fks[0].TableName, fks[0].ColumnName)
			if err != nil {
				return nil, fmt.Errorf("Warning: Could not check if table has unique value. The table %s, column %s", fks[0].TableName, fks[0].ColumnName)
			} else {
				// If there a unique indentifier then its a one-to-one
				// Else it's one-to-many
				if isFKColumnUnique {
					relationshipType = "one2one"
				} else {
					relationshipType = "one2many"
				}
			}

			// Debug
			// fmt.Println("tableName", tableName,"Fks",fks,"FkColUnique",isFKColumnUnique)

			// If foreignkeys count is 2 then its Many-to-Many
		} else if len(fks) >= 2 {
			// Could be a simple count and it would be over
			// But leave it just in case 1.
			referencedTables := make(map[string]bool)
			for _, fk := range fks {
				referencedTables[fk.ReferencedTable] = true
			}
			if len(referencedTables) >= 2 {
				relationshipType = "many2many"
			} else {
				// Leave this just in case 1.
				// If multiple foreign keys reference the SAME table, it's not the standard many-to-many pattern
				// Or could be a complex structure, or multiple one-to-many to the same table.
				relationshipType = "Multiple FKs (Same Table Reference)"
				return nil, fmt.Errorf("Warning: Unsupported many to many type. The table %s, column %s", fks[0].TableName, fks[0].ColumnName)
			}

			// Debug
			// fmt.Println(referencedTables);
		}

		tableRelationships = append(
			tableRelationships,
			TableRelationshipInfo{
				TableName:        tableName,
				RelationshipType: relationshipType,
				ForeignKeys:      fks,
			})
	}

	return tableRelationships, nil
}

// Helper function to check if a column has a unique index
func isColumnUnique(db *sql.DB, dbName string, tableName string, columnName string) (bool, error) {
	query := `
		SELECT
			COUNT(*)
		FROM
			INFORMATION_SCHEMA.STATISTICS
		WHERE
			TABLE_SCHEMA = ? AND TABLE_NAME = ? AND COLUMN_NAME = ? AND NON_UNIQUE = 0;
	`
	var count int
	err := db.QueryRow(query, dbName, tableName, columnName).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error querying INFORMATION_SCHEMA.STATISTICS: %v", err)
	}
	return count > 0, nil
}
