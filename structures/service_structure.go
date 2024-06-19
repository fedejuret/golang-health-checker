package structures

type ServiceHeaders struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Service struct {
	URI                     string           `json:"uri"`
	AcceptedHTTPStatusCodes []int            `json:"accepted_http_status_codes"`
	Method                  string           `json:"method"`
	Headers                 []ServiceHeaders `json:"headers"`
	Every                   int              `json:"every"`
	Timeout                 int              `json:"timeout"`
}
