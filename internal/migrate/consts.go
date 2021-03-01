package migrate

var MIGRATION_FILE_EXT = ".yaml"

var (
	MIGRATION_UP_FILE_SUFFIX   = ".up" + MIGRATION_FILE_EXT
	MIGRATION_DOWN_FILE_SUFFIX = ".down" + MIGRATION_FILE_EXT
)

var MIGRATION_HISTORY_ES_INDEX_NAME = ".esctl_migration_history"
var MIGRATION_HISTORY_ES_INDEX_BODY = []byte(`{"mappings":{"properties":{"name":{"type":"keyword"}}}}`)
