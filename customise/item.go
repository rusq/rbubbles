package customise

import "strconv"

type Item interface {
	Name() string
	Value() string
	Description() string
	Group() string
	Validate(s string) error
	Set(s string) error
	AllowedValues() []string
}

type VarWrapper struct {
	ItemName          string
	ItemDescr         string
	ItemGroup         string
	ValueFunc         func() string
	SetFunc           func(string) error
	ValidateFunc      func(string) error
	AllowedValuesFunc func() []string
}

func (w VarWrapper) Name() string            { return w.ItemName }
func (w VarWrapper) Value() string           { return w.ValueFunc() }
func (w VarWrapper) Description() string     { return w.ItemDescr }
func (w VarWrapper) Group() string           { return w.ItemGroup }
func (w VarWrapper) Validate(s string) error { return w.ValidateFunc(s) }
func (w VarWrapper) Set(s string) error      { return w.SetFunc(s) }
func (w VarWrapper) AllowedValues() []string { return w.AllowedValuesFunc() }

func StringVar(value *string, name, descr, group string, allowed ...string) VarWrapper {
	return VarWrapper{
		ItemName:  name,
		ItemDescr: descr,
		ItemGroup: group,
		ValueFunc: func() string { return *value },
		SetFunc: func(s string) error {
			*value = s
			return nil
		},
		ValidateFunc:      func(s string) error { return nil },
		AllowedValuesFunc: func() []string { return allowed },
	}
}

func IntVar[T ~int | ~int8 | ~int16 | ~int32 | ~int64](value *T, name, descr, group string, allowed ...string) VarWrapper {
	return VarWrapper{
		ItemName:  name,
		ItemDescr: descr,
		ItemGroup: group,
		ValueFunc: func() string { return strconv.FormatInt(int64(*value), 10) },
		SetFunc: func(s string) error {
			v, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			*value = T(v)
			return nil
		},
		ValidateFunc:      func(s string) error { return nil },
		AllowedValuesFunc: func() []string { return allowed },
	}
}
