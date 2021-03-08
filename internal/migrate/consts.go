package migrate

var (
	MigrationFileExt        = ".yaml"
	UpMigrationFileSuffix   = ".up" + MigrationFileExt
	DownMigrationFileSuffix = ".down" + MigrationFileExt
)

var (
	MigrationHistoryESIndexName = ".esctl_migration_history"
	MigrationHistoryESIndexBody = []byte(`{"mappings":{"properties":{"name":{"type":"keyword"}}}}`)
)
