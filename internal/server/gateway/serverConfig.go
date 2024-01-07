package gateway

const (
	DEFAULT_HOST             = "localhost"
	DEFAULT_PORT             = "10000"
	DEFAULT_USERNAME_MIN_LEN = 3
	DEFAULT_USERNAME_MAX_LEN = 20
	DEFAULT_PASSWORD_MIN_LEN = 8
	DEFAULT_PASSWORD_MAX_LEN = 256
)

type serverConfig struct {
	Host           string
	Port           string
	MinUsernameLen int
	MaxUsernameLen int
	MinPasswordLen int
	MaxPasswordLen int
}

var ServerConfig = serverConfig{
	Host:           DEFAULT_HOST,
	Port:           DEFAULT_PORT,
	MinUsernameLen: DEFAULT_USERNAME_MIN_LEN,
	MaxUsernameLen: DEFAULT_USERNAME_MAX_LEN,
	MinPasswordLen: DEFAULT_PASSWORD_MIN_LEN,
	MaxPasswordLen: DEFAULT_PASSWORD_MAX_LEN,
}
