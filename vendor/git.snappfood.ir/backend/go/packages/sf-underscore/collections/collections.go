package collections

import (
	"math/rand"
	"reflect"
	"sort"
	"time"
)

// Each iterates over a collection and calls the iteratee function for each element
func Each[T any](collection interface{}, iteratee func(item T, index int, collection interface{}) bool) {
	val := reflect.ValueOf(collection)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			if !iteratee(item, i, collection) {
				break
			}
		}
	} else if val.Kind() == reflect.Map {
		for i, key := range val.MapKeys() {
			item := val.MapIndex(key).Interface().(T)
			if !iteratee(item, i, collection) {
				break
			}
		}
	}
}

// Map creates a new slice with the results of calling the iteratee function on each element
func Map[T any, R any](collection interface{}, iteratee func(item T, index int, collection interface{}) R) []R {
	val := reflect.ValueOf(collection)
	result := make([]R, 0)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			result = append(result, iteratee(item, i, collection))
		}
	} else if val.Kind() == reflect.Map {
		for i, key := range val.MapKeys() {
			item := val.MapIndex(key).Interface().(T)
			result = append(result, iteratee(item, i, collection))
		}
	}

	return result
}

// Reduce boils down a collection to a single value
func Reduce[T any, R any](collection interface{}, iteratee func(result R, item T, index int, collection interface{}) R, initial R) R {
	val := reflect.ValueOf(collection)
	result := initial

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			result = iteratee(result, item, i, collection)
		}
	} else if val.Kind() == reflect.Map {
		for i, key := range val.MapKeys() {
			item := val.MapIndex(key).Interface().(T)
			result = iteratee(result, item, i, collection)
		}
	}

	return result
}

// ReduceRight is like Reduce but iterates from right to left
func ReduceRight[T any, R any](collection interface{}, iteratee func(result R, item T, index int, collection interface{}) R, initial R) R {
	val := reflect.ValueOf(collection)
	result := initial

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := val.Len() - 1; i >= 0; i-- {
			item := val.Index(i).Interface().(T)
			result = iteratee(result, item, i, collection)
		}
	} else if val.Kind() == reflect.Map {
		keys := val.MapKeys()
		for i := len(keys) - 1; i >= 0; i-- {
			item := val.MapIndex(keys[i]).Interface().(T)
			result = iteratee(result, item, i, collection)
		}
	}

	return result
}

// Find returns the first element that passes the predicate test
func Find[T any](collection interface{}, predicate func(item T, index int, collection interface{}) bool) (T, bool) {
	val := reflect.ValueOf(collection)
	var zero T

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			if predicate(item, i, collection) {
				return item, true
			}
		}
	} else if val.Kind() == reflect.Map {
		for i, key := range val.MapKeys() {
			item := val.MapIndex(key).Interface().(T)
			if predicate(item, i, collection) {
				return item, true
			}
		}
	}

	return zero, false
}

// Filter returns all elements that pass the predicate test
func Filter[T any](collection interface{}, predicate func(item T, index int, collection interface{}) bool) []T {
	val := reflect.ValueOf(collection)
	result := make([]T, 0)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			if predicate(item, i, collection) {
				result = append(result, item)
			}
		}
	} else if val.Kind() == reflect.Map {
		for i, key := range val.MapKeys() {
			item := val.MapIndex(key).Interface().(T)
			if predicate(item, i, collection) {
				result = append(result, item)
			}
		}
	}

	return result
}

// Where returns all elements that match the properties
func Where[T any](collection interface{}, properties map[string]interface{}) []T {
	return Filter(collection, func(item T, _ int, _ interface{}) bool {
		itemVal := reflect.ValueOf(item)
		if itemVal.Kind() == reflect.Struct {
			for key, val := range properties {
				field := itemVal.FieldByName(key)
				if !field.IsValid() || field.Interface() != val {
					return false
				}
			}
			return true
		}
		return false
	})
}

