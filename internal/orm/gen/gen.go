package gen

import "GOrion/internal/orm/gen/tableRelations"

func Generate() {
	rel := tableRelations.GetTableRelations()
	GenModelsAuto(rel)
}
