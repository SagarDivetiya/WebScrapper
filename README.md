**Web Scraper**
================

A robust and flexible web scraper written in Go that extracts data from websites using CSS selectors.

**Usage**
-----

To run the web scraper, use the following command: 

```bash
go run web_scraper.go -base_url=https://example.com -start_page=/start-page -selectors=title=h1,content=.article-content -max_pages=10
```

**Flags**
------

* `-base_url`: The base URL of the website to scrape.
* `-start_page`: The starting page URL or path on the website.
* `-selectors`: A comma-separated list of CSS selectors to extract data from the webpage. For example, `title=h1,content=.article-content`.
* `-max_pages`: The maximum number of pages to scrape from the website.

How It Works
============
### Fetching the Webpage

The scraper sends an HTTP request to the URL specified by the `-base_url` flag.
It then parses the HTML response using the **goquery** library.

### Extracting Data

The scraper uses the `-selectors` flag to determine which data to extract.
It finds matching elements on the page and extracts their **text**.

### Storing Data

Extracted data is stored in a **cache directory** to avoid re-fetching the same page.

### Following Pagination Links

The scraper identifies and follows **pagination links** to scrape additional pages.

### Printing Data

Extracted data is printed to the **console**.

**Example Sites for Testing**
---------------------------

* Example Domain: `http://example.com`
	+ Start Page: `/`
	+ Selectors: `title=h1`, `content=p`
* Books to Scrape: `http://books.toscrape.com`
	+ Start Page: `/`
	+ Selectors: `title=h3`, `price=.price_color`
* Quotes to Scrape: `http://quotes.toscrape.com`
	+ Start Page: `/`
	+ Selectors: `quote=.text`, `author=.author`

**Possible Improvements**
-------------------------

* Formatting: Output data in CSV or JSON format for further analysis.
* Pagination Handling: Ensure pagination logic correctly follows links and handles multiple pages.
* Error Handling: Implement better error handling and logging to identify and debug issues.
* Rate Limiting: Adjust rate limiting to respect website scraping policies.

**CSV Saving Function**
---------------------

The `saveToCSV` function creates a new CSV file, writes headers, and populates rows with the extracted data.

**License**
-------

This project is licensed under the MIT License. See `LICENSE` for details.

**Contributing**
------------

Contributions are welcome! If you have any issues or suggestions, please open an issue or submit a pull request.
