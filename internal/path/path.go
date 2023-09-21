package path

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _   = runtime.Caller(0)
	Root         = filepath.Join(filepath.Dir(b), "../..")
	ConfigPath   = filepath.Join(Root, "config.json")
	SQLSetupPath = filepath.Join(Root, "store", "setup.sql")
)
