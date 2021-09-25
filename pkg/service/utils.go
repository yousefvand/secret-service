package service

import (
	"fmt"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// UUID returns uuid without dashes
func UUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// Path2Name takes a dbus path like /a/b/c/d and returns
// a dbus object name like a.b.c and last part d
// Example: path=/a/b/c/xyz, name=Foo -> (a.b.c.Foo, xyz)
func Path2Name(path string, name string) (string, string) {

	path = strings.TrimSpace(path)
	name = strings.TrimSpace(name)

	result := strings.ReplaceAll(path, "/", ".")
	result = strings.TrimLeft(result, ".")
	idx := strings.LastIndex(result, ".")
	child := result[idx+1:]
	result = result[:idx]
	result += "." + name

	return result, child
}

// MemUsageOS returns memory usage on OS in MB
func MemUsageOS() uint64 {
	var memoryUsage runtime.MemStats
	runtime.ReadMemStats(&memoryUsage)
	return memoryUsage.Sys / 1024 / 1024
}

// IsMapSubsetSingleMatch returns true if only one
// key/value of mapSubset exists in mapSet otherwise false
func IsMapSubsetSingleMatch(mapSet map[string]string,
	mapSubset map[string]string, lock *sync.RWMutex) bool {

	lock.RLock()
	defer lock.RUnlock()

	if len(mapSubset) == 0 {
		return true
	}

	if len(mapSubset) > len(mapSet) {
		return false
	}

	for k, v := range mapSubset {
		if _, ok := mapSet[k]; ok {
			if mapSet[k] == v {
				return true
			}
		}
	}
	return false
}

// IsMapSubsetFullMatch returns true if mapSubset is a full subset of mapSet otherwise false
func IsMapSubsetFullMatch(mapSet map[string]string,
	mapSubset map[string]string, lock *sync.RWMutex) bool {

	lock.RLock()
	defer lock.RUnlock()

	if len(mapSubset) == 0 {
		return true
	}

	if len(mapSubset) > len(mapSet) {
		return false
	}

	for k, v := range mapSubset {
		if _, ok := mapSet[k]; !ok {
			return false
		}
		if mapSet[k] != v {
			return false
		}
	}
	return true
}

// IsMapSubsetFullMatchGeneric returns true if
// mapSubset is a full subset of mapSet otherwise false
func IsMapSubsetFullMatchGeneric(mapSet interface{},
	mapSubset interface{}, lock *sync.RWMutex) bool {

	lock.RLock()
	defer lock.RUnlock()

	mapSetValue := reflect.ValueOf(mapSet)
	mapSubsetValue := reflect.ValueOf(mapSubset)

	if fmt.Sprintf("%T", mapSet) != fmt.Sprintf("%T", mapSubset) {
		return false
	}

	if len(mapSetValue.MapKeys()) < len(mapSubsetValue.MapKeys()) {
		return false
	}

	if len(mapSubsetValue.MapKeys()) == 0 {
		return true
	}

	iterMapSubset := mapSubsetValue.MapRange()

	for iterMapSubset.Next() {
		k := iterMapSubset.Key()
		v := iterMapSubset.Value()

		value := mapSetValue.MapIndex(k)

		if !value.IsValid() || v.Interface() != value.Interface() {
			return false
		}
	}

	return true
}

/* TODO: Needs go 1.18

func IsMapSubsetFullMatchGeneric[K, V comparable](m, sub map[K]V) bool {
    if len(sub) > len(m) {
        return false
    }
    for k, vsub := range sub {
        if vm, found := m[k]; !found || vm != vsub {
            return false
        }
    }
    return true
}

*/

// CommandExists returns true if command exists on OS otherwise false
func CommandExists(cmdName string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+cmdName)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// Epoch returns seconds from midnight 1/1/1970
func Epoch() uint64 {
	return uint64(time.Now().Unix())
}
