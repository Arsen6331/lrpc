package reflectutil

import (
	"encoding"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// Convert attempts to convert the given value to the given type
func Convert(in reflect.Value, toType reflect.Type) (reflect.Value, error) {
	// Get input type
	inType := in.Type()

	// If input is already the desired type, return
	if inType == toType {
		return in, nil
	}

	// If input can be converted to desired type, convert and return
	if in.CanConvert(toType) {
		return in.Convert(toType), nil
	}

	// Create new value of desired type
	to := reflect.New(toType).Elem()

	// If type is a pointer
	if to.Kind() == reflect.Pointer {
		// Initialize value
		to.Set(reflect.New(to.Type().Elem()))
	}

	switch val := in.Interface().(type) {
	case string:
		// If desired type satisfies text unmarshaler
		if u, ok := to.Interface().(encoding.TextUnmarshaler); ok {
			// Use text unmarshaler to get value
			err := u.UnmarshalText([]byte(val))
			if err != nil {
				return reflect.Value{}, err
			}

			// Return unmarshaled value
			return reflect.ValueOf(any(u)), nil
		}
	case []byte:
		// If desired type satisfies binary unmarshaler
		if u, ok := to.Interface().(encoding.BinaryUnmarshaler); ok {
			// Use binary unmarshaler to get value
			err := u.UnmarshalBinary(val)
			if err != nil {
				return reflect.Value{}, err
			}

			// Return unmarshaled value
			return reflect.ValueOf(any(u)), nil
		}
	}

	// If input is a map
	if in.Kind() == reflect.Map {
		// Use mapstructure to decode value
		err := mapstructure.Decode(in.Interface(), to.Addr().Interface())
		if err == nil {
			return to, nil
		}
	}

	// If input is a slice of any, and output is an array or slice
	if in.Type() == reflect.TypeOf([]any{}) &&
		to.Kind() == reflect.Slice || to.Kind() == reflect.Array {
		// Use ConvertSlice to convert value
		to.Set(reflect.ValueOf(ConvertSlice(
			in.Interface().([]any),
			toType,
		)))
	}

	return to, fmt.Errorf("cannot convert %s to %s", inType, toType)
}

// ConvertSlice converts []any to an array or slice, as provided
// in the "to" argument.
func ConvertSlice(in []any, to reflect.Type) any {
	// Create new value for output
	out := reflect.New(to).Elem()

	// If output value is a slice
	if out.Kind() == reflect.Slice {
		// Get type of slice elements
		outType := out.Type().Elem()

		// For every value provided
		for i := 0; i < len(in); i++ {
			// Get value of input type
			inVal := reflect.ValueOf(in[i])
			// Create new output type
			outVal := reflect.New(outType).Elem()

			// If types match
			if inVal.Type() == outType {
				// Set output value to input value
				outVal.Set(inVal)
			} else {
				// If input value can be converted to output type
				if inVal.CanConvert(outType) {
					// Convert and set output value to input value
					outVal.Set(inVal.Convert(outType))
				} else {
					// Set output value to its zero value
					outVal.Set(reflect.Zero(outVal.Type()))
				}
			}

			// Append output value to slice
			out = reflect.Append(out, outVal)
		}
	} else if out.Kind() == reflect.Array && out.Len() == len(in) {
		//If output type is array and lengths match

		// For every input value
		for i := 0; i < len(in); i++ {
			// Get matching output index
			outVal := out.Index(i)
			// Get input value
			inVal := reflect.ValueOf(in[i])

			// If types match
			if inVal.Type() == outVal.Type() {
				// Set output value to input value
				outVal.Set(inVal)
			} else {
				// If input value can be converted to output type
				if inVal.CanConvert(outVal.Type()) {
					// Convert and set output value to input value
					outVal.Set(inVal.Convert(outVal.Type()))
				} else {
					// Set output value to its zero value
					outVal.Set(reflect.Zero(outVal.Type()))
				}
			}
		}
	}

	// Return created value
	return out.Interface()
}
