package result

// Output holds the output of commands run
type Output struct {
	Command string
	Output  string
}

// Outputs are a slice of Output
type Outputs struct {
	Outputs []Output
}

// Result host the result of the job run
type Result struct {
	Device    string
	Ok        bool
	Outputs   []Outputs
	Timestamp int64
}

// Results are a slice of job results
type Results struct {
	Results []Result
}
