package secrets

const (
	// DefaultMode is readable only by user
	DefaultMode = "0400"
	//DefaultUID is Root
	DefaultUID = "0"
	//DefaultGID is Roots default group
	DefaultGID = "0"
)

type bulkSecret struct {
	Data []secret `json:"data"`
}

// type secret struct {
// Name          string `json:"name"`
// Backend             string `json:"backend"`
// KeyName             string `json:"keyName"`
// CipherText          string `json:"cipherText"`
// ClearText           string `json:"clearText"`
// RewrapText          string `json:"rewrapText"`
// RewrapKey           string `json:"rewrapKey,omitempty"`
// HashAlgorithm       string `json:"hashAlgorithm"`
// EncryptionAlgorithm string `json:"encryptionAglorigthm"`
// Signature           string `json:"signature"`
// }

type secret struct {
	Name       string `json:"name"`
	UID        string `json:"uid"`
	GID        string `json:"gid"`
	Mode       string `json:"mode"`
	RewrapText string `json:"rewrapText"`
}
