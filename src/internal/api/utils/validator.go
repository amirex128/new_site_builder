package utils

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

type ValidationError struct {
	Property string `json:"property"`
	Tag      string `json:"tag"`
	Value    string `json:"value"`
	Message  string `json:"message"`
}

func ValidateStruct(s interface{}) error {
	if validate == nil {
		validate = initValidator()
	}

	err := validate.Struct(s)

	return err
}

func Validate(field any, rule string) error {
	if validate == nil {
		validate = initValidator()
	}

	return validate.Var(field, rule)
}

func GetValidationErrors(err error) *[]ValidationError {
	var validationErrors []ValidationError
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, err := range err.(validator.ValidationErrors) {
			var el ValidationError
			el.Property = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			validationErrors = append(validationErrors, el)
		}
		return &validationErrors
	}

	return nil
}

// initValidator initializes and configures the validator
func initValidator() *validator.Validate {
	v := validator.New(
		validator.WithRequiredStructEnabled(),
		validator.WithPrivateFieldValidation(),
	)

	// Register custom validation functions
	_ = v.RegisterValidation("iranian_mobile", IranianMobileNumberValidator)
	_ = v.RegisterValidation("required_text", RequiredTextValidator)
	_ = v.RegisterValidation("optional_text", OptionalTextValidator)
	_ = v.RegisterValidation("required_bool", RequiredBoolValidator)
	_ = v.RegisterValidation("optional_bool", OptionalBoolValidator)
	_ = v.RegisterValidation("domain", RequiredDomainValidator)
	_ = v.RegisterValidation("domain_optional", OptionalDomainValidator)
	_ = v.RegisterValidation("slug", RequiredSlugValidator)
	_ = v.RegisterValidation("slug_optional", OptionalSlugValidator)
	_ = v.RegisterValidation("comma_numbers", RequiredCommaSeparatedNumbersValidator)
	_ = v.RegisterValidation("comma_numbers_optional", OptionalCommaSeparatedNumbersValidator)
	_ = v.RegisterValidation("pattern", RequiredPatternValidator)
	_ = v.RegisterValidation("pattern_optional", OptionalPatternValidator)
	_ = v.RegisterValidation("array_string", RequiredArrayStringValidator)
	_ = v.RegisterValidation("array_string_optional", OptionalArrayStringValidator)
	_ = v.RegisterValidation("array_number", RequiredArrayNumberValidator)
	_ = v.RegisterValidation("array_number_optional", OptionalArrayNumberValidator)
	_ = v.RegisterValidation("enum", RequiredEnumValidator)
	_ = v.RegisterValidation("enum_optional", OptionalEnumValidator)
	_ = v.RegisterValidation("enum_string_map", EnumStringMapValidator)
	_ = v.RegisterValidation("enum_string_map_optional", OptionalEnumStringMapValidator)

	return v
}

// IranianMobileNumberValidator validates Iranian mobile numbers
func IranianMobileNumberValidator(fld validator.FieldLevel) bool {
	iranianMobileNumberPattern := `^09(1[0-9]|2[0-2]|3[0-9]|9[0-9])[0-9]{7}$`
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return false
	}

	res, err := regexp.MatchString(iranianMobileNumberPattern, value)
	if err != nil {
		log.Print(err.Error())
	}

	return res
}

// RequiredTextValidator validates that a text field is not empty and meets length requirements
func RequiredTextValidator(fld validator.FieldLevel) bool {
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return false
	}

	if strings.TrimSpace(value) == "" {
		return false
	}

	params := strings.Split(fld.Param(), ",")
	if len(params) >= 2 {
		minLen, err1 := strconv.Atoi(params[0])
		maxLen, err2 := strconv.Atoi(params[1])

		if err1 == nil && len(value) < minLen {
			return false
		}
		if err2 == nil && len(value) > maxLen {
			return false
		}
	}

	return true
}

// OptionalTextValidator validates that a text field meets length requirements if it's not empty
func OptionalTextValidator(fld validator.FieldLevel) bool {
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return true // Optional field can be of wrong type
	}

	if strings.TrimSpace(value) == "" {
		return true // Empty is valid for optional
	}

	params := strings.Split(fld.Param(), ",")
	if len(params) >= 2 {
		minLen, err1 := strconv.Atoi(params[0])
		maxLen, err2 := strconv.Atoi(params[1])

		if err1 == nil && len(value) < minLen {
			return false
		}
		if err2 == nil && len(value) > maxLen {
			return false
		}
	}

	return true
}

