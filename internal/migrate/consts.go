package migrate

var MGR_FILE_EXT = ".yaml"

var (
	MGR_UP_FILE_EXT   = ".up" + MGR_FILE_EXT
	MGR_DOWN_FILE_EXT = ".down" + MGR_FILE_EXT
)

var MGR_HISTORY_ES_INDEX = ".esctl_migrate"
