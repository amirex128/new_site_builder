package utils

import (
	"log"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	validationmessage "github.com/amirex128/new_site_builder/src/internal/api/utils/validation_message"
)

// ValidationHelper provides methods for handling validation in handlers
type ValidationHelper struct{}

type ValidationErrorBag struct {
	Property   string `json:"property"`
	PropertyFa string `json:"propertyFa"`
	Tag        string `json:"tag"`
	Value      string `json:"value"`
	Message    string `json:"message"`
}

// NewValidationHelper creates a new ValidationHelper instance
func NewValidationHelper() *ValidationHelper {
	return &ValidationHelper{}
}

// ValidateCommand handles binding and validating the input struct
// Returns true if validation passes, false otherwise
func (h *ValidationHelper) ValidateCommand(c *gin.Context, params interface{}) bool {
	if err := c.ShouldBindJSON(params); err != nil {
		ValidationErrorString(c, err.Error(), map[string]any{})
		return false
	}

	errs := ValidateStruct(params)
	if len(errs) > 0 {
		ValidationError(c, errs...)
		return false
	}
	return true
}

// ValidateQuery handles binding and validating query parameters
// Returns true if validation passes, false otherwise
func (h *ValidationHelper) ValidateQuery(c *gin.Context, params interface{}) bool {
	if err := c.ShouldBindQuery(params); err != nil {
		ValidationErrorString(c, err.Error(), map[string]any{})
		return false
	}
	errs := ValidateStruct(params)
	if len(errs) > 0 {
		ValidationError(c, errs...)
		return false
	}
	return true
}

// ValidateStruct validates all fields and collects all errors for all rules
func ValidateStruct(s interface{}) []ValidationErrorBag {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	var errorsBag []ValidationErrorBag

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}
		// Split by comma, but handle params with = or space
		tags := parseValidateTag(validateTag)
		errorsTag := field.Tag.Get("errors")
		errorMsgs := parseErrorsTag(errorsTag)
		nameFa := field.Tag.Get("nameFa")
		if nameFa == "" {
			nameFa = field.Name
		}
		for _, tag := range tags {
			if tag == "" {
				continue
			}
			tagName, tagParam := splitTagAndParam(tag)
			ok := runCustomValidator(tagName, tagParam, fieldValue)
			if !ok {
				msg := errorMsgs[tagName]
				if msg == "" {
					// Use default message from ValidationMessages
					msgTemplate, found := validationmessage.ValidationMessages[tagName]
					if found {
						// Prepare args for Sprintf: first is nameFa, then tagParam if present
						args := []interface{}{nameFa}
						if tagParam != "" {
							params := splitValidatorParams(tagParam)
							for _, p := range params {
								args = append(args, p)
							}
						}
						msg = sprintfWithArgs(msgTemplate, args)
					} else {
						msg = defaultErrorMessage(nameFa, tagName, tagParam)
					}
				}
				errorsBag = append(errorsBag, ValidationErrorBag{
					Property:   field.Name,
					PropertyFa: nameFa,
					Tag:        tagName,
					Value:      tagParam,
					Message:    msg,
				})
			}
		}
	}
	return errorsBag
}

// sprintfWithArgs formats a string with a variable number of arguments, similar to fmt.Sprintf, but without importing fmt for performance.
func sprintfWithArgs(format string, args []interface{}) string {
	// Only support up to 3 %s for our use case
	res := format
	for _, arg := range args {
		res = strings.Replace(res, "%s", toString(arg), 1)
	}
	return res
}

func toString(arg interface{}) string {
	switch v := arg.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}

// parseValidateTag splits a validate tag into individual rules
func parseValidateTag(tag string) []string {
	var tags []string
	var curr strings.Builder
	inParam := false
	for _, r := range tag {
		if r == ',' && !inParam {
			tags = append(tags, curr.String())
			curr.Reset()
			continue
		}
		if r == '=' || r == ' ' {
			inParam = true
		}
		if inParam && r == ',' {
			inParam = false
		}
		curr.WriteRune(r)
	}
	if curr.Len() > 0 {
		tags = append(tags, curr.String())
	}
	return tags
}

