package messageformat

import (
	"bytes"
	"fmt"
	"strconv"
)

// formatOrdinal is the format function associated with the "selectordinal" type.
//
// It will returns an error if :
// - the associated value can't be convert to string or to an int (i.e. bool, ...)
// - the pluralFunc is not defined (MessageFormat.getNamedKey)
//
// It will falls back to the "other" choice if :
// - its key can't be found in the given map
// - the computed named key (MessageFormat.getNamedKey) is not a key of the given map
func (f *formatter) formatOrdinal(expr Expression, ptr_output *bytes.Buffer, data map[string]any) error {
	o, ok := expr.(*SelectExpr)
	if !ok {
		return fmt.Errorf("InvalidExprType: want SelectExpr, got %T", expr)
	}

	value, err := toString(data, o.Key)
	if err != nil {
		return err
	}

	var choice *ParseTree

	if v, ok := data[o.Key]; ok {
		switch val := v.(type) {
		default:
			return fmt.Errorf("UnsupportedOrdinalKeyType: %T", val)

		case int, float64:

		case string:
			_, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
		}

		key, err := f.getNamedKey(v, true)
		if err != nil {
			return err
		}
		choice = o.Choices[key]
	}

	if choice == nil {
		choice = o.Choices["other"]
	}

	return f.format(choice, ptr_output, data, value)
}
