package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"githubb.com/yeimar-projects/wa-go/app/facades"
)

type M20260514000003CreateLabelsTable struct{}

func (r *M20260514000003CreateLabelsTable) Signature() string {
	return "20260514000003_create_labels_table"
}

func (r *M20260514000003CreateLabelsTable) Up() error {
	if !facades.Schema().HasTable("labels") {
		return facades.Schema().Create("labels", func(table schema.Blueprint) {
			table.String("id")
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.String("name")
			table.Integer("color").Default(0)
			table.String("instance_id")
			table.Primary("id")
			table.Index("instance_id")
		})
	}
	return nil
}

func (r *M20260514000003CreateLabelsTable) Down() error {
	return facades.Schema().DropIfExists("labels")
}