// splitTagAndParam splits a tag like "required_text=1 100" into ("required_text", "1 100")
func splitTagAndParam(tag string) (string, string) {
	if idx := strings.Index(tag, "="); idx != -1 {
		return tag[:idx], strings.TrimSpace(tag[idx+1:])
	}
	return tag, ""
}

// parseErrorsTag parses the errors struct tag into a map[tag]message
func parseErrorsTag(tag string) map[string]string {
	result := make(map[string]string)
	pairs := strings.Split(tag, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			result[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return result
}

// runCustomValidator runs the custom validator by name
func runCustomValidator(tag, param string, value reflect.Value) bool {
	switch tag {
	case "required":
		return RequiredValidatorValue(value)
	case "min":
		return MinValidatorValue(value, param)
	case "max":
		return MaxValidatorValue(value, param)
	case "len":
		return LenValidatorValue(value, param)
	case "email":
		return EmailValidatorValue(value)
	case "oneof":
		return OneOfValidatorValue(value, param)
	case "gt":
		return GtValidatorValue(value, param)
	case "lt":
		return LtValidatorValue(value, param)
	case "gte":
		return GteValidatorValue(value, param)
	case "lte":
		return LteValidatorValue(value, param)
	case "eq":
		return EqValidatorValue(value, param)
	case "ne":
		return NeValidatorValue(value, param)
	case "url":
		return URLValidatorValue(value)
	case "uuid":
		return UUIDValidatorValue(value)
	case "required_text":
		return RequiredTextValidatorValue(value, param)
	case "optional_text":
		return OptionalTextValidatorValue(value, param)
	case "required_bool":
		return RequiredBoolValidatorValue(value)
	case "optional_bool":
		return true
	case "domain":
		return RequiredDomainValidatorValue(value)
	case "domain_optional":
		return OptionalDomainValidatorValue(value)
	case "slug":
		return RequiredSlugValidatorValue(value)
	case "slug_optional":
		return OptionalSlugValidatorValue(value)
	case "comma_numbers":
		return RequiredCommaSeparatedNumbersValidatorValue(value)
	case "comma_numbers_optional":
		return OptionalCommaSeparatedNumbersValidatorValue(value)
	case "pattern":
		return RequiredPatternValidatorValue(value, param)
	case "pattern_optional":
		return OptionalPatternValidatorValue(value, param)
	case "array_string":
		return RequiredArrayStringValidatorValue(value, param)
	case "array_string_optional":
		return OptionalArrayStringValidatorValue(value, param)
	case "array_number":
		return RequiredArrayNumberValidatorValue(value, param)
	case "array_number_optional":
		return OptionalArrayNumberValidatorValue(value, param)
	case "enum":
		return RequiredEnumValidatorValue(value)
	case "enum_optional":
		return OptionalEnumValidatorValue(value)
	case "enum_string_map":
		return EnumStringMapValidatorValue(value)
	case "enum_string_map_optional":
		return OptionalEnumStringMapValidatorValue(value)
	case "iranian_mobile":
		return IranianMobileNumberValidatorValue(value)
	default:
		return true // unknown tags are ignored
	}
}

// defaultErrorMessage returns a default error message for a field/tag
func defaultErrorMessage(field, tag, param string) string {
	return "Validation failed for " + field + " with rule " + tag + " (" + param + ")"
}

// --- Custom validator helpers below ---

func RequiredValidatorValue(value reflect.Value) bool {
	if !value.IsValid() {
		return false
	}
	switch value.Kind() {
	case reflect.String:
		return value.String() != ""
	case reflect.Ptr, reflect.Interface:
		return !value.IsNil()
	case reflect.Slice, reflect.Array, reflect.Map:
		return value.Len() > 0
	default:
		zero := reflect.Zero(value.Type())
		return !reflect.DeepEqual(value.Interface(), zero.Interface())
	}
}

func MinValidatorValue(value reflect.Value, param string) bool {
	min, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return true // ignore invalid param
	}
	switch value.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		return float64(value.Len()) >= min
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(value.Int()) >= min
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(value.Uint()) >= min
	case reflect.Float32, reflect.Float64:
		return value.Float() >= min
	}
	return true
}

func MaxValidatorValue(value reflect.Value, param string) bool {
	max, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return true
	}
	switch value.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		return float64(value.Len()) <= max
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(value.Int()) <= max
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(value.Uint()) <= max
	case reflect.Float32, reflect.Float64:
		return value.Float() <= max
	}
	return true
}

