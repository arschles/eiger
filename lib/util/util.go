package util

import (
    "log"
    "fmt"
)

func LogWarnf(fmtstr string, i ...interface{}) {
    s := fmt.Sprintf(fmtstr, i...)
    log.Printf("[WARN] %s", s)
}
