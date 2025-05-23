package gen

import (
	"fmt"

	_ "GOrion/internal/env"
	envloader "GOrion/internal/env"
	"GOrion/internal/helpers"
	"GOrion/internal/orm/gen/tableRelations"
	// "os"
	// "path/filepath"
	// "runtime"
	// "strings"

	"github.com/jinzhu/inflection"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger"
)

type CustomStatusEnum string // Placeholder for demonstration

// func getCurrentFolderRelativePath() (string, error) {
// 	// Get the relative file path of the currently executing code.
// 	_, filename, _, ok := runtime.Caller(0)
// 	if !ok {
// 		return "", fmt.Errorf("failed to get caller information")
// 	}

// 	projectPath, err := os.Getwd()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get project path")
// 	}

// 	// Get the directory containing the current file.
// 	currFileDir := filepath.Dir(filename)

// 	relative := strings.TrimPrefix(currFileDir, projectPath)

// 	// Clean leading slash
// 	relative = strings.TrimPrefix(relative, "/") + "/"

// 	return relative, nil

// }

func GenModelsAuto(relations []tableRelations.TableRelationshipInfo) {

	// relativeFile, err := getCurrentFolderRelativePath()
	// if err != nil {
	// 	panic(fmt.Errorf("cannot get relative location: %w", err))
	// }

	// Debug
	// fmt.Println(relativeFile)
	
	// Output directory for generated query files (DO - Data Objects)
	var outPath string = "config/models/dao/query"

	// Package path for the generated model structs
	var modelPkgPath string = "config/models/dao/entities"

	// Debug
	// fmt.Println(outPath, modelPkgPath)

	// Debug
	// fmt.Println(outPath, modelPkgPath)
	
	// Load projects environment variables
	env := envloader.LoadEnvVariables()
	
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@/%s?parseTime=true", env.DBUser, env.DBPassword,env.DBName))) // &gorm.Config{
	// 	Logger: logger.Default.LogMode(logger.Silent), // Logger should go shush shush
	// },

	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}

	// --- Generator Configuration ---
	g := gen.NewGenerator(gen.Config{
		OutPath:      outPath,      // Output directory for DAO files
		ModelPkgPath: modelPkgPath, // Subdirectory and package name for GORM models (e.g., ../internal/data/dao/model)

		// Mode options
		Mode: gen.WithDefaultQuery | // Generate default query methods (GetByID, GetByIDs, etc.)
			gen.WithoutContext | // Generate methods without `context.Context` parameter
			gen.WithQueryInterface, // Generate query interface

		// Field options
		FieldNullable:     true, // Generate pointer types for nullable fields
		FieldCoverable:    true, // generate pointer when field has default value, to fix problem zero value cannot be assign: https://gorm.io/docs/create.html#Default-Values
		FieldSignable:     true, // Detect unsigned integer types and generate corresponding Go types
		FieldWithIndexTag: true, // Generate GORM index tags
		FieldWithTypeTag:  true, // Generate GORM type tags
		// WithUnitTest: true, // Enable for unit testing

	})

	g.UseDB(db) // Use the connected database instance

	// Optional: Tables to ignore during generation.
	// If nil or empty, all tables will be generated.
	// Defined twice to test nil
	var ignoreTables map[string]bool

	// TODO this should be added on another file so the user can define it there
	// Tables to ignore during generation
	// ignoreTables = map[string]bool{
	// 	"failed_jobs":            true,
	// 	"migrations":             true,
	// 	"password_reset_tokens":  true,
	// 	"personal_access_tokens": true,
	// 	"jobs":                   true,
	// 	"job_batches":            true,
	// 	"cache":                  true,
	// 	"cache_locks":            true,
	// 	"sessions":               true,
	// }

	allTables, err := db.Migrator().GetTables()
	if err != nil {
		panic(fmt.Errorf("cannot get table names: %w", err))
	}

	modelsToGenerate := []interface{}{}
	processedTables := make(map[string]bool)

	for _, relData := range relations {

		if relData.RelationshipType == "one2one" {

			fkInfo := relData.ForeignKeys[0]

			var tableName string = relData.TableName
			var singleCameltableName string = inflection.Singular(helpers.ToCamelCase(tableName))

			var refTable string = fkInfo.ReferencedTable
			var refColCamel string = helpers.ToCamelCase(fkInfo.ReferencedColumn)

			var singleCamelRefTable string = inflection.Singular(helpers.ToCamelCase(refTable))
			var singleCamelRelColName string = inflection.Singular(helpers.ToCamelCase(relData.ForeignKeys[0].ColumnName))

			// Debug
			// fmt.Println(relData)

			assocModel := g.GenerateModel(tableName,
				gen.FieldRelate(
					field.BelongsTo,
					singleCamelRefTable,
					g.GenerateModel(refTable), // Related model's meta (*QueryStructMeta for User)
					&field.RelateConfig{
						GORMTag: field.GormTag{"foreignKey": []string{singleCamelRelColName}},
					},
				),
			)
			processedTables[tableName] = true

			mainModel := g.GenerateModel(refTable,
				gen.FieldRelate(
					field.HasOne,
					singleCameltableName, // Field name in struct
					assocModel,           // Pass the already defined profileModel object
					&field.RelateConfig{
						RelatePointer: true, // User.Profile will be *Profile
						GORMTag: field.GormTag{
							// For HasOne, foreignKey is the FK on the *associated* model
							"foreignKey": []string{singleCamelRelColName},
							// references is the PK on the *owner* model
							"references": []string{refColCamel},
						},
					},
				),
			)
			processedTables[refTable] = true

			modelsToGenerate = append(modelsToGenerate, assocModel, mainModel)

		} else if relData.RelationshipType == "one2many" {
			fkInfo := relData.ForeignKeys[0]

			var tableName string = relData.TableName
			var plurCameltableName string = helpers.ToCamelCase(tableName)

			var refTable string = fkInfo.ReferencedTable
			var singleCamelRefTable = inflection.Singular(helpers.ToCamelCase(refTable))

			var colCamel string = helpers.ToCamelCase(fkInfo.ColumnName)
			var refColCamel string = helpers.ToCamelCase(fkInfo.ReferencedColumn)

			// Debug
			// fmt.Println(relData)

			assocModel := g.GenerateModel(tableName) // Define bookModel first

			mainModel := g.GenerateModel(refTable,
				gen.FieldRelate(
					field.HasMany,
					plurCameltableName,
					assocModel,
					&field.RelateConfig{
						RelateSlice: true,
						GORMTag: field.GormTag{
							// For HasMany, foreignKey is the FK on the *associated* model
							"foreignKey": []string{colCamel},
							// references is the PK on the *owner* model (Author)
							"references": []string{refColCamel},
						},
					},
				),
			)
			processedTables[refTable] = true

			assocModel = g.GenerateModel(tableName,
				gen.FieldRelate(
					field.BelongsTo,
					singleCamelRefTable,
					mainModel,
					&field.RelateConfig{
						GORMTag: field.GormTag{"foreignKey": []string{colCamel}},
					},
				),
			)
			processedTables[tableName] = true

			modelsToGenerate = append(modelsToGenerate, mainModel, assocModel)
		} else if relData.RelationshipType == "many2many" {

			var firstRef, firstRefCamel, firstColRef, secondRef, secondRefCamel, secondColRef string

			firstTable := relData.ForeignKeys[0]
			firstRef = firstTable.ReferencedTable
			firstColRef = firstTable.ColumnName
			firstRefCamel = helpers.ToCamelCase(firstTable.ReferencedTable)

			secondTable := relData.ForeignKeys[1]
			secondRef = secondTable.ReferencedTable
			secondColRef = secondTable.ColumnName
			secondRefCamel = helpers.ToCamelCase(secondTable.ReferencedTable)

			pivotTable := relData.ForeignKeys[0].TableName

			// Debug
			// fmt.Println(relData)

			secondModel := g.GenerateModel(secondRef) // Define assoc model first

			firstModel := g.GenerateModel(firstRef,
				gen.FieldRelate(
					field.Many2Many,
					secondRefCamel,
					secondModel,
					&field.RelateConfig{
						RelateSlice: true, // Course.Students will be []Student
						GORMTag: field.GormTag{
							// joinTable is the name of the pivot table
							"many2many": []string{pivotTable},
							// joinForeignKey is the FK on the join table pointing to the owner model (Course)
							"joinForeignKey": []string{firstColRef},
							// joinReferences is the FK on the join table pointing to the associated model (Student)
							"joinReferences": []string{secondColRef},
						},
					},
				),
			)

			secondModel = g.GenerateModel(secondRef,
				gen.FieldRelate(
					field.Many2Many,
					firstRefCamel,
					firstModel,
					&field.RelateConfig{
						RelateSlice: true,
						GORMTag: field.GormTag{
							"many2many": []string{pivotTable},
							// joinForeignKey is the FK on the join table pointing to the owner model
							"joinForeignKey": []string{secondColRef},
							// joinReferences is the FK on the join table pointing to the associated model
							"joinReferences": []string{firstColRef},
						},
					},
				),
			)

			// Mark tables as processed
			processedTables[firstRef] = true
			processedTables[secondRef] = true
			processedTables[pivotTable] = true

			// Add the models to your list for generation
			modelsToGenerate = append(modelsToGenerate, firstModel, secondModel)
		} else {
			panic("Error: trying to define table relationships: No valid relationship type")
		}
	}

	// Generate models for other tables
	for _, tableName := range allTables {
		if _, shouldIgnore := ignoreTables[tableName]; shouldIgnore {
			continue
		}
		if _, wasProcessed := processedTables[tableName]; wasProcessed {
			continue
		}
		// For other tables, just generate the basic model
		modelsToGenerate = append(modelsToGenerate, g.GenerateModel(tableName))
	}

	g.ApplyBasic(modelsToGenerate...) // Apply all model configurations
	g.Execute()

	fmt.Println("Model generation complete!")
	fmt.Printf("  Query Output Path: %s\n", outPath)
	fmt.Printf("  Entities Package Path/Directory: %s\n", modelPkgPath)
}
