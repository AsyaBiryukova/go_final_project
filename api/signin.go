package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	targetPassword = os.Getenv("TODO_PASSWORD")
)

// PostSigninHandler обрабатывает запросы к api/signin.
// При корректном вводе пароля, возвращает JSON {"token":JWT}. В случае ошибки возвращает JSON {"error":error}
func PostSigninHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeErr(err, w)
		return
	}

	var body map[string]string
	if err = json.Unmarshal(buf.Bytes(), &body); err != nil {
		writeErr(err, w)
		return
	}
	password := body["password"]
	if len(password) == 0 {
		writeErr(fmt.Errorf("пустая строка вместо password"), w)
		return
	} else if password != targetPassword {
		writeErr(fmt.Errorf("неправильный пароль"), w)
		return
	}

	claims := jwt.MapClaims{
		"password": sha256.Sum256([]byte(targetPassword)),
		"Exp":      1550946689,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString([]byte(targetPassword))
	if err != nil {
		writeErr(err, w)
		return
	}

	tokenResp := map[string]string{
		"token": signedToken,
	}
	resp, err := json.Marshal(tokenResp)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}
