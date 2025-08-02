package autoenv

var (
	loader *Loader
)

func init() {
	loader = NewLoader()
}

func Load(i any) error {
	if loader == nil {
		return ErrNilLoader
	}

	return loader.Load(i)
}
