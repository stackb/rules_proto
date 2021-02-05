package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

func mustGetSha256(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	h := sha256.New()
	if _, err = io.Copy(h, response.Body); err != nil {
		log.Fatal(err)
	}

	sha256sum := fmt.Sprintf("%x", h.Sum(nil))

	log.Printf("sha256 for %s is %q", url, sha256sum)

	return sha256sum
}

func stringInSlice(search string, slice []string) bool {
	for _, item := range slice {
		if item == search {
			return true
		}
	}
	return false
}
