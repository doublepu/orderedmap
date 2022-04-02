package orderedmap

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
)

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		m: map[K]*list.Element{},
		l: list.New(),
	}
}

type kv[K comparable, V any] struct {
	K K
	V V
}

type OrderedMap[K comparable, V any] struct {
	m map[K]*list.Element
	l *list.List
}

func (me *OrderedMap[K, V]) Map() map[K]V {
	r := map[K]V{}
	for k, e := range me.m {
		r[k] = e.Value.(kv[K, V]).V
	}
	return r
}

func (me *OrderedMap[K, V]) List() []kv[K, V] {
	l := make([]kv[K, V], 0, me.l.Len())
	for v := me.l.Front(); v != nil; v = v.Next() {
		l = append(l, v.Value.(kv[K, V]))
	}
	return l
}

func (me *OrderedMap[K, V]) Set(k K, v V) {
	if _, ok := me.m[k]; ok {
		return
	}
	e := me.l.PushBack(kv[K, V]{
		K: k,
		V: v,
	})
	me.m[k] = e
}

func (me *OrderedMap[K, V]) Get(k K) (v V, ok bool) {
	e, ok := me.m[k]
	if ok {
		return e.Value.(kv[K, V]).V, true
	}
	return v, false
}

func (me *OrderedMap[K, V]) Delete(k K) {
	if e, ok := me.m[k]; ok {
		delete(me.m, k)
		me.l.Remove(e)
	}
}

func (me *OrderedMap[K, V]) Reset() {
	me.m = map[K]*list.Element{}
	me.l = list.New()
}

func (me *OrderedMap[K, V]) Len() int {
	return len(me.m)
}

func (me *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	kvs := me.List()
	l := len(kvs)
	if l == 0 {
		return []byte("{}"), nil
	}
	buf := bytes.NewBuffer(nil)
	buf.WriteString("{")
	for i, kv := range kvs {
		kBytes, err := json.Marshal(kv.K)
		if err != nil {
			return nil, err

		}
		buf.Write(kBytes)
		buf.WriteString(":")
		vBytes, err := json.Marshal(kv.V)
		if err != nil {
			return nil, err
		}
		buf.Write(vBytes)
		if i < l-1 {
			buf.Write([]byte(","))
		}
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (me *OrderedMap[K, V]) UnmarshalJSON(b []byte) error {
	tmp := map[string]V{}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	objectKeys, err := me.objectKeys(b)
	if err != nil {
		return err
	}
	me.Reset()
	for _, objectKey := range objectKeys {
		var k K
		err := json.Unmarshal([]byte(fmt.Sprintf(`"%v"`, objectKey)), &k)
		if err != nil {
			return err
		}
		me.Set(k, tmp[objectKey])
	}
	return nil
}

func (me *OrderedMap[K, V]) objectKeys(b []byte) ([]string, error) {
	d := json.NewDecoder(bytes.NewReader(b))
	t, err := d.Token()
	if err != nil {
		return nil, err
	}
	if t != json.Delim('{') {
		return nil, errors.New("expected start of object")
	}
	var keys []string
	for {
		t, err := d.Token()
		if err != nil {
			return nil, err
		}
		if t == json.Delim('}') {
			return keys, nil
		}
		keys = append(keys, t.(string))
		if err := me.skipValue(d); err != nil {
			return nil, err
		}
	}
}
func (me *OrderedMap[K, V]) skipValue(d *json.Decoder) error {
	t, err := d.Token()
	if err != nil {
		return err
	}
	switch t {
	case json.Delim('['), json.Delim('{'):
		for {
			if err := me.skipValue(d); err != nil {
				if err == end {
					break
				}
				return err
			}
		}
	case json.Delim(']'), json.Delim('}'):
		return end
	}
	return nil
}

var end = errors.New("invalid end of array or object")
