package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lmas/feedloggr2"
	"github.com/mmcdole/gofeed"
)

const (
	testDir string = "testdata"
	//maxItems int    = 50
)

var (
	parser *gofeed.Parser
)

func main() {
	parser = gofeed.NewParser()

	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".feed" {
			continue
		}

		path := filepath.Join(testDir, f.Name())
		err := handleFeed(path)
		if err != nil {
			panic(err)
		}
	}
}

func handleFeed(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	feed, err := parser.Parse(f)
	if err != nil {
		return err
	}

	var items []feedloggr2.Item
	//for _, i := range feed.Items[:maxItems] {
	for _, i := range feed.Items {
		items = append(items, feedloggr2.Item{
			i.Title,
			i.Link,
		})
	}

	b, err := json.Marshal(items)
	if err != nil {
		return err
	}

	newPath := strings.Replace(path, ".feed", ".wanted", 1)
	err = ioutil.WriteFile(newPath, b, 0644)
	if err != nil {
		return err
	}

	log.Printf("Parsed %s, wrote %s", path, newPath)
	return nil
}