// FindWhere returns the first element that matches the properties
func FindWhere[T any](collection interface{}, properties map[string]interface{}) (T, bool) {
	return Find(collection, func(item T, _ int, _ interface{}) bool {
		itemVal := reflect.ValueOf(item)
		if itemVal.Kind() == reflect.Struct {
			for key, val := range properties {
				field := itemVal.FieldByName(key)
				if !field.IsValid() || field.Interface() != val {
					return false
				}
			}
			return true
		}
		return false
	})
}

// Reject returns all elements that don't pass the predicate test
func Reject[T any](collection interface{}, predicate func(item T, index int, collection interface{}) bool) []T {
	return Filter(collection, func(item T, index int, collection interface{}) bool {
		return !predicate(item, index, collection)
	})
}

// Every returns true if all elements pass the predicate test
func Every[T any](collection interface{}, predicate func(item T, index int, collection interface{}) bool) bool {
	val := reflect.ValueOf(collection)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			if !predicate(item, i, collection) {
				return false
			}
		}
	} else if val.Kind() == reflect.Map {
		for i, key := range val.MapKeys() {
			item := val.MapIndex(key).Interface().(T)
			if !predicate(item, i, collection) {
				return false
			}
		}
	}

	return true
}

// Some returns true if any element passes the predicate test
func Some[T any](collection interface{}, predicate func(item T, index int, collection interface{}) bool) bool {
	val := reflect.ValueOf(collection)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			if predicate(item, i, collection) {
				return true
			}
		}
	} else if val.Kind() == reflect.Map {
		for i, key := range val.MapKeys() {
			item := val.MapIndex(key).Interface().(T)
			if predicate(item, i, collection) {
				return true
			}
		}
	}

	return false
}

// Contains returns true if the value is present in the collection
func Contains[T comparable](collection interface{}, value T) bool {
	val := reflect.ValueOf(collection)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface()
			if item, ok := item.(T); ok && item == value {
				return true
			}
		}
	} else if val.Kind() == reflect.Map {
		for _, key := range val.MapKeys() {
			item := val.MapIndex(key).Interface()
			if item, ok := item.(T); ok && item == value {
				return true
			}
		}
	}

	return false
}

// Invoke calls the method named by methodName on each element in the collection
func Invoke(collection interface{}, methodName string, args ...interface{}) []interface{} {
	val := reflect.ValueOf(collection)
	result := make([]interface{}, 0)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			method := item.MethodByName(methodName)
			if method.IsValid() {
				callArgs := make([]reflect.Value, len(args))
				for j, arg := range args {
					callArgs[j] = reflect.ValueOf(arg)
				}
				res := method.Call(callArgs)
				if len(res) > 0 {
					result = append(result, res[0].Interface())
				} else {
					result = append(result, nil)
				}
			}
		}
	}

	return result
}

// Pluck extracts a list of property values
func Pluck[T any](collection interface{}, propertyName string) []interface{} {
	val := reflect.ValueOf(collection)
	result := make([]interface{}, 0)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			if item.Kind() == reflect.Struct {
				field := item.FieldByName(propertyName)
				if field.IsValid() {
					result = append(result, field.Interface())
				}
			} else if item.Kind() == reflect.Map {
				field := item.MapIndex(reflect.ValueOf(propertyName))
				if field.IsValid() {
					result = append(result, field.Interface())
				}
			}
		}
	}

	return result
}

// Max returns the maximum value in the collection
func Max[T any](collection interface{}, iteratee func(item T) float64) (T, bool) {
	val := reflect.ValueOf(collection)
	var maxItem T
	var maxValue float64
	found := false

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			value := iteratee(item)
			if !found || value > maxValue {
				maxValue = value
				maxItem = item
				found = true
			}
		}
	}

	return maxItem, found
}

