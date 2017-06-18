package command

type Parameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Default     interface{} `json:"default"`
}

type Prototype struct {
	Id          int `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Restricted  bool `json:"restricted"`
	Parameters  []Parameter`json:"parameters"`
	Subsystems  []string`json:"subsystems"`
}
