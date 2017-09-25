package tools

import (
	"math/rand"
	"sync"
	"time"
)

// Source: https://play.golang.org/p/WIgH7GRnN1
// 		   https://play.golang.org/p/HOZ8ox1E6P
// Reference: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var mutex sync.Mutex

func int63(src rand.Source) int64 {
	mutex.Lock()
	v := src.Int63()
	mutex.Unlock()
	return v
}

// Random randomizes a string with specific length
func Random(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits,
	// enough for letterIdxMax characters!
	for i, cache, remain := n-1, int63(src), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = int63(src), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
