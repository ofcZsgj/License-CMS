package license

import (
	"math/rand"
	"time"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const licenseLen = 632

func RandStringBytes() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, licenseLen)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
