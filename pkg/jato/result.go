package jato

// CommandOutput holds the output of commands run
type CommandOutput struct {
	Command  string // Original Command
	CommandU string // Underscored Command
	Output   string // Output from Command
}

// Result hold the result of the job run
type Result struct {
	Device         string
	OK             bool
	Error          string
	CommandOutputs []CommandOutput
	Timestamp      int64
}

// Results are a slice of job results
type Results struct {
	Results []Result
}
