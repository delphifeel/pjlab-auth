package db

import (
	"encoding/base64"
	"io"
	"log"
	"sync"
	"crypto/sha256"
)

const usersFile = "users.csv"

func TestCreds(mutex *sync.RWMutex, username string, password string) bool {
	mutex.RLock()
	defer mutex.RUnlock()

	usersReader, closeFile := newDBReader(usersFile)
	defer closeFile()
	sum := sha256.Sum256([]byte(password))
	passwordB64 := base64.URLEncoding.EncodeToString(sum[:])

	for {
		record, err := usersReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		expectedUsername := record[0]
		expectedPasswordB64 := record[1]

		if expectedUsername == username {
			return expectedPasswordB64 == passwordB64
		}
	}

	return false
}
