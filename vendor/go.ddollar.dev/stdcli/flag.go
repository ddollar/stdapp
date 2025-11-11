package stdcli

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go.ddollar.dev/ddl"
	"go.ddollar.dev/errors"
)

type FlagType string

const (
	FlagBool     FlagType = "bool"
	FlagDuration FlagType = "duration"
	FlagInt      FlagType = "int"
	FlagString   FlagType = "string"
)

type Flag struct {
	Default     any
	Description string
	Name        string
	Short       string
	Value       any

	kind FlagType
}

type Flags []*Flag

func BoolFlag(name, short, description string) Flag {
	return Flag{
		Description: description,
		Name:        name,
		Short:       short,
		kind:        FlagBool,
	}
}

func DurationFlag(name, short, description string) Flag {
	return Flag{
		Description: description,
		Name:        name,
		Short:       short,
		kind:        FlagDuration,
	}
}

func IntFlag(name, short, description string) Flag {
	return Flag{
		Description: description,
		Name:        name,
		Short:       short,
		kind:        FlagInt,
	}
}

func StringFlag(name, short, description string) Flag {
	return Flag{
		Description: description,
		Name:        name,
		Short:       short,
		kind:        FlagString,
	}
}

func (f *Flag) Set(v string) error {
	switch f.Kind() {
	case FlagBool:
		f.Value = (v == "true")
	case FlagDuration:
		d, err := time.ParseDuration(v)
		if err != nil {
			return errors.Wrap(err)
		}
		f.Value = d
	case FlagInt:
		i, err := strconv.Atoi(v)
		if err != nil {
			return errors.Wrap(err)
		}
		f.Value = i
	case FlagString:
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
	return string(f.kind)
}

func (f *Flag) Kind() FlagType {
	return f.kind
}

func (f *Flag) Usage() string {
	command := ddl.If(f.Short != "", fmt.Sprintf("-%s", f.Short), "  ")
	command += fmt.Sprintf(" --%s", f.Name)

	switch f.Kind() {
	case FlagBool:
		return command
	case FlagDuration, FlagInt, FlagString:
		return fmt.Sprintf("%s <u><info><%s></info></u>", command, f.Name)
	default:
		panic(fmt.Sprintf("unknown flag type: %s", f.Type()))
	}
}

func (fs Flags) Bool(name string) bool {
	if f, ok := fs.find(name, FlagBool); ok {
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
	if f, ok := fs.find(name, FlagInt); ok {
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
	if f, ok := fs.find(name, FlagString); ok {
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

func (fs Flags) find(name string, kind FlagType) (*Flag, bool) {
	for _, f := range fs {
		if f.Name == name && f.Kind() == kind {
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

func typeString(v reflect.Type) FlagType {
	switch v.String() {
	case "bool":
		return FlagBool
	case "int":
		return FlagInt
	case "string":
		return FlagString
	case "time.Duration":
		return FlagDuration
	default:
		panic(fmt.Sprintf("unknown flag type: %s", v))
	}
}
