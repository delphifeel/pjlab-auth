package main

import (
	"fmt"
	"net/http"
)

func loginHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello Mr. Someone\n")
}

func main() {
	http.HandleFunc("/login", loginHandler)

	fmt.Println("pjlab_auth service started.")
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println(err)
	}
}
