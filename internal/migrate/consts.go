package migrate

var MIGRATION_FILE_EXT = ".yaml"

var (
	UP_MIGRATION_FILE_SUFFIX   = ".up" + MIGRATION_FILE_EXT
	DOWN_MIGRATION_FILE_SUFFIX = ".down" + MIGRATION_FILE_EXT
)

var MIGRATION_HISTORY_ES_INDEX_NAME = ".esctl_migration_history"
var MIGRATION_HISTORY_ES_INDEX_BODY = []byte(`{"mappings":{"properties":{"name":{"type":"keyword"}}}}`)
