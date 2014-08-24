package client

import (
    "strings"
)

func stripEndSlash(s string) string {
    return strings.TrimRight(s, "/")
}