// Min returns the minimum value in the collection
func Min[T any](collection interface{}, iteratee func(item T) float64) (T, bool) {
	val := reflect.ValueOf(collection)
	var minItem T
	var minValue float64
	found := false

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			value := iteratee(item)
			if !found || value < minValue {
				minValue = value
				minItem = item
				found = true
			}
		}
	}

	return minItem, found
}

// SortBy returns a sorted copy of the collection
func SortBy[T any](collection interface{}, iteratee func(item T) float64) []T {
	val := reflect.ValueOf(collection)
	result := make([]T, 0)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			result = append(result, item)
		}

		sort.Slice(result, func(i, j int) bool {
			return iteratee(result[i]) < iteratee(result[j])
		})
	}

	return result
}

// GroupBy groups the collection by the result of the iteratee function
func GroupBy[T any, K comparable](collection interface{}, iteratee func(item T) K) map[K][]T {
	val := reflect.ValueOf(collection)
	result := make(map[K][]T)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			key := iteratee(item)
			result[key] = append(result[key], item)
		}
	}

	return result
}

// IndexBy indexes the collection by the result of the iteratee function
func IndexBy[T any, K comparable](collection interface{}, iteratee func(item T) K) map[K]T {
	val := reflect.ValueOf(collection)
	result := make(map[K]T)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			key := iteratee(item)
			result[key] = item
		}
	}

	return result
}

// CountBy counts the collection by the result of the iteratee function
func CountBy[T any, K comparable](collection interface{}, iteratee func(item T) K) map[K]int {
	val := reflect.ValueOf(collection)
	result := make(map[K]int)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			key := iteratee(item)
			result[key]++
		}
	}

	return result
}

// Shuffle returns a shuffled copy of the collection
func Shuffle[T any](collection interface{}) []T {
	val := reflect.ValueOf(collection)
	result := make([]T, 0)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			result = append(result, item)
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(result), func(i, j int) {
			result[i], result[j] = result[j], result[i]
		})
	}

	return result
}

// Sample returns n random elements from the collection
func Sample[T any](collection interface{}, n int) []T {
	val := reflect.ValueOf(collection)
	result := make([]T, 0)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		if n >= val.Len() {
			return Shuffle[T](collection)
		}
		// Create a copy of the collection
		items := make([]T, val.Len())
		for i := 0; i < val.Len(); i++ {
			items[i] = val.Index(i).Interface().(T)
		}

		// Shuffle and take the first n elements
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(items), func(i, j int) {
			items[i], items[j] = items[j], items[i]
		})

		return items[:n]
	}

	return result
}

// ToArray converts the collection to an array
func ToArray[T any](collection interface{}) []T {
	val := reflect.ValueOf(collection)
	result := make([]T, 0)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			result = append(result, item)
		}
	} else if val.Kind() == reflect.Map {
		for _, key := range val.MapKeys() {
			item := val.MapIndex(key).Interface().(T)
			result = append(result, item)
		}
	}

	return result
}

// Size returns the size of the collection
func Size(collection interface{}) int {
	val := reflect.ValueOf(collection)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array || val.Kind() == reflect.Map {
		return val.Len()
	}

	return 0
}

// Partition splits the collection into two arrays: one with elements that pass the predicate test and one with elements that don't
func Partition[T any](collection interface{}, predicate func(item T, index int, collection interface{}) bool) [][]T {
	val := reflect.ValueOf(collection)
	pass := make([]T, 0)
	fail := make([]T, 0)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface().(T)
			if predicate(item, i, collection) {
				pass = append(pass, item)
			} else {
				fail = append(fail, item)
			}
		}
	} else if val.Kind() == reflect.Map {
		for i, key := range val.MapKeys() {
			item := val.MapIndex(key).Interface().(T)
			if predicate(item, i, collection) {
				pass = append(pass, item)
			} else {
				fail = append(fail, item)
			}
		}
	}

	return [][]T{pass, fail}
}
