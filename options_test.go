package autoenv

import (
	"reflect"
	"slices"
	"testing"
)

func TestOptions_Apply(t *testing.T) {
	tests := []struct {
		name    string
		options []Option
		want    options
	}{
		{
			name:    "empty",
			options: nil,
			want:    defaultOptions,
		},
		{
			name:    "single option",
			options: []Option{WithVerbose()},
			want: options{
				prefix:     defaultOptions.prefix,
				logger:     defaultOptions.logger,
				filesPaths: defaultOptions.filesPaths,
				ignores:    defaultOptions.ignores,
				onlyEnvTag: defaultOptions.onlyEnvTag,
				withFiles:  defaultOptions.withFiles,
				verbose:    true,
			},
		},
		{
			name: "multiple options",
			options: []Option{
				WithVerbose(),
				WithVerbose(),
				WithPath(".env.prod"),
				WithPath(".env.dev"),
			},
			want: options{
				prefix:     defaultOptions.prefix,
				logger:     defaultOptions.logger,
				filesPaths: slices.Concat(defaultOptions.filesPaths, []string{".env.prod", ".env.dev"}),
				ignores:    defaultOptions.ignores,
				onlyEnvTag: defaultOptions.onlyEnvTag,
				withFiles:  true,
				verbose:    true,
			},
		},
		{
			name: "all options",
			options: []Option{
				WithVerbose(),
				WithPath(".env.unknow"),
				WithPaths([]string{".env.dev", ".env.prod"}),
				WithIgnores([]string{"ignored_field"}),
				WithPath(".env.test"),
				WithLogger(defaultOptions.logger),
				WithOnlyEnvTag(),
				WithPrefix("VAR"),
				WithFiles(),
			},
			want: options{
				prefix:     "VAR",
				logger:     defaultOptions.logger,
				filesPaths: []string{".env.dev", ".env.prod", ".env.test"},
				ignores:    []string{"ignored_field"},
				onlyEnvTag: true,
				verbose:    true,
				withFiles:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newOptions()
			got.apply(tt.options...)

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
