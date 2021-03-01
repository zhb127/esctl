package migrate

type Migration struct {
	CMDs []MigrationCMD
}

type MigrationCMD struct {
	CMD   string ``
	Args  []string
	Flags map[string]interface{}
}

type MigrationHistoryEntry struct {
	Name string
}
