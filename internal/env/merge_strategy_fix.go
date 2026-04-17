package env

import "fmt"

// fmtErrorf is an alias kept here so the strategy file can call fmt.Errorf
// without a separate import block (the real import lives in this file).
var _ = fmt.Errorf // ensure import is used
