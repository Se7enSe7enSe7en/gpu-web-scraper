package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"web-scraper/util"

	"github.com/gocolly/colly"
)

// initialize a data structure to keep the scraped data
type Product struct {
	Url, Image, Name, Price string
}

func main() {
	// // get url from terminal
	// args := os.Args
	// url := args[1]
	// // DEBUG
	// fmt.Println("url: ", url)

	// initialize the slice of structs that will contain the scraped data
	var products []Product

	// new instance of collector (Q1: what is a collector?)
	collector := colly.NewCollector(
	// colly.AllowedDomains("www.scrapingcourse.com"), // allowed domains
	)

	// called before an HTTP request is triggered
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// triggered when the scraper encounters an error
	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Crikey, an error occurred!: ", e)
	})

	// fired when the server responds
	collector.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
		// fmt.Println("Got a response from", r)
	})

	// // triggered when a CSS selector matches an element
	// collector.OnHTML("a", func(e *colly.HTMLElement) {
	// 	// printing all URLs associated with the <a> tag on the page
	// 	fmt.Println("%v", e.Attr("href"))
	// })

	// // Get products from the website: https://www.scrapingcourse.com/ecommerce
	// collector.OnHTML(".product", func(element *colly.HTMLElement) {
	// 	// initialize a new Product instance
	// 	product := Product{}

	// 	// scrape the target data
	// 	product.Url = element.ChildAttr("a", "href")
	// 	product.Image = element.ChildAttr("img", "src")
	// 	product.Name = element.ChildText(".product-name")
	// 	product.Price = element.ChildText(".price")

	// 	// add the product instance with scraped data to the list of products
	// 	products = append(products, product)
	// })

	// Get products from the website: https://pcx.com.ph/collections/graphics-cards
	collector.OnHTML(".t4s-product-info", func(element *colly.HTMLElement) {
		// initialize a new Product instance
		product := Product{}

		// scrape the target data

		// TODO: add checker if it is a valid url, if not make it into valid url
		url := element.ChildAttr("a", "href")
		if !util.IsUrl(url) {
			url = "https://pcx.com.ph/collections/graphics-cards" + url
			if !util.IsUrl(url) {
				panic("url stll not valid even after adding base path")
			}
		}

		product.Url = url
		product.Image = element.ChildAttr("img", "src")
		product.Name = element.ChildText(".t4s-product-title")
		product.Price = element.ChildText(".t4s-product-price")

		// add the product instance with scraped data to the list of products
		products = append(products, product)
	})

	// triggered once scraping is done (e.g., write data to a CSV file)
	collector.OnScraped(func(r *colly.Response) {
		// fmt.Println(r.Request.URL, " scraped")

		// open the CSV file
		file, err := os.Create("test-products.csv")
		if err != nil {
			log.Fatalln("Failed to create output CSV file", err)
		}
		defer file.Close()

		// initialize a file writer
		writer := csv.NewWriter(file)

		// write the CSV headers
		headers := []string{
			"Url",
			"Image",
			"Name",
			"Price",
		}
		writer.Write(headers)

		// write each product as a CSV row
		for _, product := range products {
			// convert a Product to an array of strings
			record := []string{
				product.Url,
				product.Image,
				product.Name,
				product.Price,
			}

			// add a CSV record to the output file
			writer.Write(record)
		}
		defer writer.Flush()
	})

	// open the target URL
	// collector.Visit("https://www.scrapingcourse.com/ecommerce")
	collector.Visit("https://pcx.com.ph/collections/graphics-cards")
	// collector.Visit(url)

	// check collected data
	// fmt.Println("VIBE CHECK products: ", products)
}

/**

Notes:
Reference Guide: https://www.zenrows.com/blog/web-scraping-golang#using-colly

These functions are executed in the following order:
1. OnRequest(): Called before performing an HTTP request with Visit().
2. OnError(): Called if an error occurred during the HTTP request.
3. OnResponse(): Called after receiving a response from the server.
4. OnHTML(): Called right after OnResponse() if the received content is HTML.
5. OnScraped(): Called after all OnHTML() callback executions are completed.

*/
