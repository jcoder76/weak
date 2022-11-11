package weak

import (
	"reflect"
	"sync"
)

var store = map[reflect.Type]any{}
var storeLock sync.RWMutex

// Extend gets or adds a type extension for a given source reference
func Extend[TValue any, TTarget any](target *TTarget) *TValue {
	return ExtendFunc(target, func(target *TTarget) *TValue {
		return new(TValue)
	})
}

// ExtendFunc gets or adds a type extension for a given source reference with
// a type factory function that will be called to create new instances
func ExtendFunc[TValue any, TTarget any](target *TTarget, factory func(target *TTarget) *TValue) *TValue {
	if target == nil {
		panic("extend requires a valid target reference")
	}

	table := getTable[TValue](target)
	if table == nil {
		panic("extend cannot create a table for target")
	}

	value := table.GetOrCreate(target, factory)

	return value
}

func getTable[TValue any, TTarget any](src *TTarget) *Table[TTarget, TValue] {
	typeOf := reflect.TypeOf(src)
	if typeOf == nil {
		return nil
	}

	var table *Table[TTarget, TValue]
	storeLock.RLock()
	tableInterface, ok := store[typeOf]
	storeLock.RUnlock()
	if ok {
		table, ok = tableInterface.(*Table[TTarget, TValue])
	}
	if !ok {
		table = newTable[TValue](src)
	}

	return table
}

func newTable[Value any, Target any](src *Target) *Table[Target, Value] {
	table := &Table[Target, Value]{}
	storeLock.Lock()
	defer storeLock.Unlock()
	store[reflect.TypeOf(src)] = table
	return table
}
