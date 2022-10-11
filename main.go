package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/delphifeel/pjlab_auth/db"
	"log"
	"net"
	"net/http"
	"time"
	"sync"
)

var loginMutex sync.RWMutex

func makeToken(username string, password string) string {
	now := time.Now().UnixMicro()
	s := fmt.Sprintf("%v:%v:%v", username, password, now)
	sum := sha256.Sum256([]byte(s))
	encoded := base64.URLEncoding.EncodeToString(sum[:])
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

	if ok := db.TestCreds(&loginMutex, username, password); !ok {
		http.Error(w, "Wrong creds", 400)
		return
	}

	token := makeToken(username, password)
	fmt.Fprintf(w, token)

	loginMutex.Lock()
	db.ChangeUserToken(username, token)
	loginMutex.Unlock()
}

func testTokenHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Connection", "close")
	
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		http.Error(w, "Error", 500)
		log.Fatal(err)
		return
	}
	// allow only localhost (for now)
	if host != "127.0.0.1" {
		http.Error(w, "Denied", 403)
		return
	}

	username := req.FormValue("username")
	if username == "" {
		http.Error(w, "No username", 400)
		return
	}
	tokenB64 := req.FormValue("tokenB64")
	if tokenB64 == "" {
		http.Error(w, "No token", 400)
		return
	}

	if ok := db.TestUserToken(username, tokenB64); !ok {
		fmt.Fprintf(w, "No")
		return
	}

	fmt.Fprintf(w, "Yes")
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/testToken", testTokenHandler)

	log.Println("pjlab_auth service started.")
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}
