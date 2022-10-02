package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func makeToken(username string, password string) string {
	now := time.Now()
	year, month, day := now.Date()
	s := fmt.Sprintf("%s:%s:%d:%d:%d", username, password, year, month, day)
	sum := sha256.Sum256([]byte(s))
	encoded := base64.StdEncoding.EncodeToString(sum[:])
	return encoded
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Connection", "close")

	username := req.FormValue("username")
	if username == "" {
		http.Error(w, "No username", 400)
		return
	}
	password := req.FormValue("password")
	if password == "" {
		http.Error(w, "No password", 400)
		return
	}

	token := makeToken(username, password)

	fmt.Fprintf(w, token)
}

func main() {
	http.HandleFunc("/login", loginHandler)

	fmt.Println("pjlab_auth service started.")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		fmt.Println(err)
	}
}
