package helpers

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateReferralCode(name string) string {
	rand.Seed(time.Now().UnixNano())

	base := strings.ToUpper(strings.ReplaceAll(name, " ", ""))
	if len(base) > 5 {
		base = base[:5]
	}

	random := rand.Intn(9000) + 1000
	return base + string(rune(random))
}