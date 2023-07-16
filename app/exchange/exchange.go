package exchange

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

const nbkzUrl = "https://nationalbank.kz/rss/rates_all.xml"

type Rate struct {
	Title   string
	Current string
}

type Item struct {
	Title       string `xml:"title"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
	Quant       string `xml:"quant"`
	Index       string `xml:"index"`
	Change      string `xml:"change"`
}

func parseXMLRate(data io.ReadCloser) ([]Item, error) {
	xmlData, err := io.ReadAll(data)
	if err != nil {
		log.Fatal(err)
	}

	var channel struct {
		Items []Item `xml:"channel>item"`
	}
	if err := xml.Unmarshal(xmlData, &channel); err != nil {
		log.Fatal(err)
	}
	return channel.Items, nil
}

func GetCurrentRate(cur string) (Rate, error) { //@todo

	resp, err := http.Get(nbkzUrl)
	if err != nil {
		log.Fatal(err) //@todo : return err
	}
	defer resp.Body.Close()

	items, err := parseXMLRate(resp.Body)
	fmt.Println()
	var chosenItem Item //@todo

	for _, item := range items {
		if item.Title == cur {
			chosenItem = item
			break
		}
	}
	result := Rate{
		Title:   chosenItem.Title,
		Current: chosenItem.Description,
	}
	return result, nil

}
