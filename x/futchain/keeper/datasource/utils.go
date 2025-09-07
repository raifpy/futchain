package datasource

import (
	"crypto/md5"
	"fmt"
	"strings"
)

func calcHash(str string) string {
	hash := md5.Sum([]byte(str))
	return strings.ToUpper(fmt.Sprintf("%x", hash))
}
