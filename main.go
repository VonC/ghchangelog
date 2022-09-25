package main

import (
	"fmt"
	"ghchangelog/version"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {
	fmt.Println("ghchangelog")
	if len(os.Args) == 1 {
		fmt.Println("Expect title part or -v, --version, -version or version")
		os.Exit(1)
	}
	firstParam := strings.ToLower(os.Args[1])
	switch firstParam {
	case "-v", "--version", "-version", "version":
		fmt.Println(version.String())
		os.Exit(0)
	}

	c := colly.NewCollector(
		colly.AllowedDomains("github.blog"),
		colly.MaxDepth(0),
	)
	// Step 2. Perform some logic before REQUEST Is made
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting ", r.URL.String())
	})

	// Step 2.1. If error occurred during the request, handle it!
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL: ", r.Request.URL, " failed with response: ", r, "\nError: ", err)
	})

	if err := c.Visit("https://github.blog/changelog/"); err != nil {
		log.Fatal(err)
	}
	// Wait until threads are finished
	c.Wait()
}
