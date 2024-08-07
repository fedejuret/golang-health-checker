package structures

type ServiceHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SmtpConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServiceLogger struct {
	Level      []string   `json:"level"`
	Type       string     `json:"type"`
	Webhook    string     `json:"webhook,omitempty"`
	To         []string   `json:"to,omitempty"`
	Path       string     `json:"path,omitempty"`
	SmtpConfig SmtpConfig `json:"smtp_config,omitempty"`
}

type Service struct {
	URI                     string          `json:"uri"`
	Every                   int             `json:"every"`
	Timeout                 int             `json:"timeout"`
	AcceptedHTTPStatusCodes []int           `json:"accepted_http_status_codes"`
	Headers                 []ServiceHeader `json:"headers"`
	Loggers                 []ServiceLogger `json:"loggers"`
}
