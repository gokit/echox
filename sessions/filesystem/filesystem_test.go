package filesystem

import (
	"testing"

	"github.com/gokit/echox/sessions"
	"github.com/gokit/echox/sessions/tester"
)

var newStore = func(_ *testing.T) sessions.Store {
	store := NewStore("./runtime/sessions/", []byte("secret"))
	return store
}

func TestFileSystem_SessionGetSet(t *testing.T) {
	tester.GetSet(t, newStore)
}

func TestFileSystem_SessionDeleteKey(t *testing.T) {
	tester.DeleteKey(t, newStore)
}

func TestFileSystem_SessionFlashes(t *testing.T) {
	tester.Flashes(t, newStore)
}

func TestFileSystem_SessionClear(t *testing.T) {
	tester.Clear(t, newStore)
}

func TestFileSystem_SessionOptions(t *testing.T) {
	tester.Options(t, newStore)
}
