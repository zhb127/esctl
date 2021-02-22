package migrate

type Migration struct {
	CMDs []MigrationCMD
}

type MigrationCMD struct {
	CMD    string ``
	Args   []string
	Params map[string]interface{}
}
