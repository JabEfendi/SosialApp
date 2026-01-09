package helpers

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateReferralCode(name string) string {
	rand.Seed(time.Now().UnixNano())

	prefix := strings.ToUpper(strings.ReplaceAll(name, " ", ""))
	if len(prefix) > 3 {
		prefix = prefix[:3]
	}

	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	suffix := ""
	for i := 0; i < 9; i++ {
		suffix += string(chars[rand.Intn(len(chars))])
	}

	return prefix + suffix
}