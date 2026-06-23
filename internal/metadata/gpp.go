package metadata

type GPPMetadata struct {
	Command      string   `json:"command"`
	CompileFlags []string `json:"compile_flags"`
	CompileArgs  []string `json:"compile_args"`
}

func gppNew(args []string) (*Metadata, error) {
}

func gppJoin(_, rhs *Metadata) (*Metadata, error) {
	return rhs, nil
}
