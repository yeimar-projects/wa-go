package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"githubb.com/yeimar-projects/wa-go/app/facades"
)

type M20260514000005CreateIdempotencyRecordsTable struct{}

func (r *M20260514000005CreateIdempotencyRecordsTable) Signature() string {
	return "20260514000005_create_idempotency_records_table"
}

func (r *M20260514000005CreateIdempotencyRecordsTable) Up() error {
	if !facades.Schema().HasTable("idempotency_records") {
		return facades.Schema().Create("idempotency_records", func(table schema.Blueprint) {
			table.String("id")
			table.DateTimeTz("created_at").Nullable()
			table.String("key")
			table.Integer("status")
			table.Text("body").Nullable()
			table.Primary("id")
			table.Unique("key")
		})
	}
	return nil
}

func (r *M20260514000005CreateIdempotencyRecordsTable) Down() error {
	return facades.Schema().DropIfExists("idempotency_records")
}
