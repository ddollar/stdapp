package migrate

type State map[string]bool

func LoadState(e *Engine) (State, error) {
	rows, err := e.db.Query("select * from _migrations")
	if err != nil {
		return nil, err
	}

	state := State{}

	for rows.Next() {
		var s string

		if err := rows.Scan(&s); err != nil {
			return nil, err
		}

		state[s] = true
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return state, nil
}
