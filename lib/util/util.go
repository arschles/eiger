package util

import (
	"fmt"
	"log"
)

func LogWarnf(fmtstr string, i ...interface{}) {
	s := fmt.Sprintf(fmtstr, i...)
	log.Printf("[WARN] %s", s)
}
