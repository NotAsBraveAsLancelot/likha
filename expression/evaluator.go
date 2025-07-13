package expression

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// Regex to find functions like $random_int(-10, 200) or $random_string(10, "abc")
	funcRegex = regexp.MustCompile(`\$random_(\w+)\(([^)]*)\)`)
	// Regex to find field references like #field_name
	fieldRegex = regexp.MustCompile(`#(\w+)`)
)

// Evaluator parses and evaluates expression strings.
type Evaluator struct {
	r *rand.Rand
}

// NewEvaluator creates a new expression evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Evaluate replaces function calls and field references in a template string with generated values.
func (e *Evaluator) Evaluate(template string, row map[string]interface{}) (string, error) {
	// Step 1: Replace all field references like #field_name with their values from the current row.
	// This is done first so that generated values can't be misinterpreted as field names.
	processedTemplate := fieldRegex.ReplaceAllStringFunc(template, func(match string) string {
		fieldName := strings.TrimPrefix(match, "#")
		if val, ok := row[fieldName]; ok {
			return fmt.Sprintf("%v", val)
		}
		// If field not found (e.g., forward reference), leave it as is for now.
		// The runner ensures order, so this shouldn't be an issue in practice.
		return match
	})

	// Step 2: Evaluate all $random_*() functions.
	var firstErr error
	result := funcRegex.ReplaceAllStringFunc(processedTemplate, func(match string) string {
		// If an error has already occurred, don't process further.
		if firstErr != nil {
			return match
		}

		parts := funcRegex.FindStringSubmatch(match)
		if len(parts) < 2 { // Should be at least $random_name()
			return match // Should not happen with this regex
		}
		funcName := parts[1]
		// Handle case where there are no arguments, e.g., $random_isodate()
		argsStr := ""
		if len(parts) > 2 {
			argsStr = parts[2]
		}

		var args []string
		if strings.TrimSpace(argsStr) != "" {
			args = strings.Split(argsStr, ",")
			for i := range args {
				args[i] = strings.TrimSpace(args[i])
			}
		}

		val, err := e.callFunc(funcName, args)
		if err != nil {
			firstErr = fmt.Errorf("error in expression '%s': %w", match, err)
			return match // In case of an error, return the original placeholder
		}
		return fmt.Sprintf("%v", val)
	})

	if firstErr != nil {
		return "", firstErr
	}

	return result, nil
}

// callFunc dispatches to the correct random generator function based on name.
func (e *Evaluator) callFunc(name string, args []string) (interface{}, error) {
	switch name {
	case "int":
		if len(args) != 2 {
			return nil, fmt.Errorf("random_int requires 2 arguments (min, max), got %d", len(args))
		}
		min, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, fmt.Errorf("invalid min for random_int: %w", err)
		}
		max, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, fmt.Errorf("invalid max for random_int: %w", err)
		}
		if min > max {
			return nil, fmt.Errorf("min cannot be greater than max for random_int")
		}
		return e.r.Intn(max-min+1) + min, nil

	case "string":
		length := 10
		charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		if len(args) >= 1 && args[0] != "" {
			l, err := strconv.Atoi(args[0])
			if err == nil {
				length = l
			}
		}
		if len(args) >= 2 && args[1] != "" {
			// Unquote the charset if it's quoted
			unquoted, err := strconv.Unquote(args[1])
			if err == nil {
				charset = unquoted
			} else {
				charset = args[1]
			}
		}
		b := make([]byte, length)
		for i := range b {
			b[i] = charset[e.r.Intn(len(charset))]
		}
		return string(b), nil

	case "decimal":
		min := 0.0
		max := 100.0
		places := 2
		if len(args) >= 1 && args[0] != "" {
			m, err := strconv.ParseFloat(args[0], 64)
			if err == nil {
				min = m
			}
		}
		if len(args) >= 2 && args[1] != "" {
			m, err := strconv.ParseFloat(args[1], 64)
			if err == nil {
				max = m
			}
		}
		if len(args) >= 3 && args[2] != "" {
			p, err := strconv.Atoi(args[2])
			if err == nil {
				places = p
			}
		}
		if min > max {
			return nil, fmt.Errorf("min cannot be greater than max for random_decimal")
		}
		val := min + e.r.Float64()*(max-min)
		return fmt.Sprintf(fmt.Sprintf("%%.%df", places), val), nil

	case "epoch":
		start := time.Now().Add(-365 * 24 * time.Hour).Unix()
		end := time.Now().Unix()
		if len(args) >= 1 && args[0] != "" {
			s, err := strconv.ParseInt(args[0], 10, 64)
			if err == nil {
				start = s
			}
		}
		if len(args) >= 2 && args[1] != "" {
			e, err := strconv.ParseInt(args[1], 10, 64)
			if err == nil {
				end = e
			}
		}
		if start > end {
			return nil, fmt.Errorf("start cannot be after end for random_epoch")
		}
		return e.r.Int63n(end-start+1) + start, nil

	case "isodate":
		start := time.Now().Add(-365 * 24 * time.Hour)
		end := time.Now()
		if len(args) >= 1 && args[0] != "" {
			s, err := time.Parse(time.RFC3339, args[0])
			if err == nil {
				start = s
			}
		}
		if len(args) >= 2 && args[1] != "" {
			e, err := time.Parse(time.RFC3339, args[1])
			if err == nil {
				end = e
			}
		}
		if start.After(end) {
			return nil, fmt.Errorf("start_date cannot be after end_date for random_isodate")
		}
		diff := end.Unix() - start.Unix()
		sec := e.r.Int63n(diff) + start.Unix()
		return time.Unix(sec, 0).Format(time.RFC3339), nil

	default:
		return nil, fmt.Errorf("unknown expression function: %s", name)
	}
}
