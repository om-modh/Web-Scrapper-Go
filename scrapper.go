package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// initializing a data structure to keep the scraped data
type PokemonProduct struct {
	url, image, name, price string
}

func main() {
	// initializing the slice of structs to store the data to scrape
	var pokemonProducts []PokemonProduct

	// creating a new Colly instance
	c := colly.NewCollector(colly.AllowedDomains("www.scrapeme.live", "scrapeme.live"))
	// flag := true
	i := 1
	limit := 5

	
	pagetoscrape := "https://scrapeme.live/shop/page/1/"
	var pagestoscrape []string
	discoveredPages := []string{pagetoscrape}

	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		newPage := e.Attr("href")

		if !contains(pagestoscrape, newPage) {
			if !contains(discoveredPages, newPage) {
				pagestoscrape = append(pagestoscrape, newPage)
			}
			discoveredPages = append(discoveredPages, newPage)
		}
	})

	// scraping logic
	c.OnHTML("li.type-product", func(e *colly.HTMLElement) {
		pokemonProduct := PokemonProduct{}
		pokemonProduct.url = e.ChildAttr("a", "href")
		pokemonProduct.image = e.ChildAttr("img", "src")
		pokemonProduct.name = e.ChildText("h2")
		pokemonProduct.price = e.ChildText(".price")

		pokemonProducts = append(pokemonProducts, pokemonProduct)
	})

	c.OnScraped(func(response *colly.Response) {
		if len(pagestoscrape) != 0 && i < limit {
			pagetoscrape = pagestoscrape[0]
			pagestoscrape = pagestoscrape[1:]
			i++
			c.Visit(pagetoscrape)
		}
	})
	// visiting the target page
	c.Visit(pagetoscrape)

	//-------------------------------------------------------------------------------------------------------------

	// opening the CSV file
	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	// initializing a file writer
	writer := csv.NewWriter(file)

	// writing the CSV headers
	headers := []string{
		"url",
		"image",
		"name",
		"price",
	}
	writer.Write(headers)

	// writing each Pokemon product as a CSV row
	for _, pokemonProduct := range pokemonProducts {
		// converting a PokemonProduct to an array of strings
		record := []string{
			pokemonProduct.url,
			pokemonProduct.image,
			pokemonProduct.name,
			pokemonProduct.price,
		}

		// adding a CSV record to the output file
		writer.Write(record)
	}

	defer writer.Flush()
}
