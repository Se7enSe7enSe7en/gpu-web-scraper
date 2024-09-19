package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"web-scraper/util"

	"github.com/gocolly/colly"
)

type Website struct {
	Url, ProductSelector, ProductNameSelector, ProductPriceSelector string
}

type Product struct {
	Name, Price, Url string
}

var websites []Website = []Website{
	{
		Url:                  "https://pcx.com.ph/collections/graphics-cards",
		ProductSelector:      ".t4s-product-info",
		ProductNameSelector:  ".t4s-product-title",
		ProductPriceSelector: ".t4s-product-price",
	},
}

func main() {
	var products []Product

	for _, website := range websites {
		collector := colly.NewCollector()

		collector.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting: ", r.URL)
		})

		collector.OnError(func(r *colly.Response, e error) {
			fmt.Println("Error: ", e)
		})

		collector.OnHTML(website.ProductSelector, func(element *colly.HTMLElement) {
			product := Product{}

			product.Name = element.ChildText(website.ProductNameSelector)
			product.Price = element.ChildText(website.ProductPriceSelector)

			url := element.ChildAttr("a", "href")
			if !util.IsUrl(url) {
				urlBackup := url
				// try adding base path
				url = website.Url + url
				if !util.IsUrl(url) {
					fmt.Println("url stll not valid even after adding base path: ", url)
					fmt.Println("reverting to: ", urlBackup)
					url = urlBackup
				}
			}
			product.Url = url

			products = append(products, product)
		})

		collector.OnScraped(func(r *colly.Response) {
			file, err := os.Create("products.csv")
			if err != nil {
				log.Fatalln("Failed to create output CSV file", err)
			}
			defer file.Close()

			writer := csv.NewWriter(file)

			headers := []string{
				"Name",
				"Price",
				"Url",
			}

			writer.Write(headers)

			for _, product := range products {
				record := []string{
					product.Name,
					product.Price,
					product.Url,
				}

				writer.Write(record)
			}
			defer writer.Flush()
		})

		collector.Visit(website.Url)
	}
}
