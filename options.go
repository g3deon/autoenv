package autoenv

type options struct {
	verbose     bool
	onlyEnvTag  bool
	withEnvFile bool

	filePath string

	ignore []string
	logger Logger
}

type Option func(*options)

func newOptions() *options {
	return &options{
		verbose:     false,
		onlyEnvTag:  false,
		withEnvFile: false,
		filePath:    ".env",
		ignore:      []string{},
		logger:      &defaultLogger{},
	}
}

func WithVerbose(verbose bool) Option {
	return func(o *options) {
		o.verbose = verbose
	}
}

func WithIgnore(field string) Option {
	return func(o *options) {
		o.ignore = append(o.ignore, field)
	}
}

func WithIgnores(ignores []string) Option {
	return func(o *options) {
		o.ignore = ignores
	}
}

func WithOnlyEnvTag() Option {
	return func(o *options) {
		o.onlyEnvTag = true
	}
}

func WithEnvFilePath(envFilePath string) Option {
	return func(o *options) {
		o.filePath = envFilePath
	}
}

func WithEnvFile() Option {
	return func(o *options) {
		o.withEnvFile = true
	}
}

func WithLogger(logger Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}
