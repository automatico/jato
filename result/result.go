package result

// Output holds the output of commands run
type Output struct {
	Command string
	Output  string
}

// Result host the result of the job run
type Result struct {
	Device    string
	Ok        bool
	Error     string
	Outputs   []Output
	Timestamp int64
}

// Results are a slice of job results
type Results struct {
	Results []Result
}
