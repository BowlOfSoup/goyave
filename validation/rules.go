package validation

// Rule function defining a validation rule.
// Passing rules should return true, false otherwise.
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

func validateString(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
	_, ok := value.(string)
	return ok
}
