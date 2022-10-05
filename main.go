package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/delphifeel/pjlab_auth/db"
	"log"
	"net/http"
	"time"
	"sync"
)

var mutex sync.RWMutex

func makeToken(username string, password string) string {
	now := time.Now().UnixMicro()
	s := fmt.Sprintf("%v:%v:%v", username, password, now)
	fmt.Println(s)
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

	if ok := db.TestCreds(&mutex, username, password); !ok {
		http.Error(w, "Wrong creds", 400)
		return
	}

	token := makeToken(username, password)
	fmt.Fprintf(w, token)

	mutex.Lock()
	db.ChangeUserToken(username, token)
	mutex.Unlock()
}

func initLogger() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	initLogger()

	http.HandleFunc("/login", loginHandler)

	log.Println("pjlab_auth service started.")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}
