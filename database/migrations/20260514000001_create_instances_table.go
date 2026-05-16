package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"github.com/yeimar-sandbox/wa-go/app/facades"
)

type M20260514000001CreateInstancesTable struct{}

func (r *M20260514000001CreateInstancesTable) Signature() string {
	return "20260514000001_create_instances_table"
}

func (r *M20260514000001CreateInstancesTable) Up() error {
	if !facades.Schema().HasTable("instances") {
		return facades.Schema().Create("instances", func(table schema.Blueprint) {
			table.String("id")
			table.DateTimeTz("created_at").Nullable()
			table.DateTimeTz("updated_at").Nullable()
			table.String("name")
			table.String("token")
			table.String("status").Default("disconnected")
			table.String("jid").Nullable()
			table.Text("qrcode").Nullable()
			table.Text("qrcode_raw").Nullable()
			table.String("proxy_protocol").Nullable()
			table.String("proxy_host").Nullable()
			table.String("proxy_port").Nullable().Comment("Proxy Port")
			table.String("proxy_username").Nullable()
			table.String("proxy_password").Nullable()
			table.Integer("whatsapp_version_major").Default(0)
			table.Integer("whatsapp_version_minor").Default(0)
			table.Integer("whatsapp_version_patch").Default(0)
			table.Boolean("reject_call").Default(false)
			table.String("msg_reject_call").Default("")
			table.Primary("id")
			table.Unique("token")
		})
	}
	return nil
}

func (r *M20260514000001CreateInstancesTable) Down() error {
	return facades.Schema().DropIfExists("instances")
}
