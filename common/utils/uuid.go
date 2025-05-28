package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func UuidHex() string {
	uuidWithHyphen := uuid.New()
	u := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return u
}

func UuidMd5() string {
	uuidWithHyphen := uuid.New()
	u := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return fmt.Sprintf("%x", md5.Sum([]byte(u)))
}
