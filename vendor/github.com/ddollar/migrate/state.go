package migrate

import "github.com/ddollar/errors"

type State map[string]bool

func LoadState(e *Engine) (State, error) {
	rows, err := e.db.Query("select * from _migrations")
	if err != nil {
		return nil, errors.Wrap(err)
	}

	state := State{}

	for rows.Next() {
		var s string

		if err := rows.Scan(&s); err != nil {
			return nil, errors.Wrap(err)
		}

		state[s] = true
	}

	if err := rows.Close(); err != nil {
		return nil, errors.Wrap(err)
	}

	return state, nil
}
