package cpolicy

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/intelsdi-x/pulse/core/ctypes"
)

// A rule validating against string-typed config
type FloatRule struct {
	rule

	key      string
	required bool
	default_ *float64
	minimum  *float64
	maximum  *float64
}

func (f *FloatRule) GobEncode() ([]byte, error) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	if err := encoder.Encode(f.key); err != nil {
		return nil, err
	}
	if err := encoder.Encode(f.required); err != nil {
		return nil, err
	}
	if f.default_ == nil {
		encoder.Encode(false)
	} else {
		encoder.Encode(true)
		if err := encoder.Encode(&f.default_); err != nil {
			return nil, err
		}
	}
	if f.minimum == nil {
		encoder.Encode(false)
	} else {
		encoder.Encode(true)
		if err := encoder.Encode(f.minimum); err != nil {
			return nil, err
		}
	}
	if f.maximum == nil {
		encoder.Encode(false)
	} else {
		encoder.Encode(true)
		if err := encoder.Encode(f.maximum); err != nil {
			return nil, err
		}
	}
	return w.Bytes(), nil
}

func (f *FloatRule) GobDecode(buf []byte) error {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	if err := decoder.Decode(&f.key); err != nil {
		return err
	}
	if err := decoder.Decode(&f.required); err != nil {
		return err
	}
	var is_default_set bool
	decoder.Decode(&is_default_set)
	if is_default_set {
		return decoder.Decode(&f.default_)
	}
	var is_minimum_set bool
	decoder.Decode(&is_minimum_set)
	if is_minimum_set {
		if err := decoder.Decode(&f.minimum); err != nil {
			return err
		}
	}
	var is_maximum_set bool
	decoder.Decode(&is_maximum_set)
	if is_maximum_set {
		if err := decoder.Decode(&f.maximum); err != nil {
			return err
		}
	}
	return nil
}

// Returns a new float-typed rule. Arguments are key(string), required(bool), default(float64), min(float64), max(float64)
func NewFloatRule(key string, req bool, opts ...float64) (*FloatRule, error) {
	// Return error if key is empty
	if key == "" {
		return nil, EmptyKeyError
	}

	options := make([]*float64, 1)
	for i, o := range opts {
		options[i] = &o
	}

	return &FloatRule{
		key:      key,
		required: req,
		default_: options[0],
	}, nil
}

// Returns the key
func (f *FloatRule) Key() string {
	return f.key
}

// Validates a config value against this rule.
func (f *FloatRule) Validate(cv ctypes.ConfigValue) error {
	// Check that type is correct
	if cv.Type() != "float" {
		return wrongType(f.key, cv.Type(), "float")
	}
	// Check minimum. Type should be safe now because of the check above.
	if f.minimum != nil && cv.(ctypes.ConfigValueFloat).Value < *f.minimum {
		return errors.New(fmt.Sprintf("value is under minimum (%s value %f < %f)", f.key, cv.(ctypes.ConfigValueFloat).Value, *f.minimum))
	}
	// Check maximum. Type should be safe now because of the check above.
	if f.maximum != nil && cv.(ctypes.ConfigValueFloat).Value > *f.maximum {
		return errors.New(fmt.Sprintf("value is over maximum (%s value %f > %f)", f.key, cv.(ctypes.ConfigValueFloat).Value, *f.maximum))
	}
	return nil
}

func (f *FloatRule) Default() ctypes.ConfigValue {
	if f.default_ != nil {
		return &ctypes.ConfigValueFloat{Value: *f.default_}
	}
	return nil
}

func (f *FloatRule) Required() bool {
	return f.required
}

func (f *FloatRule) SetMinimum(m float64) {
	f.minimum = &m
}

func (f *FloatRule) SetMaximum(m float64) {
	f.maximum = &m
}
