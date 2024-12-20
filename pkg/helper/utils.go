package helper

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"regexp"
	"strings"
)

func OnlyNumbers(data string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(data, "")
}

func GenerateStrongPassword(length int) string {
	if length == 0 {
		length = 20
	}
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[{]}\\|;:,<.>/?"
	var password []byte
	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password = append(password, charset[randomIndex.Int64()])
	}
	return string(password)
}

func GetIP(r *http.Request) string {
	// Try to get the IP address from the X-Forwarded-For header
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		// If there is no X-Forwarded-For header, get the IP address from the request's remote address
		ip = r.RemoteAddr
	} else {
		// If there are multiple IPs in the X-Forwarded-For header, take the first one
		ip = strings.Split(ip, ",")[0]
	}
	// Trim any port information from the IP address
	ip = strings.Split(ip, ":")[0]
	return ip
}
