package stdcli

import (
	"fmt"

	"github.com/pkg/errors"
)

type Validator func(ctx Context) error

func Args(num int) Validator {
	return func(ctx Context) error {
		if len(ctx.Args()) != num {
			return errors.Errorf("%d %s required", num, plural("arg", num))
		}
		return nil
	}
}

func ArgsBetween(min, max int) Validator {
	return func(ctx Context) error {
		if err := ArgsMin(min)(ctx); err != nil {
			return err
		}
		if err := ArgsMax(max)(ctx); err != nil {
			return err
		}
		return nil
	}
}

func ArgsMin(min int) Validator {
	return func(ctx Context) error {
		if len(ctx.Args()) < min {
			return errors.Errorf("at least %d %s required", min, plural("arg", min))
		}
		return nil
	}
}

func ArgsMax(max int) Validator {
	return func(ctx Context) error {
		if len(ctx.Args()) > max {
			return errors.Errorf("no more than %d %s expected", max, plural("arg", max))
		}
		return nil
	}
}

func plural(noun string, num int) string {
	if num == 1 {
		return noun
	}

	return fmt.Sprintf("%ss", noun)
}
