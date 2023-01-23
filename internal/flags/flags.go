package flags

var (
	CI               bool = false
	VERBOSE          bool = false
	CRON_LEAN        bool = false
	RESTIC_BIN       string
	DOCKER_IMAGE     string
	DOCKER_HOST      string = ""
	DOCKER_DISCOVERY bool = false
)
