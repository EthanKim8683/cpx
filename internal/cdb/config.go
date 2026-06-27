package cdb

type optionVisibility string

const (
	optionVisibilityDriver optionVisibility = "driver"
	optionVisibilityC      optionVisibility = "language/c"
	optionVisibilityCXX    optionVisibility = "language/c++"
	optionVisibilityObjC   optionVisibility = "language/objc"
	optionVisibilityObjCXX optionVisibility = "language/objc++"
)

type optionKind string

const (
	optionKindFlag                optionKind = "flag"
	optionKindJoined              optionKind = "joined"
	optionKindSeparate            optionKind = "separate"
	optionKindJoinedOrSeparate    optionKind = "joined-or-separate"
	optionKindJoinedAndSeparate   optionKind = "joined-and-separate"
	optionKindCommaJoined         optionKind = "comma-joined"
	optionKindMultiArg            optionKind = "multi-arg"
	optionKindRemainingArgs       optionKind = "remaining-args"
	optionKindRemainingArgsJoined optionKind = "remaining-args-joined"
)

type option struct {
	visibility []optionVisibility
	kind       optionKind
	numArgs    int
	alias      *string
	aliasArgs  []string
}

type config struct {
	options map[string]option
}
