package crawler

import (
	"YellowBloomKnapsack/mini-yektanet/adserver/kvstorage"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type Crawler struct {
	kvstorage        kvstorage.KVStorageInterface
	persianStopWords map[string]bool
	numOftopwords    int
	urls             map[string]int
}

// Publisher IDs.
var publisherIDs = map[string]int{
	"varzesh3": 1,
	"digikala": 2,
	"zoomit":   3,
	"sheypoor": 4,
	"filimo":   5,
}

func NewCrawler(kvstorage kvstorage.KVStorageInterface) *Crawler {
	baseUrl := "http://" + os.Getenv("PUBLISHER_WEBSITE_HOSTNAME") + ":" + os.Getenv("PUBLISHER_WEBSITE_PORT")
	return &Crawler{
		kvstorage: kvstorage,
		persianStopWords: map[string]bool{
			"و": true, "در": true, "به": true, "از": true, "که": true, "این": true, "را": true, "اینجا": true,
			"با": true, "برای": true, "است": true, "آن": true, "یک": true, "تا": true, "هم": true, "کنیم": true,
			"می": true, "بر": true, "بود": true, "شد": true, "یا": true, "وی": true, "اما": true, "داریم": true, "اولین": true,
			"اگر": true, "هر": true, "من": true, "ما": true, "شما": true, "او": true, "آنها": true, "دهیم": true, "آخرین": true,
			"ایشان": true, "بودن": true, "باشند": true, "نیز": true, "چون": true, "چه": true, "نیست": true, "های": true,
			"هیچ": true, "همین": true, "چیزی": true, "دارند": true, "کنند": true, "خواهد": true, "آیا": true, "ها": true,
			"کنید": true, "بدانید": true, "خوش": true, "آمدید": true, "خود": true, "زیاد": true, "کم": true, "زیادی": true,
		},
		numOftopwords: 5,
		urls: map[string]int{
			baseUrl + "/varzesh3": 1,
			baseUrl + "/digikala": 2,
			baseUrl + "/zoomit":   3,
			baseUrl + "/sheypoor": 4,
			baseUrl + "/filimo":   5,
		},
	}
}

func (c *Crawler) readHTML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching URL: received status code %d", resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return content, nil
}

// extractText extracts text content from the HTML node tree.
func (crawler *Crawler) extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		buf.WriteString(crawler.extractText(c))
	}
	return buf.String()
}

func (c *Crawler) normalizeText(text string) string {
	text = strings.ReplaceAll(text, "ي", "ی") // Arabic Yeh to Persian Yeh
	text = strings.ReplaceAll(text, "ك", "ک") // Arabic Kaf to Persian Kaf

	// Remove punctuation using a regex.
	reg, err := regexp.Compile("[^\\p{L}\\p{N}\\s]+")
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return ""
	}
	normalizedText := reg.ReplaceAllString(text, " ")

	return normalizedText
}

func (c *Crawler) findTopWords(text string) []string {
	words := strings.Fields(text)
	wordFreq := make(map[string]int)

	for _, word := range words {
		if !c.persianStopWords[word] {
			wordFreq[word]++
		}
	}

	type wordPair struct {
		word  string
		count int
	}
	var pairs []wordPair
	for word, count := range wordFreq {
		pairs = append(pairs, wordPair{word, count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	// Extract the top words.
	topWords := []string{}
	for i, pair := range pairs {
		if i >= c.numOftopwords {
			break
		}
		topWords = append(topWords, pair.word)
	}

	return topWords
}

func (c *Crawler) Crawl() {
	c.crawlOnce()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				c.crawlOnce()
			}
		}
	}()
}

func (c *Crawler) crawlOnce() {
	var wg sync.WaitGroup

	for filePath, publisherID := range c.urls {
		wg.Add(1)

		go func(filePath string, publisherID int) {
			defer wg.Done()

			content, err := c.readHTML(filePath)
			if err != nil {
				log.Printf("%d+%s: Failed to read file: %v", publisherID, filePath, err)
			}

			node, err := html.Parse(bytes.NewReader(content))
			if err != nil {
				log.Printf("%d+%s: Failed to parse HTML: %v", publisherID, filePath, err)
			}

			rawText := c.extractText(node)
			normalizedText := c.normalizeText(rawText)

			topWords := c.findTopWords(normalizedText)
			publisherIDStr := strconv.Itoa(int(publisherID))

			c.kvstorage.Set(publisherIDStr, strings.Join(topWords, ","))
			log.Printf("%d+%s: %s", publisherID, filePath, strings.Join(topWords, ", "))
		}(filePath, publisherID)
	}

	go func() {
		wg.Wait()
	}()
}