func LenValidatorValue(value reflect.Value, param string) bool {
	l, err := strconv.Atoi(param)
	if err != nil {
		return true
	}
	switch value.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		return value.Len() == l
	}
	return true
}

func EmailValidatorValue(value reflect.Value) bool {
	s, ok := value.Interface().(string)
	if !ok || s == "" {
		return false
	}
	_, err := mail.ParseAddress(s)
	return err == nil
}

func OneOfValidatorValue(value reflect.Value, param string) bool {
	s := value.String()
	options := strings.Fields(param)
	for _, opt := range options {
		if s == opt {
			return true
		}
	}
	return false
}

func GtValidatorValue(value reflect.Value, param string) bool {
	return compareNumeric(value, param, ">")
}

func LtValidatorValue(value reflect.Value, param string) bool {
	return compareNumeric(value, param, "<")
}

func GteValidatorValue(value reflect.Value, param string) bool {
	return compareNumeric(value, param, ">=")
}

func LteValidatorValue(value reflect.Value, param string) bool {
	return compareNumeric(value, param, "<=")
}

func EqValidatorValue(value reflect.Value, param string) bool {
	return compareNumeric(value, param, "==")
}

func NeValidatorValue(value reflect.Value, param string) bool {
	return compareNumeric(value, param, "!=")
}

func compareNumeric(value reflect.Value, param string, op string) bool {
	p, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return true
	}
	var v float64
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = float64(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v = float64(value.Uint())
	case reflect.Float32, reflect.Float64:
		v = value.Float()
	default:
		return true
	}
	switch op {
	case ">":
		return v > p
	case "<":
		return v < p
	case ">=":
		return v >= p
	case "<=":
		return v <= p
	case "==":
		return v == p
	case "!=":
		return v != p
	}
	return true
}

func URLValidatorValue(value reflect.Value) bool {
	s, ok := value.Interface().(string)
	if !ok || s == "" {
		return false
	}
	_, err := url.ParseRequestURI(s)
	return err == nil
}

func UUIDValidatorValue(value reflect.Value) bool {
	s, ok := value.Interface().(string)
	if !ok || s == "" {
		return false
	}
	_, err := uuid.Parse(s)
	return err == nil
}

func RequiredTextValidatorValue(value reflect.Value, param string) bool {
	v, ok := value.Interface().(string)
	if !ok {
		return false
	}
	if strings.TrimSpace(v) == "" {
		return false
	}
	params := splitValidatorParams(param)
	if len(params) >= 2 {
		minLen, err1 := strconv.Atoi(params[0])
		maxLen, err2 := strconv.Atoi(params[1])
		if err1 == nil && len(v) < minLen {
			return false
		}
		if err2 == nil && len(v) > maxLen {
			return false
		}
	}
	return true
}

func OptionalTextValidatorValue(value reflect.Value, param string) bool {
	v, ok := value.Interface().(string)
	if !ok {
		return true
	}
	if strings.TrimSpace(v) == "" {
		return true
	}
	params := splitValidatorParams(param)
	if len(params) >= 2 {
		minLen, err1 := strconv.Atoi(params[0])
		maxLen, err2 := strconv.Atoi(params[1])
		if err1 == nil && len(v) < minLen {
			return false
		}
		if err2 == nil && len(v) > maxLen {
			return false
		}
	}
	return true
}

func RequiredBoolValidatorValue(value reflect.Value) bool {
	_, ok := value.Interface().(bool)
	return ok
}

func RequiredDomainValidatorValue(value reflect.Value) bool {
	v, ok := value.Interface().(string)
	if !ok {
		return false
	}
	if strings.TrimSpace(v) == "" {
		return false
	}
	return validateDomain(v)
}

func OptionalDomainValidatorValue(value reflect.Value) bool {
	v, ok := value.Interface().(string)
	if !ok {
		return true
	}
	if strings.TrimSpace(v) == "" {
		return true
	}
	return validateDomain(v)
}

func RequiredSlugValidatorValue(value reflect.Value) bool {
	v, ok := value.Interface().(string)
	if !ok {
		return false
	}
	if strings.TrimSpace(v) == "" {
		return false
	}
	return validateSlug(v)
}

