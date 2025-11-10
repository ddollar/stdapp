package stdcli

func New(name, version string) *Engine {
	e := &Engine{
		Executor: &defaultExecutor{},
		Name:     name,
		Reader:   DefaultReader,
		Version:  version,
		Writer:   DefaultWriter,
	}

	e.Command("help", "list commands", help(e), CommandOptions{
		Validate: ArgsBetween(0, 1),
	})

	return e
}
