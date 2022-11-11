package weak

import (
	"sync"
	"unsafe"

	weakext "github.com/KarpelesLab/weak"
)

type valueRef[Key any, Value any] struct {
	*weakext.Ref[Key]
	value *Value
}

// Table is a table of Go objects as weak reference keys which map to Go objects
// and is useful for tying values transiently to the existence of another Go object
type Table[Key any, Value any] struct {
	table map[uintptr]valueRef[Key, Value]
	lock  sync.RWMutex
}

func (w *Table[Key, Value]) value(key uintptr) *Value {
	value, ok := w.tryGetValue(key)
	if ok {
		return value
	}

	// Key still present, but not valid
	w.deleteKey(key)
	return nil
}

// Get returns the value for a given weak reference pointer
func (w *Table[Key, Value]) Get(key *Key) *Value {
	w.initTable()
	if key == nil {
		return nil
	}

	ref := w.getRef(key)
	return w.value(ref)
}

// NewRef returns a reference to the object v that may be cleared by the garbage collector
func (w *Table[Key, Value]) GetOrCreate(key *Key, factory func(key *Key) *Value) *Value {
	w.initTable()
	if key == nil {
		return nil
	}

	value := w.Get(key)
	if value != nil {
		return value
	}

	keyRef := weakext.NewRefDestroyer(key, func(key *Key, wr *weakext.Ref[Key]) {
		w.deleteKey(w.getRef(key))
	})
	if keyRef.Get() == nil {
		return nil
	}

	value = factory(key)
	w.table[w.getRef(key)] = valueRef[Key, Value]{
		Ref:   keyRef,
		value: value,
	}

	return value
}

// Delete deletes the specified key
func (w *Table[Key, Value]) Delete(key *Key) {
	ref := w.getRef(key)
	w.deleteKey(ref)
}

// Size gets the number of elements in the table
func (w *Table[Key, Value]) Size() int {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return len(w.table)
}

func (w *Table[Key, Value]) deleteKey(key uintptr) {
	w.lock.Lock()
	defer w.lock.Unlock()
	delete(w.table, key)
}

func (w *Table[Key, Value]) tryGetValue(key uintptr) (*Value, bool) {
	w.lock.RLock()
	defer w.lock.RUnlock()

	valueRef, ok := w.table[key]
	if !ok {
		return nil, true
	}

	// Make sure the key is still valid
	if valueRef.Get() != nil {
		return valueRef.value, true
	}

	return nil, false
}

func (w *Table[Key, Value]) initTable() {
	if w.table != nil {
		return
	}

	w.table = map[uintptr]valueRef[Key, Value]{}
}

func (w *Table[Key, Value]) getRef(key *Key) uintptr {
	return uintptr(unsafe.Pointer(key))
}