func OptionalSlugValidatorValue(value reflect.Value) bool {
	v, ok := value.Interface().(string)
	if !ok {
		return true
	}
	if strings.TrimSpace(v) == "" {
		return true
	}
	return validateSlug(v)
}

func RequiredCommaSeparatedNumbersValidatorValue(value reflect.Value) bool {
	v, ok := value.Interface().(string)
	if !ok {
		return false
	}
	if strings.TrimSpace(v) == "" {
		return false
	}
	return validateCommaSeparatedNumbers(v)
}

func OptionalCommaSeparatedNumbersValidatorValue(value reflect.Value) bool {
	v, ok := value.Interface().(string)
	if !ok {
		return true
	}
	if strings.TrimSpace(v) == "" {
		return true
	}
	return validateCommaSeparatedNumbers(v)
}

func RequiredPatternValidatorValue(value reflect.Value, param string) bool {
	v, ok := value.Interface().(string)
	if !ok {
		return false
	}
	if strings.TrimSpace(v) == "" {
		return false
	}
	pattern := param
	if pattern == "" {
		pattern = "^[a-z0-9]+(?:-[a-z0-9]+)*$"
	}
	match, err := regexp.MatchString(pattern, v)
	if err != nil {
		return false
	}
	return match && len(v) <= 200
}

func OptionalPatternValidatorValue(value reflect.Value, param string) bool {
	v, ok := value.Interface().(string)
	if !ok {
		return true
	}
	if strings.TrimSpace(v) == "" {
		return true
	}
	pattern := param
	if pattern == "" {
		pattern = "^[a-z0-9]+(?:-[a-z0-9]+)*$"
	}
	match, err := regexp.MatchString(pattern, v)
	if err != nil {
		return false
	}
	return match && len(v) <= 200
}

func RequiredArrayStringValidatorValue(value reflect.Value, param string) bool {
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return false
	}
	if value.IsNil() || value.Len() == 0 {
		return false
	}
	params := splitValidatorParams(param)
	var minLength, maxLength int
	if len(params) >= 2 {
		minLength, _ = strconv.Atoi(params[0])
		maxLength, _ = strconv.Atoi(params[1])
	}
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		if item.Kind() != reflect.String {
			return false
		}
		strValue := item.String()
		if strValue == "" {
			return false
		}
		if minLength > 0 && len(strValue) < minLength {
			return false
		}
		if maxLength > 0 && len(strValue) > maxLength {
			return false
		}
	}
	return true
}

func OptionalArrayStringValidatorValue(value reflect.Value, param string) bool {
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return true
	}
	if value.IsNil() {
		return true
	}
	if value.Len() == 0 {
		return true
	}
	params := splitValidatorParams(param)
	var minLength, maxLength int
	if len(params) >= 2 {
		minLength, _ = strconv.Atoi(params[0])
		maxLength, _ = strconv.Atoi(params[1])
	}
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		if item.Kind() != reflect.String {
			return false
		}
		strValue := item.String()
		if minLength > 0 && len(strValue) < minLength {
			return false
		}
		if maxLength > 0 && len(strValue) > maxLength {
			return false
		}
	}
	return true
}

func RequiredArrayNumberValidatorValue(value reflect.Value, param string) bool {
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return false
	}
	if value.IsNil() || value.Len() == 0 {
		return false
	}
	params := splitValidatorParams(param)
	var minValue, maxValue, minLength, maxLength int64
	var canBeZero bool
	if len(params) >= 5 {
		minLength, _ = strconv.ParseInt(params[0], 10, 64)
		maxLength, _ = strconv.ParseInt(params[1], 10, 64)
		minValue, _ = strconv.ParseInt(params[2], 10, 64)
		maxValue, _ = strconv.ParseInt(params[3], 10, 64)
		canBeZero = params[4] == "true"
	}
	if minLength > 0 && int64(value.Len()) < minLength {
		return false
	}
	if maxLength > 0 && int64(value.Len()) > maxLength {
		return false
	}
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		var numValue int64
		switch item.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			numValue = item.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			numValue = int64(item.Uint())
		case reflect.Float32, reflect.Float64:
			numValue = int64(item.Float())
		default:
			return false
		}
		if !canBeZero && numValue == 0 {
			return false
		}
		if minValue != 0 && numValue < minValue {
			return false
		}
		if maxValue != 0 && numValue > maxValue {
			return false
		}
	}
	return true
}

