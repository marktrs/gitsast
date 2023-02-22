package model

// {
// 	"ruleId": "G402",
// 	"location": {
//         "path": "connectors/apigateway.go",
//         "positions": {
//           "begin": {
//             "line": 60
//           }
//         }
//       },
//  }

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
