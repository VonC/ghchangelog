package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/VonC/ghchangelog/version"

	"github.com/PuerkitoBio/goquery"
	"github.com/atotto/clipboard"
	"github.com/gocolly/colly/v2"

	"gopkg.in/toast.v1"
)

type article struct {
	e       *colly.HTMLElement
	title   string
	href    string
	date    string
	fromURL bool
}

type articles []*article

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

	ghurl := "https://github.blog/changelog/"
	iquery := 1

	if len(os.Args) >= 2 && strings.HasPrefix(os.Args[1], "https") {
		ghurl = os.Args[1]
		iquery = 2
	}

	query := ""
	if len(os.Args) > 2 || iquery == 1 {
		query = strings.Join(os.Args[iquery:], " ")
		query = strings.ToLower(query)
	}
	u, err := url.Parse(ghurl)
	if err != nil {
		log.Fatal(err)
	}
	ghdomain := u.Hostname()

	c := colly.NewCollector(
		colly.AllowedDomains(ghdomain),
		colly.MaxDepth(0),
	)
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL: ", r.Request.URL, " failed with response: ", r.StatusCode, "\nError: ", err)
	})

	articles := make(articles, 0)
	goqueryselector := "article"
	if iquery == 2 {
		goqueryselector = "html"
	}
	// https://htmlcheatsheet.com/jquery/
	c.OnHTML(goqueryselector, func(e *colly.HTMLElement) { //class that contains wanted info
		article := &article{e: e}
		if iquery == 2 {
			article.href = ghurl
			article.fromURL = true
		}
		title := article.getTitle()
		//fmt.Printf("Check title '%s' with query '%s'\n", title, query)
		if strings.Contains(strings.ToLower(title), query) {
			articles = append(articles, article)
		}
	})

	if err := c.Visit(ghurl); err != nil {
		log.Fatal(err)
	}
	// Wait until threads are finished
	c.Wait()

	if len(articles) == 0 {
		fmt.Printf("\nNo article with title query '%s' found in '%s'", query, ghurl)
		os.Exit(0)
	}
	if len(articles) > 1 {
		fmt.Printf("\nWARNING:\n  Multiples articles with title query '%s' found in '%s':\n\n", query, ghurl)
		for _, article := range articles {
			fmt.Println("- " + article.getTitle())
		}
		os.Exit(0)
	}
	text := articles[0].markdown()
	fmt.Println(text)
	if err = clipboard.WriteAll(text); err != nil {
		log.Fatal(err)
	}
	if runtime.GOOS == "windows" {
		notification := toast.Notification{
			AppID:   "ghchangelog",
			Title:   "Copied",
			Message: "GitHub Changelog entry copied to the clipboard",
		}
		err = notification.Push()
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func (a *article) getTitle() string {
	if a.title != "" {
		return a.title
	}
	if a.fromURL {
		a.title = a.e.ChildText("h1")
	} else {
		a.title = a.e.ChildText("h2.h5-mktg")

	}
	if a.title == "" {
		a.title = "<no title detected>"
	}
	return a.title
}
func (a *article) getHref() string {
	if a.href != "" {
		return a.href
	}
	a.href = a.e.ChildAttr("h2 > a", "href")
	return a.href
}

// reDate: https://regex101.com/r/hQ3tKo/1
var reDate = regexp.MustCompile(`(?m)\d{4}-\d{2}-\d{2}`)

func (a *article) getDate() string {
	if a.date != "" {
		return a.date
	}
	href := a.getHref()
	datetime := reDate.FindString(href)
	if datetime == "" {
		a.date = fmt.Sprintf("<No date detected in href '%s'>", href)
	} else {
		t, err := time.Parse("2006-01-02", datetime)
		if err != nil {
			a.date = fmt.Sprintf("<No date detected in href '%s'>, extracted datetime '%s'", href, datetime)
		}
		a.date = t.Format("Jan. 2006")
	}
	return a.date
}

func (a *article) markdown() string {
	e := a.e
	m := fmt.Sprintf("> ## [%s](%s) (%s)\n", a.getTitle(), a.getHref(), a.getDate())

	body := e.DOM.Find(".post__content").First()
	//fmt.Printf("Body '%+v'\n", body)
	m = m + visitNodes(body)
	return m
}

var ignored = map[string]bool{
	"br": true,
}

func hasParentNamed(parentName string, names ...string) bool {
	//fmt.Printf("Parent name '%s'\n", selname)
	return hasSelNamed(parentName, names...)
}

func hasPrevNamed(prevName string, names ...string) bool {
	//fmt.Printf("Prev name '%s'\n", selname)
	return hasSelNamed(prevName, names...)
}

func hasSelNamed(selname string, names ...string) bool {
	for _, name := range names {
		if selname == name {
			return true
		}
	}
	return false
}

func visitNodes(sel *goquery.Selection) string {
	m := ""
	parentName := goquery.NodeName(sel)
	prevName := ""
	sel.Contents().Each(func(i int, sel *goquery.Selection) {
		nodeName := goquery.NodeName(sel)
		switch nodeName {
		case "#text":
			txt := sel.Text()
			if strings.TrimSpace(txt) != "" {
				r := strings.NewReplacer(". ", ".  \n> ")
				txt = r.Replace(txt)
				if !hasParentNamed(parentName, "li", "a") && !hasPrevNamed(prevName, "code", "a", "strong") {
					m = m + "> "
				}
				m = m + txt
			}
		case "br":
			m = m + "  \n"
		case "p":
			nodes := visitNodes(sel)
			m = m + ">\n" + nodes + "\n"
		case "img":
			alt := sel.AttrOr("alt", "")
			src := sel.AttrOr("src", "")
			//fmt.Printf("PREV '%s'\n", prevName)
			if !hasPrevNamed(prevName, "br", "") {
				m = m + ">\n"
			}
			m = m + "> " + src
			if alt != "" {
				m = m + " -- " + alt
			}
		case "ul":
			m = m + ">\n" + visitNodes(sel)
		case "li":
			m = m + "> - " + visitNodes(sel) + "\n"
		case "a":
			txt := visitNodes(sel)
			href := sel.AttrOr("href", "")
			if !hasPrevNamed(prevName, "#text") {
				m = m + "> "
			}
			m = m + fmt.Sprintf("[%s](%s)", txt, href)
		case "pre":
			m = m + ">\n> " + visitNodes(sel) + "\n"
		case "code":
			txt := sel.Text()
			m = m + fmt.Sprintf("`%s`", txt)
		case "strong":
			txt := sel.Text()
			m = m + fmt.Sprintf("**%s**", txt)
		default:
			if !ignored[nodeName] {
				fmt.Printf("Unknown node '%s'\n", nodeName)
			}
		}
		//fmt.Printf("m for name '%s': '%s'\n", nodeName, m)
		prevName = goquery.NodeName(sel)
	})
	return m
}