// RequiredBoolValidator validates that a boolean field is not nil
func RequiredBoolValidator(fld validator.FieldLevel) bool {
	_, ok := fld.Field().Interface().(bool)
	return ok
}

// OptionalBoolValidator always returns true for optional boolean fields
func OptionalBoolValidator(fld validator.FieldLevel) bool {
	return true
}

// RequiredDomainValidator validates required domain format
func RequiredDomainValidator(fld validator.FieldLevel) bool {
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return false
	}

	if strings.TrimSpace(value) == "" {
		return false
	}

	return validateDomain(value)
}

// OptionalDomainValidator validates optional domain format
func OptionalDomainValidator(fld validator.FieldLevel) bool {
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return true
	}

	if strings.TrimSpace(value) == "" {
		return true
	}

	return validateDomain(value)
}

// validateDomain is a helper function for domain validation
func validateDomain(value string) bool {
	// Check if domain starts with http://, https://, or www.
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "www.") {
		return false
	}

	// Check domain format
	domainPattern := `^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`
	match, err := regexp.MatchString(domainPattern, value)
	if err != nil || !match {
		return false
	}

	// Check if domain can be parsed as host
	_, err = url.Parse("https://" + value)
	return err == nil
}

// RequiredSlugValidator validates required slug format
func RequiredSlugValidator(fld validator.FieldLevel) bool {
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return false
	}

	if strings.TrimSpace(value) == "" {
		return false
	}

	return validateSlug(value)
}

// OptionalSlugValidator validates optional slug format
func OptionalSlugValidator(fld validator.FieldLevel) bool {
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return true
	}

	if strings.TrimSpace(value) == "" {
		return true
	}

	return validateSlug(value)
}

// validateSlug is a helper function for slug validation
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

// RequiredCommaSeparatedNumbersValidator validates required comma-separated numbers
func RequiredCommaSeparatedNumbersValidator(fld validator.FieldLevel) bool {
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return false
	}

	if strings.TrimSpace(value) == "" {
		return false
	}

	return validateCommaSeparatedNumbers(value)
}

// OptionalCommaSeparatedNumbersValidator validates optional comma-separated numbers
func OptionalCommaSeparatedNumbersValidator(fld validator.FieldLevel) bool {
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return true
	}

	if strings.TrimSpace(value) == "" {
		return true
	}

	return validateCommaSeparatedNumbers(value)
}

// validateCommaSeparatedNumbers is a helper function for comma-separated numbers validation
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

// RequiredPatternValidator validates a string against a required pattern
func RequiredPatternValidator(fld validator.FieldLevel) bool {
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return false
	}

	if strings.TrimSpace(value) == "" {
		return false
	}

	pattern := fld.Param()
	if pattern == "" {
		pattern = "^[a-z0-9]+(?:-[a-z0-9]+)*$" // Default pattern
	}

	match, err := regexp.MatchString(pattern, value)
	if err != nil {
		return false
	}

	return match && len(value) <= 200
}

// OptionalPatternValidator validates a string against an optional pattern
func OptionalPatternValidator(fld validator.FieldLevel) bool {
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return true
	}

	if strings.TrimSpace(value) == "" {
		return true
	}

	pattern := fld.Param()
	if pattern == "" {
		pattern = "^[a-z0-9]+(?:-[a-z0-9]+)*$" // Default pattern
	}

	match, err := regexp.MatchString(pattern, value)
	if err != nil {
		return false
	}

	return match && len(value) <= 200
}

