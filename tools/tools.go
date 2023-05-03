//go:build tools
// +build tools

// Place any runtime dependencies as imports in this file.
// Go modules will be forced to download and install them.
// https://github.com/google/wire/issues/299#issuecomment-908416248
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package tools

import (
	_ "github.com/google/wire/cmd/wire"
)
