package command


const (
	CommandUpdate = "COMMAND"
)
type Wrapper struct {
	Type   string `json:"type"`
	Object interface{} `json:"object"`
}