// RequiredArrayStringValidator validates a required array of strings
func RequiredArrayStringValidator(fld validator.FieldLevel) bool {
	field := fld.Field()

	// Check if it's an array/slice
	if field.Kind() != reflect.Slice && field.Kind() != reflect.Array {
		return false
	}

	// Check if the array is not nil and not empty
	if field.IsNil() || field.Len() == 0 {
		return false
	}

	params := strings.Split(fld.Param(), ",")
	var minLength, maxLength int
	if len(params) >= 2 {
		minLength, _ = strconv.Atoi(params[0])
		maxLength, _ = strconv.Atoi(params[1])
	}

	// Check each element
	for i := 0; i < field.Len(); i++ {
		item := field.Index(i)
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

// OptionalArrayStringValidator validates an optional array of strings
func OptionalArrayStringValidator(fld validator.FieldLevel) bool {
	field := fld.Field()

	// Check if it's an array/slice
	if field.Kind() != reflect.Slice && field.Kind() != reflect.Array {
		return true // Accept non-slice types for optional
	}

	// If nil, it's valid (optional)
	if field.IsNil() {
		return true
	}

	// Empty array is valid for optional
	if field.Len() == 0 {
		return true
	}

	params := strings.Split(fld.Param(), ",")
	var minLength, maxLength int
	if len(params) >= 2 {
		minLength, _ = strconv.Atoi(params[0])
		maxLength, _ = strconv.Atoi(params[1])
	}

	// Check each element
	for i := 0; i < field.Len(); i++ {
		item := field.Index(i)
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

// RequiredArrayNumberValidator validates a required array of numbers
func RequiredArrayNumberValidator(fld validator.FieldLevel) bool {
	field := fld.Field()

	// Check if it's an array/slice
	if field.Kind() != reflect.Slice && field.Kind() != reflect.Array {
		return false
	}

	// Check if the array is not nil and not empty
	if field.IsNil() || field.Len() == 0 {
		return false
	}

	params := strings.Split(fld.Param(), ",")
	var minValue, maxValue, minLength, maxLength int64
	var canBeZero bool

	if len(params) >= 5 {
		// Format: minLength,maxLength,minValue,maxValue,canBeZero
		minLength, _ = strconv.ParseInt(params[0], 10, 64)
		maxLength, _ = strconv.ParseInt(params[1], 10, 64)
		minValue, _ = strconv.ParseInt(params[2], 10, 64)
		maxValue, _ = strconv.ParseInt(params[3], 10, 64)
		canBeZero = params[4] == "true"
	}

	// Check array length constraints
	if minLength > 0 && int64(field.Len()) < minLength {
		return false
	}

	if maxLength > 0 && int64(field.Len()) > maxLength {
		return false
	}

	// Check each element
	for i := 0; i < field.Len(); i++ {
		item := field.Index(i)
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

// OptionalArrayNumberValidator validates an optional array of numbers
func OptionalArrayNumberValidator(fld validator.FieldLevel) bool {
	field := fld.Field()

	// Check if it's an array/slice
	if field.Kind() != reflect.Slice && field.Kind() != reflect.Array {
		return true
	}

	// If nil, it's valid (optional)
	if field.IsNil() {
		return true
	}

	params := strings.Split(fld.Param(), ",")
	var minValue, maxValue, minLength, maxLength int64
	var canBeZero bool

	if len(params) >= 5 {
		// Format: minLength,maxLength,minValue,maxValue,canBeZero
		minLength, _ = strconv.ParseInt(params[0], 10, 64)
		maxLength, _ = strconv.ParseInt(params[1], 10, 64)
		minValue, _ = strconv.ParseInt(params[2], 10, 64)
		maxValue, _ = strconv.ParseInt(params[3], 10, 64)
		canBeZero = params[4] == "true"
	}

	// Check array length constraints
	if minLength > 0 && int64(field.Len()) < minLength {
		return false
	}

	if maxLength > 0 && int64(field.Len()) > maxLength {
		return false
	}

	// Empty array is valid for optional
	if field.Len() == 0 {
		return true
	}

	// Check each element
	for i := 0; i < field.Len(); i++ {
		item := field.Index(i)
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

// RequiredEnumValidator validates that a value is a valid enum and not nil/zero
func RequiredEnumValidator(fld validator.FieldLevel) bool {
	field := fld.Field()

	// If the field is a pointer, dereference it
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return false
		}
		field = field.Elem()
	}

	// For enums, we simply check if it's not zero value since Go doesn't have a direct enum type
	return !field.IsZero()
}

// OptionalEnumValidator validates that a value is a valid enum if not nil
func OptionalEnumValidator(fld validator.FieldLevel) bool {
	field := fld.Field()

	// If the field is a pointer, it can be nil for optional
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	// Zero is acceptable for optional enum
	return true
}

// ValidateRequiredText validates a required text field with optional min and max length
func ValidateRequiredText(field string, value string, entityName string, minLength, maxLength *int) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s الزامی است", entityName)
	}

	if minLength != nil && len(value) < *minLength {
		return fmt.Errorf("%s باید حداقل %d کاراکتر باشد", entityName, *minLength)
	}

	if maxLength != nil && len(value) > *maxLength {
		return fmt.Errorf("%s نمی‌تواند بیشتر از %d کاراکتر باشد", entityName, *maxLength)
	}

	return nil
}

// ValidateOptionalText validates an optional text field with optional min and max length
func ValidateOptionalText(field string, value string, entityName string, minLength, maxLength *int) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	if minLength != nil && len(value) < *minLength {
		return fmt.Errorf("%s باید حداقل %d کاراکتر باشد", entityName, *minLength)
	}

	if maxLength != nil && len(value) > *maxLength {
		return fmt.Errorf("%s نمی‌تواند بیشتر از %d کاراکتر باشد", entityName, *maxLength)
	}

	return nil
}

// ValidateArrayString validates an array of strings
func ValidateArrayString(value []string, entityName string, required bool, minLength, maxLength *int) error {
	if required && (value == nil || len(value) == 0) {
		return fmt.Errorf("%s الزامی است", entityName)
	}

	if value == nil {
		return nil
	}

	for i, item := range value {
		if item == "" && required {
			return fmt.Errorf("آیتم %d از %s نمی‌تواند خالی باشد", i+1, entityName)
		}

		if minLength != nil && len(item) < *minLength {
			return fmt.Errorf("آیتم %d از %s باید حداقل %d کاراکتر باشد", i+1, entityName, *minLength)
		}

		if maxLength != nil && len(item) > *maxLength {
			return fmt.Errorf("آیتم %d از %s نمی‌تواند بیشتر از %d کاراکتر باشد", i+1, entityName, *maxLength)
		}
	}

	return nil
}

// ValidateRequiredEnum validates that an enum value is valid and not nil
func ValidateRequiredEnum(field string, value interface{}, entityName string) error {
	if value == nil {
		return fmt.Errorf("%s الزامی است", entityName)
	}

	// Check if enum value is zero (default)
	v := reflect.ValueOf(value)
	if v.IsZero() {
		return fmt.Errorf("%s نامعتبر است", entityName)
	}

	return nil
}

// ValidateOneOfPropertiesRequired checks if at least one of the provided properties has a valid value
func ValidateOneOfPropertiesRequired(object interface{}, properties []string, propertyNames []string) error {
	if object == nil {
		return fmt.Errorf("حداقل یکی از فیلدها باید مقدار داشته باشد")
	}

	objValue := reflect.ValueOf(object)
	if objValue.Kind() == reflect.Ptr {
		if objValue.IsNil() {
			return fmt.Errorf("حداقل یکی از فیلدها باید مقدار داشته باشد")
		}
		objValue = objValue.Elem()
	}

	if objValue.Kind() != reflect.Struct {
		return fmt.Errorf("ورودی باید یک ساختار باشد")
	}

	hasValidProperty := false

	for _, propName := range properties {
		field := objValue.FieldByName(propName)
		if !field.IsValid() {
			continue
		}

		if hasValidValue(field) {
			hasValidProperty = true
			break
		}
	}

	if !hasValidProperty {
		var displayNames string
		if propertyNames != nil && len(propertyNames) == len(properties) {
			displayNames = strings.Join(propertyNames, "، ")
		} else {
			displayNames = strings.Join(properties, "، ")
		}

		return fmt.Errorf("حداقل یکی از فیلدهای %s باید مقدار داشته باشد", displayNames)
	}

	return nil
}

// hasValidValue checks if a reflect.Value contains a valid non-zero value
func hasValidValue(field reflect.Value) bool {
	if !field.IsValid() {
		return false
	}

	// Check nil for interfaces, pointers, maps, slices, etc.
	if field.Kind() == reflect.Ptr ||
		field.Kind() == reflect.Interface ||
		field.Kind() == reflect.Slice ||
		field.Kind() == reflect.Map ||
		field.Kind() == reflect.Chan {
		if field.IsNil() {
			return false
		}
	}

	// Special cache for empty strings
	if field.Kind() == reflect.String && field.String() == "" {
		return false
	}

	// Special cache for zero numeric values
	if (field.Kind() == reflect.Int ||
		field.Kind() == reflect.Int8 ||
		field.Kind() == reflect.Int16 ||
		field.Kind() == reflect.Int32 ||
		field.Kind() == reflect.Int64 ||
		field.Kind() == reflect.Uint ||
		field.Kind() == reflect.Uint8 ||
		field.Kind() == reflect.Uint16 ||
		field.Kind() == reflect.Uint32 ||
		field.Kind() == reflect.Uint64 ||
		field.Kind() == reflect.Float32 ||
		field.Kind() == reflect.Float64) && field.IsZero() {
		return false
	}

	// Special cache for slices - check if empty
	if field.Kind() == reflect.Slice && field.Len() == 0 {
		return false
	}

	// Special cache for booleans - only true is valid
	if field.Kind() == reflect.Bool {
		return field.Bool()
	}

	// For other types, non-zero value is considered valid
	return !field.IsZero()
}

// ValidateStructCustom validates a struct and returns custom validation errors
func ValidateStructCustom(s interface{}) *[]ValidationError {
	err := ValidateStruct(s)
	if err != nil {
		return GetValidationErrors(err)
	}
	return nil
}

// Helper functions for use in struct tags

// Required returns a string for the required validator tag
func Required(entityName string) string {
	return "required"
}

// OptionalText returns a string for the optional text validator tag with min and max length
func OptionalText(minLength, maxLength int) string {
	return fmt.Sprintf("optional_text=%d,%d", minLength, maxLength)
}

// RequiredText returns a string for the required text validator tag with min and max length
func RequiredText(minLength, maxLength int) string {
	return fmt.Sprintf("required_text=%d,%d", minLength, maxLength)
}

// Slug returns a string for the slug validator tag
func Slug(optional bool) string {
	if optional {
		return "slug_optional"
	}
	return "slug"
}

// Domain returns a string for the domain validator tag
func Domain(optional bool) string {
	if optional {
		return "domain_optional"
	}
	return "domain"
}

// CommaSeparatedNumbers returns a string for the comma-separated numbers validator tag
func CommaSeparatedNumbers(optional bool) string {
	if optional {
		return "comma_numbers_optional"
	}
	return "comma_numbers"
}

// Pattern returns a string for the pattern validator tag
func Pattern(optional bool, pattern string) string {
	if optional {
		return fmt.Sprintf("pattern_optional=%s", pattern)
	}
	return fmt.Sprintf("pattern=%s", pattern)
}

// Enum returns a string for the enum validator tag
func Enum(optional bool) string {
	if optional {
		return "enum_optional"
	}
	return "enum"
}

// ArrayString returns a string for the array string validator tag
func ArrayString(optional bool, minLength, maxLength int) string {
	if optional {
		return fmt.Sprintf("array_string_optional=%d,%d", minLength, maxLength)
	}
	return fmt.Sprintf("array_string=%d,%d", minLength, maxLength)
}

// ArrayNumber returns a string for the array number validator tag
func ArrayNumber(optional bool, minLength, maxLength, minValue, maxValue int, canBeZero bool) string {
	if optional {
		return fmt.Sprintf("array_number_optional=%d,%d,%d,%d,%t", minLength, maxLength, minValue, maxValue, canBeZero)
	}
	return fmt.Sprintf("array_number=%d,%d,%d,%d,%t", minLength, maxLength, minValue, maxValue, canBeZero)
}

// EnumStringMapValidator validates a map with enum keys to string array values
func EnumStringMapValidator(fld validator.FieldLevel) bool {
	field := fld.Field()

	// Check if it's a map
	if field.Kind() != reflect.Map {
		return false
	}

	// Check if the map is not nil and not empty
	if field.IsNil() || field.Len() == 0 {
		return false
	}

	// Get all keys in the map
	iter := field.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		// Check if value is a slice of strings
		if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
			return false
		}

		// Validate key is an int/enum (non-zero)
		if !isValidEnumKey(key) {
			return false
		}

		// Value must not be empty
		if value.Len() == 0 {
			return false
		}

		// All items in the slice must be non-empty strings
		for i := 0; i < value.Len(); i++ {
			item := value.Index(i)
			if item.Kind() != reflect.String || item.String() == "" {
				return false
			}
		}
	}

	return true
}

// OptionalEnumStringMapValidator validates an optional map with enum keys to string array values
func OptionalEnumStringMapValidator(fld validator.FieldLevel) bool {
	field := fld.Field()

	// Check if it's a map
	if field.Kind() != reflect.Map {
		return true // Optional can be different type
	}

	// If nil or empty, it's valid for optional
	if field.IsNil() || field.Len() == 0 {
		return true
	}

	// Get all keys in the map
	iter := field.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		// Check if value is a slice of strings
		if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
			return false
		}

		// Validate key is an int/enum (non-zero)
		if !isValidEnumKey(key) {
			return false
		}

		// Value must not be empty
		if value.Len() == 0 {
			return false
		}

		// All items in the slice must be non-empty strings
		for i := 0; i < value.Len(); i++ {
			item := value.Index(i)
			if item.Kind() != reflect.String || item.String() == "" {
				return false
			}
		}
	}

	return true
}

// isValidEnumKey checks if a key is a valid enum (int type)
func isValidEnumKey(key reflect.Value) bool {
	switch key.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true // Any integer type is valid for enum
	default:
		return false
	}
}

// EnumStringMap returns a string for the enum string map validator tag
func EnumStringMap(optional bool) string {
	if optional {
		return "enum_string_map_optional"
	}
	return "enum_string_map"
}
