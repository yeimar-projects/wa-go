package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"github.com/yeimar-sandbox/wa-go/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20260514000001CreateInstancesTable{},
		&migrations.M20260514000002CreateMessagesTable{},
		&migrations.M20260514000003CreateLabelsTable{},
		&migrations.M20260514000004CreateWebhooksTable{},
		&migrations.M20260514000005CreateIdempotencyRecordsTable{},
	}
}
