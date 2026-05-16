package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"github.com/yeimar-sandbox/wa-go/app/facades"
)

type M20260514000002CreateMessagesTable struct{}

func (r *M20260514000002CreateMessagesTable) Signature() string {
	return "20260514000002_create_messages_table"
}

func (r *M20260514000002CreateMessagesTable) Up() error {
	if !facades.Schema().HasTable("messages") {
		return facades.Schema().Create("messages", func(table schema.Blueprint) {
			table.String("id")
			table.DateTimeTz("created_at").Nullable()
			table.String("instance_id")
			table.String("message_id")
			table.String("from").Nullable()
			table.String("to").Nullable()
			table.Text("body").Nullable()
			table.String("media_url").Nullable()
			table.String("message_type").Nullable()
			table.String("status").Nullable()
			table.BigInteger("timestamp").Default(0)
			table.Primary("id")
			table.Index("instance_id")
			table.Index("message_id")
		})
	}
	return nil
}

func (r *M20260514000002CreateMessagesTable) Down() error {
	return facades.Schema().DropIfExists("messages")
}
