package autoenv

var defaultOptions = options{
	prefix:     "",
	logger:     &defaultLogger{},
	filesPaths: []string{".env", ".env.local"},
	ignores:    []string{},
	onlyEnvTag: false,
	withFiles:  false,
	verbose:    false,
}

type options struct {
	prefix string
	logger Logger

	filesPaths []string
	ignores    []string

	onlyEnvTag bool
	withFiles  bool
	verbose    bool
}

type Option func(*options)

func newOptions() options {
	return defaultOptions
}

func (o *options) apply(options ...Option) {
	for _, option := range options {
		option(o)
	}
}

func WithPaths(fileNames []string) Option {
	return func(o *options) {
		o.withFiles = true
		o.filesPaths = fileNames
	}
}

func WithPath(fileName string) Option {
	return func(o *options) {
		o.withFiles = true
		o.filesPaths = append(o.filesPaths, fileName)
	}
}

func WithPrefix(prefix string) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}

func WithIgnores(ignores []string) Option {
	return func(o *options) {
		o.ignores = ignores
	}
}

func WithLogger(logger Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func WithOnlyEnvTag() Option {
	return func(o *options) {
		o.onlyEnvTag = true
	}
}

func WithVerbose() Option {
	return func(o *options) {
		o.verbose = true
	}
}

func WithFiles() Option {
	return func(o *options) {
		o.withFiles = true
	}
}
