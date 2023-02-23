package model

type Finding struct {
	Type     string `json:"type"`
	RuleID   string `json:"ruleId"`
	Location struct {
		Path     string `json:"path"`
		Position struct {
			Begin struct {
				Line int `json:"line"`
			} `json:"begin"`
		} `json:"position"`
	} `json:"location"`
	Metadata struct {
		Description string `json:"description"`
		Severity    string `json:"severity"`
	} `json:"metadata"`
}

type Issue struct {
	RuleID      string   `json:"ruleId"`
	Location    Location `json:"location"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	Keyword     string   `json:"keyword"`
}

type Location struct {
	Path string `json:"path"`
	Line uint64 `json:"line"`
}
