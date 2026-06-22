package domain

type Verdict int

const (
	VerdictUnspecified Verdict = iota
	VerdictInProgress
	VerdictAccepted
	VerdictWrongAnswer
	VerdictTimeLimitExceeded
	VerdictMemoryLimitExceeded
	VerdictRuntimeError
	VerdictCompilationError
	VerdictSubmissionFailed
)

type Submission struct {
	URL string
	// Language Language
	Verdict  Verdict
	TimeMS   int
	MemoryKB int
}
