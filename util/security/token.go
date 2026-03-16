package security

import (
	"errors"
	"fmt"
	"log"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/noti"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
)

const (
	normal_action_duration time.Duration = time.Minute * 30    // 30'
	access_duration        time.Duration = time.Hour * 24      // 1 ngày
	refresh_duration       time.Duration = access_duration * 7 // 1 tuần
)

func GenerateActionToken(address, sub, role string, logger *log.Logger) (string, int64, error) {
	var bytes = []byte(os.Getenv(env.SECRET_KEY))
	var errMsg string = "Error while generating token - "

	var exp int64 = time.Now().Add(normal_action_duration).Unix()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"address": address,
		"sub":     sub,
		"role":    role,
		"expire":  exp,
	}).SignedString(bytes)
	if err != nil {
		logger.Print(errMsg + fmt.Sprint(err))
		return "", 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return token, exp, nil
}

func GenerateActionTokenV2(address, sub string, roles []string, logger *log.Logger) (string, int64, error) {
	var bytes = []byte(os.Getenv(env.SECRET_KEY))
	var errMsg string = "Error while generating token - "

	var exp int64 = time.Now().Add(normal_action_duration).Unix()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"address": address,
		"sub":     sub,
		"roles":   roles,
		"expire":  exp,
	}).SignedString(bytes)
	if err != nil {
		logger.Print(errMsg + fmt.Sprint(err))
		return "", 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return token, exp, nil
}

func ExtractDataFromToken(tokenString string, logger *log.Logger) (string, string, string, time.Time, error) {
	var errRes error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	var errLogMsg string = "Error at ExtractDataFromToken - "

	// Check for empty token
	if tokenString == "" {
		logger.Println(errLogMsg + "empty token")
		return "", "", "", time.Time{}, errRes
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimSpace(tokenString[7:])
	}

	// Trim any whitespace and non-printable characters
	tokenString = strings.TrimSpace(tokenString)

	// Check if token contains invalid characters
	for i, c := range tokenString {
		if !unicode.IsPrint(c) || c == ' ' {
			logger.Printf(errLogMsg+"invalid character at position %d: %q\n", i, c)
			return "", "", "", time.Time{}, errRes
		}
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secretKey := os.Getenv(env.SECRET_KEY)
		if secretKey == "" {
			logger.Println(errLogMsg + "secret key not found in environment")
			return nil, errors.New("missing secret key")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		logger.Println(errLogMsg + err.Error())
		return "", "", "", time.Time{}, errRes
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		logger.Println(errLogMsg + "invalid claims or token")
		return "", "", "", time.Time{}, errRes
	}

	// Extract address
	address, ok := claims["address"].(string)
	if !ok || address == "" {
		logger.Println(errLogMsg + "missing or invalid address claim")
		return "", "", "", time.Time{}, errRes
	}

	// Extract sub
	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		logger.Println(errLogMsg + "missing or invalid sub claim")
		return "", "", "", time.Time{}, errRes
	}

	// Extract role
	role, ok := claims["role"].(string)
	if !ok || sub == "" {
		logger.Println(errLogMsg + "missing or invalid role claim")
		return "", "", "", time.Time{}, errRes
	}

	// Extract expiration
	expFloat, ok := claims["expire"].(float64)
	if !ok {
		logger.Println(errLogMsg + "missing or invalid expiration claim")
		return "", "", "", time.Time{}, errRes
	}

	// Convert Unix timestamp to time.Time
	exp := time.Unix(int64(expFloat), 0)

	return address, sub, role, exp, nil
}

func ExtractDataFromTokenV2(tokenString string, logger *log.Logger) (string, string, []string, time.Time, error) {
	var errRes error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	var errLogMsg string = "Error at ExtractDataFromToken - "

	// Check for empty token
	if tokenString == "" {
		logger.Println(errLogMsg + "empty token")
		return "", "", nil, time.Time{}, errRes
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimSpace(tokenString[7:])
	}

	// Trim any whitespace and non-printable characters
	tokenString = strings.TrimSpace(tokenString)

	// Check if token contains invalid characters
	for i, c := range tokenString {
		if !unicode.IsPrint(c) || c == ' ' {
			logger.Printf(errLogMsg+"invalid character at position %d: %q\n", i, c)
			return "", "", nil, time.Time{}, errRes
		}
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secretKey := os.Getenv(env.SECRET_KEY)
		if secretKey == "" {
			logger.Println(errLogMsg + "secret key not found in environment")
			return nil, errors.New("missing secret key")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		logger.Println(errLogMsg + err.Error())
		return "", "", nil, time.Time{}, errRes
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		logger.Println(errLogMsg + "invalid claims or token")
		return "", "", nil, time.Time{}, errRes
	}

	// Extract address
	address, ok := claims["address"].(string)
	if !ok || address == "" {
		logger.Println(errLogMsg + "missing or invalid address claim")
		return "", "", nil, time.Time{}, errRes
	}

	// Extract sub
	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		logger.Println(errLogMsg + "missing or invalid sub claim")
		return "", "", nil, time.Time{}, errRes
	}

	// Extract roles
	rawRoles, ok := claims["roles"].([]interface{})
	if !ok || len(rawRoles) == 0 {
		logger.Println(errLogMsg + "missing or invalid roles claim")
		return "", "", nil, time.Time{}, errRes
	}

	var roles []string
	for _, r := range rawRoles {
		if s, ok := r.(string); ok && s != "" {
			roles = append(roles, s)
		}
	}

	// Extract expiration
	expFloat, ok := claims["expire"].(float64)
	if !ok {
		logger.Println(errLogMsg + "missing or invalid expiration claim")
		return "", "", nil, time.Time{}, errRes
	}

	// Convert Unix timestamp to time.Time
	exp := time.Unix(int64(expFloat), 0)

	return address, sub, roles, exp, nil
}

// func GenerateActionTokenV2(address, sub, roles []string, logger *log.Logger) (string, int64, error)
// token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"address": address,
// 		"sub":     sub,
// 		"roles":    roles, // ← plural, []string stored as JSON array
// 		"expire":  exp,
// 	}).SignedString(bytes)
// 	if err != nil {
// 		logger.Print(errMsg + fmt.Sprint(err))
// 		return "", 0, errors.New(noti.INTERNALL_ERR_MSG)
// 	}

// ---------------------------------------------

// func ExtractDataFromTokenV2(tokenString string, logger *log.Logger) (string, string, []string, time.Time, error) {
// // Extract roles
// 	// JWT stores JSON arrays as []interface{} — convert to []string
// 	rawRoles, ok := claims["roles"].([]interface{})
// 	if !ok || len(rawRoles) == 0 {
// 		logger.Println(errLogMsg + "missing or invalid roles claim")
// 		return "", "", nil, time.Time{}, errRes
// 	}
// 	var roles []string
// 	for _, r := range rawRoles {
// 		if s, ok := r.(string); ok && s != "" {
// 			roles = append(roles, s)
// 		}
// 	}
