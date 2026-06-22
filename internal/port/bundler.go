package port

type Bundler interface {
	Bundle(sourcePath string) (string, error)
}