func OptionalArrayNumberValidatorValue(value reflect.Value, param string) bool {
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return true
	}
	if value.IsNil() {
		return true
	}
	params := splitValidatorParams(param)
	var minValue, maxValue, minLength, maxLength int64
	var canBeZero bool
	if len(params) >= 5 {
		minLength, _ = strconv.ParseInt(params[0], 10, 64)
		maxLength, _ = strconv.ParseInt(params[1], 10, 64)
		minValue, _ = strconv.ParseInt(params[2], 10, 64)
		maxValue, _ = strconv.ParseInt(params[3], 10, 64)
		canBeZero = params[4] == "true"
	}
	if minLength > 0 && int64(value.Len()) < minLength {
		return false
	}
	if maxLength > 0 && int64(value.Len()) > maxLength {
		return false
	}
	if value.Len() == 0 {
		return true
	}
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		var numValue int64
		switch item.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			numValue = item.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			numValue = int64(item.Uint())
		case reflect.Float32, reflect.Float64:
			numValue = int64(item.Float())
		default:
			return false
		}
		if !canBeZero && numValue == 0 {
			return false
		}
		if minValue != 0 && numValue < minValue {
			return false
		}
		if maxValue != 0 && numValue > maxValue {
			return false
		}
	}
	return true
}

func RequiredEnumValidatorValue(value reflect.Value) bool {
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return false
		}
		value = value.Elem()
	}
	return !value.IsZero()
}

func OptionalEnumValidatorValue(value reflect.Value) bool {
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return true
		}
		value = value.Elem()
	}
	return true
}

func EnumStringMapValidatorValue(value reflect.Value) bool {
	if value.Kind() != reflect.Map {
		return false
	}
	if value.IsNil() || value.Len() == 0 {
		return false
	}
	iter := value.MapRange()
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			return false
		}
		if !isValidEnumKey(key) {
			return false
		}
		if val.Len() == 0 {
			return false
		}
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			if item.Kind() != reflect.String || item.String() == "" {
				return false
			}
		}
	}
	return true
}

func OptionalEnumStringMapValidatorValue(value reflect.Value) bool {
	if value.Kind() != reflect.Map {
		return true
	}
	if value.IsNil() || value.Len() == 0 {
		return true
	}
	iter := value.MapRange()
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			return false
		}
		if !isValidEnumKey(key) {
			return false
		}
		if val.Len() == 0 {
			return false
		}
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			if item.Kind() != reflect.String || item.String() == "" {
				return false
			}
		}
	}
	return true
}

func IranianMobileNumberValidatorValue(value reflect.Value) bool {
	v, ok := value.Interface().(string)
	if !ok {
		return false
	}
	iranianMobileNumberPattern := `^09(1[0-9]|2[0-2]|3[0-9]|9[0-9])[0-9]{7}$`
	res, err := regexp.MatchString(iranianMobileNumberPattern, v)
	if err != nil {
		log.Print(err.Error())
	}
	return res
}

func splitValidatorParams(param string) []string {
	if strings.Contains(param, " ") {
		return strings.Fields(param)
	}
	return strings.Split(param, ",")
}

func validateDomain(value string) bool {
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "www.") {
		return false
	}
	domainPattern := `^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`
	match, err := regexp.MatchString(domainPattern, value)
	if err != nil || !match {
		return false
	}
	_, err = url.Parse("https://" + value)
	return err == nil
}

func validateSlug(value string) bool {
	if len(value) > 200 {
		return false
	}
	slugPattern := `^[a-z0-9]+(?:-[a-z0-9]+)*$`
	match, err := regexp.MatchString(slugPattern, value)
	if err != nil {
		return false
	}
	return match
}

func validateCommaSeparatedNumbers(value string) bool {
	items := strings.Split(value, ",")
	if len(items) == 0 {
		return false
	}
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			return false
		}
		num, err := strconv.ParseInt(item, 10, 64)
		if err != nil || num <= 0 {
			return false
		}
	}
	return true
}

func isValidEnumKey(key reflect.Value) bool {
	switch key.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}
