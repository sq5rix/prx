package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"jaytaylor.com/html2text"
)

func main() {
	// Read URL from stdin
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Fprintf(os.Stderr, "Error: no URL provided in stdin\n")
		os.Exit(1)
	}
	url := strings.TrimSpace(scanner.Text())

	// Initialize the collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// Create a channel to receive the text content
	textChan := make(chan string, 1)

	// Handle the HTML content
	c.OnHTML("html", func(e *colly.HTMLElement) {
		// Get the document
		doc := e.DOM
 b
		// Remove unwanted elements
		doc.Find("script, style, noscript, iframe, img").Remove()

		// Get the HTML content
		html, err := doc.Html()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting HTML: %v\n", err)
			return
		}

		// Convert HTML to plain text
		text, err := html2text.FromString(html, html2text.Options{
			PrettyTables: false,
			OmitLinks:    true,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting to text: %v\n", err)
			return
		}

		// Clean up the text thoroughly
		text = strings.TrimSpace(text)
		
		// Replace multiple spaces with single space
		space := regexp.MustCompile(`\s+`)
		text = space.ReplaceAllString(text, " ")
		
		// Clean empty lines and normalize line endings
		lines := strings.Split(text, "\n")
		var cleanLines []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				cleanLines = append(cleanLines, line)
			}
		}
		text = strings.Join(cleanLines, "\n")
		
		textChan <- text
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		fmt.Fprintf(os.Stderr, "Error scraping %s: %v\n", url, err)
		close(textChan)
	})

	// Start the scraping
	go func() {
		err := c.Visit(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error visiting URL: %v\n", err)
		}
		close(textChan)
	}()

	// Print the result to stdout
	text := <-textChan
	fmt.Println(text)
}
