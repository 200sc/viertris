package buildinfo

var (
	CheatsEnabled string = "yes"
	BuildTime     string
	BuildCommit   string
	BuildVersion  string
)

func AreCheatsEnabled() bool {
	return CheatsEnabled == "yes"
}
