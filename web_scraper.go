package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/time/rate"
)

func main() {
	// Define command-line flags for the program
	baseURL := flag.String("base_url", "", "Base URL of the website")
	startPage := flag.String("start_page", "", "Starting page URL or path")
	selectors := flag.String("selectors", "", "CSS selectors for scraping data")
	maxPages := flag.Int("max_pages", 5, "Maximum number of pages to scrape")
	flag.Parse()

	// Check if required flags are set
	if *baseURL == "" || *startPage == "" || *selectors == "" {
		log.Fatal("base_url, start_page, and selectors are required")
	}

	// Parse CSS selectors from the command-line flag
	cssSelectors := make(map[string]string)
	for _, selector := range strings.Split(*selectors, ",") {
		parts := strings.Split(selector, "=")
		if len(parts) == 2 {
			cssSelectors[parts[0]] = parts[1]
		}
	}

	// Create a temporary directory for caching webpage content
	cacheDir, err := os.MkdirTemp("", "web_cache")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(cacheDir) // Remove the cache directory when done

	// Scrape the website and store the data in the allData variable
	data := scrapeWebsite(*baseURL, *startPage, cssSelectors, *maxPages, cacheDir)
	for _, pageData := range data {
		fmt.Println(pageData) // Print the scraped data
	}

	// Save to CSV
	if err := saveToCSV("books.csv", data[0]); err != nil {
		log.Fatal(err)
	}

}

// Scrape the website and return the scraped data
func scrapeWebsite(baseURL, startPage string, cssSelectors map[string]string, maxPages int, cacheDir string) []map[string][]string {
	url := baseURL + startPage
	allData := []map[string][]string{}

	// Create a rate limiter to limit the number of requests per second
	limiter := rate.NewLimiter(5, 1)

	for i := 0; i < maxPages; i++ {
		limiter.Wait(context.Background())       // Wait for the rate limiter to allow the request
		content, err := fetchPage(url, cacheDir) // Fetch the webpage content
		if err != nil {
			log.Println(err)
			break
		}

		data := parsePage(content, cssSelectors) // Parse the webpage content using CSS selectors
		allData = append(allData, data)

		// Find the next page link
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			log.Println(err)
			break
		}
		nextPage := doc.Find("a.next").AttrOr("href", "")
		if nextPage == "" {
			break
		}
		url = nextPage

		time.Sleep(1 * time.Second) // Wait 1 second before making the next request
	}
	return allData
}

// Fetch the webpage content from the cache or from the internet
func fetchPage(url, cacheDir string) (string, error) {
	cacheFile := filepath.Join(cacheDir, strings.ReplaceAll(url, "/", "_"))
	if _, err := os.Stat(cacheFile); err == nil {
		return readFile(cacheFile) // Return the cached content
	}

	resp, err := http.Get(url) // Fetch the webpage content from the internet
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = writeFile(cacheFile, bodyBytes) // Cache the webpage content
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

// Read a file from disk
func readFile(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Write a file to disk
func writeFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}

// Parse the webpage content using CSS selectors
func parsePage(content string, cssSelectors map[string]string) map[string][]string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string][]string)
	for key, selector := range cssSelectors {
		elements := doc.Find(selector)
		elements.Each(func(i int, s *goquery.Selection) {
			data[key] = append(data[key], s.Text())
		})
	}
	return data
}

// Function to save data to a CSV file
func saveToCSV(filename string, data map[string][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{"Title", "Price"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write data rows
	for i := 0; i < len(data["title"]); i++ {
		row := []string{data["title"][i], data["price"][i]}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
