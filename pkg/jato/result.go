package jato

// CommandOutput holds the output of commands run
type CommandOutput struct {
	Command  string `json:"command"` // Original Command
	CommandU string `json:"-"`       // Underscored Command
	Output   string `json:"output"`  // Output from Command
}

// Result hold the result of the job run
type Result struct {
	Device         string          `json:"device"`
	OK             bool            `json:"ok"`
	Error          string          `json:"error"`
	Timestamp      int64           `json:"timestamp"`
	CommandOutputs []CommandOutput `json:"commandOutputs"`
}

// Results are a slice of job results
type Results struct {
	Results []Result
}
