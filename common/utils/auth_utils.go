package utils

import (
	"entry_task/config"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

func HashPwd(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash), err
}

func GenerateJWT(userId *uint, username *string, topicId *uint) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      strconv.FormatUint(uint64(*userId), 10),
		"username": username,
		"group":    strconv.FormatUint(uint64(*topicId), 10),
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.ProjectConfig.Jwt.Key))
	return "Bearer " + tokenString, err

}

func CheckPassword(password *string, hash *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*hash), []byte(*password))
	return err == nil
}

func ParseJWT(tokenString string) (uint, string, uint, error) {
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.ProjectConfig.Jwt.Key), nil
	})
	if err != nil {
		return 0, "", 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		groupRaw, ok := claims["group"].(string)
		if !ok {
			return 0, "", 0, errors.New("invalid group claim")
		}

		group, err := strconv.ParseUint(groupRaw, 10, 32)
		if err != nil {
			return 0, "", 0, errors.New("invalid group claim2")
		}

		idRaw, ok := claims["sub"].(string)
		if !ok {
			return 0, "", 0, errors.New("invalid id claim")
		}

		id, err := strconv.ParseUint(idRaw, 10, 32)
		if err != nil {
			return 0, "", 0, errors.New("invalid id claim2")
		}

		username, ok := claims["username"].(string)
		if !ok {
			return 0, "", 0, errors.New("invalid username claim")
		}

		return uint(id), username, uint(group), nil
	}

	return 0, "", 0, errors.New("token is not valid")
}
