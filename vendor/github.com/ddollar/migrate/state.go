package migrate

type State map[string]bool

func LoadState(e *Engine) (State, error) {
	var ss []string

	if _, err := e.db.Query(&ss, "select * from _migrations"); err != nil {
		return nil, err
	}

	state := State{}

	for _, s := range ss {
		state[s] = true
	}

	return state, nil
}
