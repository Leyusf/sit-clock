package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
)

func ToBase64(path string) string {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	res := base64.StdEncoding.EncodeToString(fileBytes)
	return res
}

//func main() {
//	var input string
//	fmt.Scanln(&input)
//	fmt.Print(ToBase64(input))
//}
