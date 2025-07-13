package builtin

import (
	"fmt"
	"math/rand"
	"time"

	"likha/generator/types"

	"likha/util"
)

// BuiltinGenerator uses predefined functions to generate data.
type BuiltinGenerator struct {
	function func() (interface{}, error)
}

// New creates a new BuiltinGenerator.
func New(settings map[string]interface{}) (types.Generator, error) {
	funcName, ok := settings["function"].(string)
	if !ok {
		return nil, fmt.Errorf("builtin generator requires a 'function' string setting")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var f func() (interface{}, error)
	switch funcName {
	case "random_epoch":
		f = makeRandomEpoch(r, settings)
	case "random_isodate":
		f = makeRandomISODate(r, settings)
	case "random_string":
		f = makeRandomString(r, settings)
	case "random_int":
		f = makeRandomInt(r, settings)
	case "random_decimal":
		f = makeRandomDecimal(r, settings)
	default:
		return nil, fmt.Errorf("unknown builtin function: %s", funcName)
	}

	return &BuiltinGenerator{function: f}, nil
}

// Generate calls the configured builtin function.
func (g *BuiltinGenerator) Generate(row map[string]interface{}) (interface{}, error) {
	return g.function()
}

// Helper functions to create the specific generator functions
func makeRandomEpoch(r *rand.Rand, s map[string]interface{}) func() (interface{}, error) {
	start := time.Now().Add(-365 * 24 * time.Hour).Unix()
	end := time.Now().Unix()
	if v, ok := s["start"]; ok {
		start, _ = util.InterfaceToInt64(v)
	}
	if v, ok := s["end"]; ok {
		end, _ = util.InterfaceToInt64(v)
	}
	return func() (interface{}, error) {
		return r.Int63n(end-start+1) + start, nil
	}
}

func makeRandomISODate(r *rand.Rand, s map[string]interface{}) func() (interface{}, error) {
	start := time.Now().Add(-365 * 24 * time.Hour)
	end := time.Now()
	if v, ok := s["start_date"]; ok {
		parsed, err := time.Parse(time.RFC3339, v.(string))
		if err == nil {
			start = parsed
		}
	}
	if v, ok := s["end_date"]; ok {
		parsed, err := time.Parse(time.RFC3339, v.(string))
		if err == nil {
			end = parsed
		}
	}
	diff := end.Unix() - start.Unix()
	return func() (interface{}, error) {
		sec := r.Int63n(diff) + start.Unix()
		return time.Unix(sec, 0).Format(time.RFC3339), nil
	}
}

func makeRandomString(r *rand.Rand, s map[string]interface{}) func() (interface{}, error) {
	length := 10
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if v, ok := s["length"]; ok {
		l, _ := util.InterfaceToInt(v)
		length = l
	}
	if v, ok := s["charset"]; ok {
		charset = v.(string)
	}
	return func() (interface{}, error) {
		b := make([]byte, length)
		for i := range b {
			b[i] = charset[r.Intn(len(charset))]
		}
		return string(b), nil
	}
}

func makeRandomInt(r *rand.Rand, s map[string]interface{}) func() (interface{}, error) {
	min := 0
	max := 100
	if v, ok := s["min"]; ok {
		min, _ = util.InterfaceToInt(v)
	}
	if v, ok := s["max"]; ok {
		max, _ = util.InterfaceToInt(v)
	}
	return func() (interface{}, error) {
		return r.Intn(max-min+1) + min, nil
	}
}

func makeRandomDecimal(r *rand.Rand, s map[string]interface{}) func() (interface{}, error) {
	min := 0.0
	max := 100.0
	places := 2
	if v, ok := s["min"]; ok {
		min, _ = util.InterfaceToFloat64(v)
	}
	if v, ok := s["max"]; ok {
		max, _ = util.InterfaceToFloat64(v)
	}
	if v, ok := s["places"]; ok {
		places, _ = util.InterfaceToInt(v)
	}
	return func() (interface{}, error) {
		val := min + r.Float64()*(max-min)
		return fmt.Sprintf(fmt.Sprintf("%%.%df", places), val), nil
	}
}
