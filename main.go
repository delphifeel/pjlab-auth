package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

func testCreds(username string, password string) bool {
	f, err := os.Open("users.csv")
	if err != nil {
		panic(err)
	}

	passwordB64 := base64.StdEncoding.EncodeToString([]byte(password))
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		expectedUsername := record[0]
		expectedPasswordB64 := record[1]

		if expectedUsername == username && expectedPasswordB64 == passwordB64 {
			return true
		}
	}

	return false
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

	if ok := testCreds(username, password); !ok {
		http.Error(w, "Wrong creds", 400)
		return
	}

	token := makeToken(username, password)
	fmt.Fprintf(w, token)
}

func main() {
	http.HandleFunc("/login", loginHandler)

	log.Println("pjlab_auth service started.")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}
