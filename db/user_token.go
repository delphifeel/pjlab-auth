package db

import (
	"io"
	"log"
)

const tokensFile = "tokens.csv"

func ChangeUserToken(username string, tokenB64 string) {
	reader, closeFile := newDBReader(tokensFile)
	allRecords, err := reader.ReadAll()
	closeFile()
	if err != nil {
		log.Fatal(err)
	}
	index := -1
	for i, record := range allRecords {
		if record[0] == username {
			index = i
			break
		}
	}

	if index != -1 {
		allRecords[index][0] = username
		allRecords[index][1] = tokenB64
	} else {
		allRecords = append(allRecords, []string {username, tokenB64})
	}

	writer, closeFile := newDBWriter(tokensFile)
	defer closeFile()
	if err := writer.WriteAll(allRecords); err != nil {
		log.Fatal(err)
	}
}

func TestUserToken(username string, tokenB64 string) bool {
	reader, closeFile := newDBReader(tokensFile)
	defer closeFile()

	// TODO: same in test_creds.go
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		expectedUsername := record[0]
		expectedToken := record[1]

		if expectedUsername == username {
			return tokenB64 == expectedToken
		}
	}

	return false
}
