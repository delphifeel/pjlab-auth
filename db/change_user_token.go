package db

import "log"

const tokensFile = "tokens.csv"

func ChangeUserToken(username string, token string) {
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
		allRecords[index][1] = token
	} else {
		allRecords = append(allRecords, []string {username, token})
	}

	writer, closeFile := newDBWriter(tokensFile)
	defer closeFile()
	if err := writer.WriteAll(allRecords); err != nil {
		log.Fatal(err)
	}
}
