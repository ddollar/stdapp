package stdcli

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Flag struct {
	Default     any
	Description string
	Name        string
	Short       string
	Value       any

	kind string
}

type Flags []*Flag

func BoolFlag(name, short, description string) Flag {
	return Flag{
		Description: description,
		Name:        name,
		Short:       short,
		kind:        "bool",
	}
}

func DurationFlag(name, short, description string) Flag {
	return Flag{
		Description: description,
		Name:        name,
		Short:       short,
		kind:        "duration",
	}
}

func IntFlag(name, short, description string) Flag {
	return Flag{
		Description: description,
		Name:        name,
		Short:       short,
		kind:        "int",
	}
}

func StringFlag(name, short, description string) Flag {
	return Flag{
		Description: description,
		Name:        name,
		Short:       short,
		kind:        "string",
	}
}

func (f *Flag) Set(v string) error {
	switch f.Type() {
	case "bool":
		f.Value = (v == "true")
	case "duration":
		d, err := time.ParseDuration(v)
		if err != nil {
			return errors.WithStack(err)
		}
		f.Value = d
	case "int":
		i, err := strconv.Atoi(v)
		if err != nil {
			return errors.WithStack(err)
		}
		f.Value = i
	case "string":
		f.Value = v
	default:
		return errors.Errorf("unknown flag type: %s", f.Type())
	}

	return nil
}

func (f *Flag) String() string {
	return fmt.Sprintf("%s: %v", f.Name, f.Value)
}

func (f *Flag) Type() string {
	return f.kind
}

func (f *Flag) Usage(v string) string {
	switch f.Type() {
	case "bool":
		return fmt.Sprintf("%s <u><info></info></u>", v)
	case "duration", "int", "string":
		return fmt.Sprintf("%s <u><info>%s</info></u>", v, f.Name)
	default:
		panic(fmt.Sprintf("unknown flag type: %s", f.Type()))
	}
}

func (f *Flag) UsageLong() string {
	return f.Usage(fmt.Sprintf("--%s", f.Name))
}

func (f *Flag) UsageShort() string {
	if f.Short == "" {
		return ""
	}

	return f.Usage(fmt.Sprintf("-%s", f.Short))
}

func (fs Flags) Bool(name string) bool {
	if f, ok := fs.find(name, "bool"); ok {
		switch t := f.Value.(type) {
		case nil:
			v, _ := f.Default.(bool)
			return v
		case bool:
			return t
		}
	}

	return false
}

func (fs Flags) Int(name string) int {
	if f, ok := fs.find(name, "int"); ok {
		switch t := f.Value.(type) {
		case nil:
			v, _ := f.Default.(int)
			return v
		case int:
			return t
		}
	}

	return 0
}

func (fs Flags) String(name string) string {
	if f, ok := fs.find(name, "string"); ok {
		switch t := f.Value.(type) {
		case nil:
			v, _ := f.Default.(string)
			return v
		case string:
			return t
		}
	}

	return ""
}

func (fs Flags) Value(name string) any {
	for _, f := range fs {
		if f.Name == name {
			return f.Value
		}
	}

	return nil
}

func (fs Flags) find(name, kind string) (*Flag, bool) {
	for _, f := range fs {
		if f.Name == name && f.Type() == kind {
			return f, true
		}
	}
	return nil, false
}

func OptionFlags(opts any) []Flag {
	flags := []Flag{}

	v := reflect.ValueOf(opts)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if n := f.Tag.Get("flag"); n != "" {
			parts := strings.Split(n, ",")
			flag := Flag{
				Default:     f.Tag.Get("default"),
				Description: f.Tag.Get("desc"),
				Name:        parts[0],
				kind:        typeString(f.Type.Elem()),
			}
			if len(parts) > 1 {
				flag.Short = parts[1]
			}
			flags = append(flags, flag)
		}
	}

	return flags
}

func typeString(v reflect.Type) string {
	switch v.String() {
	case "bool":
		return "bool"
	case "int":
		return "int"
	case "string":
		return "string"
	case "time.Duration":
		return "duration"
	default:
		panic(fmt.Sprintf("unknown flag type: %s", v))
	}
}
