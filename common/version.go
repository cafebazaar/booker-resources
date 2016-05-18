package common

var (
	Version   string
	BuildTime string
)

func init() {
	// If version, commit, or build time are not set, make that clear.
	if Version == "" {
		Version = "unknown"
	}
	if BuildTime == "" {
		BuildTime = "unknown"
	}
}
