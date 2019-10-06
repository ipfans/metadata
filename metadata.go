package metadata

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// MetaData is a mapping from metadata keys to values.
// Users should use the following two convenience functions
// New and Pairs to generate MetaData.
type MetaData struct {
	locker sync.RWMutex
	data   map[string]interface{}
}

// New creates an MD from a given key-value map.
//
// Only the following ASCII characters are allowed in keys:
//  - digits: 0-9
//  - uppercase letters: A-Z (normalized to lower)
//  - lowercase letters: a-z
//  - special characters: -_.
// Uppercase letters are automatically converted to lowercase.
func New(m map[string]interface{}) *MetaData {
	md := &MetaData{data: map[string]interface{}{}}
	for k, val := range m {
		key := strings.ToLower(k)
		md.data[key] = val
	}
	return md
}

// Pairs returns an MD formed by the mapping of key, value ...
// Pairs panics if len(kv) is odd.
//
// Only the following ASCII characters are allowed in keys:
//  - digits: 0-9
//  - uppercase letters: A-Z (normalized to lower)
//  - lowercase letters: a-z
//  - special characters: -_.
// Uppercase letters are automatically converted to lowercase.
func Pairs(kv ...interface{}) *MetaData {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: Pairs got the odd number of input pairs for metadata: %d", len(kv)))
	}
	md := &MetaData{data: map[string]interface{}{}}
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = strings.ToLower(s.(string))
			continue
		}
		md.data[key] = s
	}
	return md
}

// Len returns the number of items in metadata.
func (md *MetaData) Len() int {
	return len(md.data)
}

// Copy returns a copy of metadata.
func (md *MetaData) Copy() *MetaData {
	return Join(md)
}

// Merge from metadatas.
func (md *MetaData) Merge(mds ...*MetaData) {
	if len(mds) == 0 {
		return
	}
	for i := range mds {
		mds[i].locker.RLock()
		for k := range mds[i].data {
			md.Set(k, mds[i].data[k])
		}
		mds[i].locker.RUnlock()
	}
}

// Get obtains the values for a given key.
func (md *MetaData) Get(k string) interface{} {
	k = strings.ToLower(k)
	md.locker.RLock()
	v := md.data[k]
	md.locker.RUnlock()
	return v
}

// Set sets the value of a given key with value.
func (md *MetaData) Set(k string, val interface{}) {
	if val == nil {
		return
	}
	k = strings.ToLower(k)
	md.locker.Lock()
	md.data[k] = val
	md.locker.Unlock()
}

// Join joins any number of mds into a single MetaData.
// The order of values for each key is determined by the order in which
// the mds containing those values are presented to Join.
func Join(mds ...*MetaData) *MetaData {
	out := &MetaData{data: map[string]interface{}{}}
	for _, md := range mds {
		md.locker.RLock()
		for k := range md.data {
			out.data[k] = md.Get(k)
		}
		md.locker.RUnlock()
	}
	return out
}

type mdContextKey struct{}

// NewContext creates a new context with metadata attached.
func NewContext(ctx context.Context, md *MetaData) context.Context {
	return context.WithValue(ctx, mdContextKey{}, md)
}

// FromContext returns the incoming metadata in ctx if it exists. The
// returned Metadata can be modified.
func FromContext(ctx context.Context) (md *MetaData, ok bool) {
	md, ok = ctx.Value(mdContextKey{}).(*MetaData)
	return
}

// AppendToContext returns a new context with the provided kv merged
// with any existing metadata in the context.
func AppendToContext(ctx context.Context, mds ...*MetaData) context.Context {
	if len(mds) == 0 {
		return ctx
	}
	md, ok := ctx.Value(mdContextKey{}).(*MetaData)
	if !ok {
		md = &MetaData{data: map[string]interface{}{}}
	}
	newMDs := make([]*MetaData, len(mds)+1)
	newMDs[0] = md
	for i := range mds {
		newMDs[i+1] = mds[i]
	}
	newMD := Join(newMDs...)
	ctx = context.WithValue(ctx, mdContextKey{}, newMD)
	return ctx
}
