package validation

import (
	"reflect"
	"strconv"

	"github.com/System-Glitch/goyave/helpers"
)

// Rule function defining a validation rule.
// Passing rules should return true, false otherwise.
//
// Rules can modifiy the validated value if needed.
// For example, the "numeric" rule converts the data to float64 if it's a string.
type Rule func(string, interface{}, []string, map[string]interface{}) bool

func validateRequired(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	val, ok := form[field]
	if ok {
		if val == nil {
			return false
		}
		if str, okStr := val.(string); okStr && str == "" {
			return false
		}
	}
	return ok
}

func validateMin(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	requireParametersCount("min", parameters, 1)
	min, err := strconv.ParseFloat(parameters[0], 64)
	if err != nil {
		panic(err)
	}
	switch getFieldType(value) {
	case "numeric":
		floatValue, _ := helpers.ToFloat64(value)
		return floatValue >= min
	case "string":
		return len(value.(string)) >= int(min)
	case "array":
		list := reflect.ValueOf(value)
		return list.Len() >= int(min)
	case "file":
		return false // TODO implement file min size
	}

	return true // Pass if field type cannot be checked (bool, dates, ...)
}

func validateMax(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	requireParametersCount("max", parameters, 1)
	max, err := strconv.ParseFloat(parameters[0], 64)
	if err != nil {
		panic(err)
	}
	switch getFieldType(value) {
	case "numeric":
		floatValue, _ := helpers.ToFloat64(value)
		return floatValue <= max
	case "string":
		return len(value.(string)) <= int(max)
	case "array":
		list := reflect.ValueOf(value)
		return list.Len() <= int(max)
	case "file":
		return false // TODO implement file max size
	}

	return true // Pass if field type cannot be checked (bool, dates, ...)
}

func validateBetween(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	requireParametersCount("between", parameters, 2)
	min, errMin := strconv.ParseFloat(parameters[0], 64)
	max, errMax := strconv.ParseFloat(parameters[1], 64)
	if errMin != nil {
		panic(errMin)
	}
	if errMax != nil {
		panic(errMax)
	}

	switch getFieldType(value) {
	case "numeric":
		floatValue, _ := helpers.ToFloat64(value)
		return floatValue >= min && floatValue <= max
	case "string":
		length := len(value.(string))
		return length >= int(min) && length <= int(max)
	case "array":
		list := reflect.ValueOf(value)
		length := list.Len()
		return length >= int(min) && length <= int(max)
	case "file":
		return false // TODO implement file between size
	}

	return true // Pass if field type cannot be checked (bool, dates, ...)
}

func validateGreaterThan(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	requireParametersCount("greater_than", parameters, 1)
	valueType := getFieldType(value)

	compared, exists := form[parameters[0]]
	if !exists || valueType != getFieldType(compared) {
		return false // Can't compare two different types or missing field
	}

	switch valueType {
	case "numeric":
		floatValue, _ := helpers.ToFloat64(value)
		comparedFloatValue, _ := helpers.ToFloat64(compared)
		return floatValue > comparedFloatValue
	case "string":
		return len(value.(string)) > len(compared.(string))
	case "array":
		return reflect.ValueOf(value).Len() > reflect.ValueOf(compared).Len()
	case "file":
		return false // TODO implement file greater than size
	}

	return false
}

func validateGreaterThanEqual(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	requireParametersCount("greater_than_equal", parameters, 1)
	valueType := getFieldType(value)

	compared, exists := form[parameters[0]]
	if !exists || valueType != getFieldType(compared) {
		return false // Can't compare two different types or missing field
	}

	switch valueType {
	case "numeric":
		floatValue, _ := helpers.ToFloat64(value)
		comparedFloatValue, _ := helpers.ToFloat64(compared)
		return floatValue >= comparedFloatValue
	case "string":
		return len(value.(string)) >= len(compared.(string))
	case "array":
		return reflect.ValueOf(value).Len() >= reflect.ValueOf(compared).Len()
	case "file":
		return false // TODO implement file greater than size
	}

	return false
}

func validateLowerThan(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	requireParametersCount("lower_than", parameters, 1)
	valueType := getFieldType(value)

	compared, exists := form[parameters[0]]
	if !exists || valueType != getFieldType(compared) {
		return false // Can't compare two different types or missing field
	}

	switch valueType {
	case "numeric":
		floatValue, _ := helpers.ToFloat64(value)
		comparedFloatValue, _ := helpers.ToFloat64(compared)
		return floatValue < comparedFloatValue
	case "string":
		return len(value.(string)) < len(compared.(string))
	case "array":
		return reflect.ValueOf(value).Len() < reflect.ValueOf(compared).Len()
	case "file":
		return false // TODO implement file greater than size
	}

	return false
}

func validateLowerThanEqual(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	requireParametersCount("lower_than_equal", parameters, 1)
	valueType := getFieldType(value)

	compared, exists := form[parameters[0]]
	if !exists || valueType != getFieldType(compared) {
		return false // Can't compare two different types or missing field
	}

	switch valueType {
	case "numeric":
		floatValue, _ := helpers.ToFloat64(value)
		comparedFloatValue, _ := helpers.ToFloat64(compared)
		return floatValue <= comparedFloatValue
	case "string":
		return len(value.(string)) <= len(compared.(string))
	case "array":
		return reflect.ValueOf(value).Len() <= reflect.ValueOf(compared).Len()
	case "file":
		return false // TODO implement file greater than size
	}

	return false
}

// -------------------------
// Misc

// boolean (accept 1, "on", "true", "yes")
// different
// nullable
// confirmed
