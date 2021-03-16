package jato

// CommandOutput holds the output of commands run
type CommandOutput struct {
	Command string
	Output  string
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
