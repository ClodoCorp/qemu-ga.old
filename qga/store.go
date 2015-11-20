package qga

var store map[string]map[interface{}]interface{} = make(map[string]map[interface{}]interface{})

func StoreSet(ns string, k interface{}, v interface{}) {
	store[ns][k] = v
}

func StoreGet(ns string, k interface{}) (interface{}, bool) {
	m, ok := store[ns]
	if !ok {
		return nil, false
	}
	v, ok := m[k]
	return v, ok
}

func StoreDel(ns string, k interface{}) {
	m, ok := store[ns]
	if !ok {
		return
	}
	delete(m, k)
}
