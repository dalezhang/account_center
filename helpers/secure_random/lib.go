package secure_random

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func NewString(l int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, l)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

func RandomNumber(l int) int64 {
	seed := rand.NewSource(time.Now().Unix())
	r := rand.New(seed)
	rendint := r.Int63()
	str := fmt.Sprintf("%d", rendint)
	bb := []byte(str)
	max := len(bb)
	bb2 := bb[(max - l) : max-1]
	result, _ := strconv.ParseInt(string(bb2), 10, 64)
	return result
}
