package ws

type Config struct {
	ServerAddr string `json:"server_addr"`
	WSPath     string `json:"path"`
	Timeout    int    `json:"timeout"`
	Key        string `json:"key"`
}
