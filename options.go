package table

// Options to configure table parser
type Options struct {
	// Selector for rows to be parsed
	RowSelector string

	// Process columns with colspan attribute as one column each
	IgnoreColspan bool

	// Possible number of columns for alignment parsing
	NumberOfColumns []int
}

func getOptions(options ...*Options) *Options {
	var opts *Options
	if len(options) > 0 {
		opts = options[0]
	}
	if opts == nil {
		opts = &Options{}
	}
	return opts
}
