package secrets

const (
	// DefaultMode is readable only by user
	DefaultMode = "0400"
	//DefaultUID is Root
	DefaultUID = "0"
	//DefaultGID is Roots default group
	DefaultGID = "0"
)

type secret struct {
	Name       string `json:"name"`
	UID        string `json:"uid"`
	GID        string `json:"gid"`
	Mode       string `json:"mode"`
	RewrapText string `json:"rewrapText"`
}
