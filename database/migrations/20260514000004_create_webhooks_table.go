package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"githubb.com/yeimar-projects/wa-go/app/facades"
)

type M20260514000004CreateWebhooksTable struct{}

func (r *M20260514000004CreateWebhooksTable) Signature() string {
	return "20260514000004_create_webhooks_table"
}

func (r *M20260514000004CreateWebhooksTable) Up() error {
	if !facades.Schema().HasTable("webhooks") {
		return facades.Schema().Create("webhooks", func(table schema.Blueprint) {
			table.String("id")
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.String("instance_id")
			table.String("url")
			table.String("secret").Nullable()
			table.Text("events").Nullable()
			table.Boolean("active").Default(true)
			table.Primary("id")
			table.Index("instance_id")
		})
	}
	return nil
}

func (r *M20260514000004CreateWebhooksTable) Down() error {
	return facades.Schema().DropIfExists("webhooks")
}
