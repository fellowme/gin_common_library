package util

import (
	"math/rand"
	"time"
)

var (
	byteArray = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	bytesLen  = 61
)

func RandomString(length int) string {
	result := make([]byte, 0)
	for i := 0; i < length; i++ {
		result = append(result, byteArray[rand.Intn(bytesLen)])
	}
	return string(result)
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max-min) + min
	return randNum
}
