package helpers

import (
	"math/rand"
	"strings"
	"time"

	"gorm.io/gorm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

func GenerateRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())

	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)

	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}

	return string(code)
}

func GenerateUniqueCorporateReferral(db *gorm.DB) (string, error) {
	for {
		code := GenerateRandomCode(8)

		var count int64

		db.Table("corporates").
			Where("referral_code = ?", code).
			Count(&count)

		if count > 0 {
			continue
		}

		db.Table("users").
			Where("referral_code = ?", code).
			Count(&count)

		if count > 0 {
			continue
		}

		return code, nil
	}
}