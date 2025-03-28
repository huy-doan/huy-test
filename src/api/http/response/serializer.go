package response

import (
	"reflect"
	"time"
)

// Serializer is an interface for serializing data before sending it to clients
type Serializer interface {
	Serialize() interface{}
}

// SerializerFunc is a function type for transforming any data structure into a client-friendly representation
type SerializerFunc func(data interface{}) interface{}

// DefaultSerialize is the default serialization function that handles common data types
func DefaultSerialize(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	// If data already implements Serializer, use that
	if serializer, ok := data.(Serializer); ok {
		return serializer.Serialize()
	}

	val := reflect.ValueOf(data)
	typ := val.Type()

	// Handle pointers
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		typ = val.Type()
	}

	// Handle different types
	switch typ.Kind() {
	case reflect.Struct:
		return serializeStruct(val)
	case reflect.Slice, reflect.Array:
		return serializeSlice(val)
	case reflect.Map:
		return serializeMap(val)
	default:
		// For basic types, return as is
		return data
	}
}

// serializeStruct converts a struct to a map, applying transformations as needed
func serializeStruct(val reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Get JSON tag name or use field name
		key := field.Name
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			// Handle json tag options like omitempty
			for i, c := range jsonTag {
				if c == ',' {
					key = jsonTag[:i]
					break
				}
			}
			if key == "" {
				key = jsonTag
			}
		}

		// Get field value and serialize it
		fieldValue := val.Field(i)

		// Handle special types like time.Time
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			timeVal := fieldValue.Interface().(time.Time)
			result[key] = timeVal.Format(time.RFC3339)
			continue
		}

		// Recursively serialize the field value
		serialized := DefaultSerialize(fieldValue.Interface())

		// Skip nil values
		if serialized != nil {
			result[key] = serialized
		}
	}

	return result
}

// serializeSlice converts a slice to a slice of serialized elements
func serializeSlice(val reflect.Value) []interface{} {
	result := make([]interface{}, val.Len())

	for i := 0; i < val.Len(); i++ {
		result[i] = DefaultSerialize(val.Index(i).Interface())
	}

	return result
}

// serializeMap converts a map to a map with serialized values
func serializeMap(val reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})

	for _, key := range val.MapKeys() {
		// Only support string keys
		if key.Kind() == reflect.String {
			result[key.String()] = DefaultSerialize(val.MapIndex(key).Interface())
		}
	}

	return result
}
