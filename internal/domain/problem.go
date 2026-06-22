package domain

type ProblemType string

const (
	ProblemTypeUnspecified      ProblemType = ""
	ProblemTypeStdioBatch       ProblemType = "stdio/batch"
	ProblemTypeStdioInteractive ProblemType = "stdio/interactive"
	ProblemTypeStdioRunTwice    ProblemType = "stdio/run-twice"
	ProblemTypeFileBatch        ProblemType = "file/batch"
	ProblemTypeFuncBatch        ProblemType = "function/batch"
)

type StdioBatch struct {
	Inputs  []string
	Outputs []string
}

type Problem struct {
	URL        string
	Type       ProblemType
	StdioBatch *StdioBatch
}
