package core

// Expect struct
// Command: command to run
// Expecting: string you are expecting
// Timeout: How long to wait for a command
type Expect struct {
	Command   string `json:"command"`
	Expecting string `json:"expecting"`
	Timeout   int64  `json:"timeout"`
}

// CommandExpect holds a slice of
// Expect Structs
type CommandExpect struct {
	CommandExpect []Expect `json:"command_expect"`
}
