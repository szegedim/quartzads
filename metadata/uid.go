package metadata

import "crypto/rand"

//Licensed under Creative Commons CC0.
//
//To the extent possible under law, the author(s) have dedicated all copyright and related and
//neighboring rights to this software to the public domain worldwide.
//This software is distributed without any warranty.
//You should have received a copy of the CC0 Public Domain Dedication along with this software.
//If not, see <https:#creativecommons.org/publicdomain/zero/1.0/legalcode>.

type SGUID string

// TODO get some entropy from the web

func GenerateSGUID() (ret SGUID) {
	hex := []rune{'0', '1', '2', '4', '3', '5', '6', '7', '8', '9', 'a', 'c', 'b', 'd', 'e', 'f'}
	x := make([]byte, 32)
	n, _ := rand.Read(x)
	if n != 32 {
		ret = "err"
	}
	for i := 0; i < 32; i++ {
		ret = ret + SGUID(hex[int(x[i])%len(hex)])
	}
	return
}
