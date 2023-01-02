package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func ToBase64(path string) string {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	res := base64.StdEncoding.EncodeToString(fileBytes)
	return res
}

func main() {
	input := os.Args[1]

	strs := strings.Split(input, "\\")
	name := strings.Split(strs[len(strs)-1], ".")[0]
	fmt.Println("var " + name + " = \"" + ToBase64(input) + "\"")
}
