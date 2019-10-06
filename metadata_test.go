package metadata

import (
	"context"
	"reflect"
	"testing"
)

func TestPairsMD(t *testing.T) {
	for _, test := range []struct {
		// input
		kv []interface{}
		// output
		md MetaData
	}{
		{[]interface{}{}, MetaData{data: map[string]interface{}{}}},
		{
			[]interface{}{"k1", "v1", "k1", "v2"},
			MetaData{
				data: map[string]interface{}{"k1": "v2"},
			},
		},
	} {
		md := Pairs(test.kv...)
		if !reflect.DeepEqual(md.data, test.md.data) {
			t.Fatalf("Pairs(%v) = %+v, want %+v", test.kv, md.data, test.md.data)
		}
	}
}

func TestCopy(t *testing.T) {
	const key, val = "key", "val"
	orig := Pairs(key, val)
	cpy := orig.Copy()
	if !reflect.DeepEqual(orig.data, cpy.data) {
		t.Errorf("copied value not equal to the original, got %+v, want %+v", cpy.data, orig.data)
	}
	orig.Set(key, "baz")
	if v := cpy.data[key]; v.(string) != val {
		t.Errorf("change in original should not affect copy, got %q, want %q", v, val)
	}
}

func TestMerge(t *testing.T) {
	const key, val, newVal = "key", "val", "newVal"
	orig := Pairs(key, val)
	orig.Merge(Pairs(key, newVal))
	if !reflect.DeepEqual(orig.Get(key), newVal) {
		t.Errorf("context's metadata is %+v, want %+v", orig.Get(key), newVal)
	}
}
func TestJoin(t *testing.T) {
	for _, test := range []struct {
		mds  []*MetaData
		want *MetaData
	}{
		{[]*MetaData{}, &MetaData{data: map[string]interface{}{}}},
		{[]*MetaData{Pairs("foo", "bar")}, Pairs("foo", "bar")},
		{[]*MetaData{Pairs("foo", "bar"), Pairs("foo", "baz")}, Pairs("foo", "baz")},
		{[]*MetaData{Pairs("foo", "bar"), Pairs("foo", "baz"), Pairs("zip", "zap")}, Pairs("foo", "baz", "zip", "zap")},
	} {
		md := Join(test.mds...)
		if !reflect.DeepEqual(md.data, test.want.data) {
			t.Errorf("context's metadata is %+v, want %+v", md.data, test.want.data)
		}
	}
}

func TestGet(t *testing.T) {
	for _, test := range []struct {
		md      *MetaData
		key     string
		wantVal interface{}
	}{
		{md: Pairs("My-Optional-Header", "42"), key: "My-Optional-Header", wantVal: "42"},
		{md: Pairs("Header", "42", "Header", "43", "Header", "44", "other", "1"), key: "HEADER", wantVal: "44"},
		{md: Pairs("HEADER", "10"), key: "HEADER", wantVal: "10"},
	} {
		vals := test.md.Get(test.key)
		if !reflect.DeepEqual(vals, test.wantVal) {
			t.Errorf("value of metadata %v is %v, want %v", test.key, vals, test.wantVal)
		}
	}
}

func TestSet(t *testing.T) {
	for _, test := range []struct {
		md     *MetaData
		setKey string
		setVal interface{}
		want   *MetaData
	}{
		{
			md:     Pairs("My-Optional-Header", "42", "other-key", "999"),
			setKey: "Other-Key",
			setVal: "1",
			want:   Pairs("my-optional-header", "42", "other-key", "1"),
		},
		{
			md:     Pairs("My-Optional-Header", "42"),
			setKey: "Other-Key",
			setVal: "3",
			want:   Pairs("my-optional-header", "42", "other-key", "3"),
		},
		{
			md:     Pairs("My-Optional-Header", "42"),
			setKey: "Other-Key",
			setVal: nil,
			want:   Pairs("my-optional-header", "42"),
		},
	} {
		test.md.Set(test.setKey, test.setVal)
		if !reflect.DeepEqual(test.md.data, test.want.data) {
			t.Errorf("value of metadata is %+v, want %+v", test.md.data, test.want.data)
		}
	}
}

func TestNewContext(t *testing.T) {
	ctx := context.TODO()
	md := Pairs("foo", "bar")
	ctx = NewContext(ctx, md)
	v := ctx.Value(mdContextKey{})
	if !reflect.DeepEqual(v.(*MetaData).data, md.data) {
		t.Errorf("value(%+v) not set to context.", md.data)
	}
}

func TestFromContext(t *testing.T) {
	ctx := context.TODO()
	const key, val, new = "foo", "bar", "baz"
	md := Pairs(key, val)
	ctx = NewContext(ctx, md)
	md, ok := FromContext(ctx)
	if !ok {
		t.Errorf("value(%+v) not set to context.", md.data)
	}
	if !reflect.DeepEqual(md.Get(key), val) {
		t.Errorf("md.Get(%s) is %v, want %v.", key, md.Get(key), val)
	}
	md.Set(key, new)
	md, ok = FromContext(ctx)
	if !reflect.DeepEqual(md.Get(key), new) {
		t.Errorf("md.Get(%s) is %v, want %v.", key, md.Get(key), new)
	}
}

func TestAppendToContext(t *testing.T) {
	ctx := context.TODO()
	const key, val, newVal = "foo", "bar", "baz"
	md := Pairs(key, val)
	ctx = NewContext(ctx, md)
	ctx = AppendToContext(ctx, Pairs(key, newVal))
	md, _ = FromContext(ctx)
	if !reflect.DeepEqual(md.Get(key), newVal) {
		t.Errorf("md.Get(%s) is %v, want %v.", key, md.Get(key), newVal)
	}
}
