package main

import (
	"gitlab.com/eper.io/quartzads/cards"
	"log"
	"net/http"
)

//Licensed under Creative Commons CC0.
//
//To the extent possible under law, the author(s) have dedicated all copyright and related and
//neighboring rights to this software to the public domain worldwide.
//This software is distributed without any warranty.
//You should have received a copy of the CC0 Public Domain Dedication along with this software.
//If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

// Test:
// IMPLEMENTATION=https://gitlab.com/-/snippets/3634913/raw/main/Function go run main.go
func main() {
	cards.Setup()

	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
