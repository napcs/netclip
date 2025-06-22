package netclip_test

import (
	"testing"

	"netclip"

	"github.com/stretchr/testify/assert"
)

func TestStoreAndGet(t *testing.T) {
	ds := netclip.NewDataStore()
	ds.Store("foo", "bar")
	value, ok := ds.GetValue("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", value)
}

func TestDelete(t *testing.T) {
	ds := netclip.NewDataStore()
	ds.Store("foo", "bar")
	ds.Delete("foo")
	_, ok := ds.GetValue("foo")
	assert.False(t, ok)
}

func TestRange(t *testing.T) {
	ds := netclip.NewDataStore()
	ds.Store("foo", "bar")
	ds.Store("baz", "qux")
	data := ds.Range()
	assert.Equal(t, 2, len(data))
	assert.Equal(t, "qux", data["baz"])
	assert.Equal(t, "bar", data["foo"])
}
