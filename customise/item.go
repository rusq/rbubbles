package customise

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
)

type Item interface {
	Name() string
	Value() string
	Description() string
	Group() string
	Validate(s string) error
	Set(s string) error
	AllowedValues() []string
	Type() Type
}

type Type int

const (
	TText Type = iota
	TMultiline
	TRadio
	TCheckbox // for booleans
	TFile
)

type VarWrapper struct {
	ItemName          string
	ItemDescr         string
	ItemGroup         string
	ItemType          Type
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
func (w VarWrapper) Type() Type              { return w.ItemType }

func StringVar(value *string, name, descr, group string) VarWrapper {
	return VarWrapper{
		ItemName:  name,
		ItemDescr: descr,
		ItemGroup: group,
		ItemType:  TText,
		ValueFunc: func() string { return *value },
		SetFunc: func(s string) error {
			*value = s
			return nil
		},
		ValidateFunc: func(s string) error { return nil },
	}
}

func MultilineVar(value *string, name, descr, group string) VarWrapper {
	t := StringVar(value, name, descr, group)
	t.ItemType = TMultiline
	return t
}

func IntVar[T ~int | ~int8 | ~int16 | ~int32 | ~int64](value *T, name, descr, group string) VarWrapper {
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
		ValidateFunc: func(s string) error { return nil },
	}
}

func BoolVar(value *bool, name, descr, group string) VarWrapper {
	return VarWrapper{
		ItemName:  name,
		ItemDescr: descr,
		ItemGroup: group,
		ItemType:  TCheckbox,
		ValueFunc: func() string {
			return strconv.FormatBool(*value)
		},
		SetFunc: func(s string) error {
			v, err := strconv.ParseBool(s)
			if err != nil {
				return err
			}
			*value = v
			return nil
		},
		ValidateFunc: func(s string) error {
			_, err := strconv.ParseBool(s)
			return err
		},
		AllowedValuesFunc: func() []string {
			return []string{sTrue, sFalse}
		},
	}
}

var ErrInvalidValue = errors.New("invalid value")

func RadioStringVar(value *string, name, descr, group string, choices []string) VarWrapper {
	validateFunc := func(s string) error {
		if !slices.Contains(choices, s) {
			return fmt.Errorf("%w: %q", ErrInvalidValue, s)
		}
		return nil
	}

	return VarWrapper{
		ItemName:  name,
		ItemDescr: descr,
		ItemGroup: group,
		ItemType:  TRadio,
		ValueFunc: func() string {
			return *value
		},
		SetFunc: func(s string) error {
			if err := validateFunc(s); err != nil {
				return err
			}
			*value = s
			return nil
		},
		ValidateFunc: validateFunc,
		AllowedValuesFunc: func() []string {
			return choices
		},
	}
}

func FilenameVar(value *string, name, descr, group string) VarWrapper {
	v := StringVar(value, name, descr, group)
	v.ItemType = TFile
	return v
}

var (
	sTrue  = strconv.FormatBool(true)
	sFalse = strconv.FormatBool(false)
)
