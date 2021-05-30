package lib

import (
	"context"
	"testing"

	"golang.org/x/net/webdav"
)

func TestFileSystem(t *testing.T) {
	f := FS{
		Scope: `/`,
		FS: func(scope string, options map[string]string) webdav.FileSystem {
			panic(`panic~~~~`)
			return webdav.Dir(scope)
		},
		Options: map[string]string{},
	}
	ctx := context.Background()
	_, err := f.Stat(ctx, `/test`)
	if err.Error() != `panic~~~~` {
		t.Fatal(`err.Error() should be panic~~~~`)
	}
}
