package main

import (
	"flag"
	"log"

	"github.com/lufia/taskfs/fs"
	"github.com/lufia/taskfs/github"
)

var (
	accessToken = flag.String("t", "", "access token")
	baseURL     = flag.String("url", "", "endpoint url")
)

func main() {
	flag.Parse()
	s, err := github.NewService(&github.Config{
		BaseURL: *baseURL,
		Token:   *accessToken,
	})
	if err != nil {
		log.Fatal(err)
	}
	root := fs.NewRoot()
	root.CreateService(s)
	if err := root.MountAndServe("x"); err != nil {
		log.Fatal(err)
	}
}
