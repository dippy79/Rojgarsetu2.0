package main

import (
    "database/sql"
    "net/http"
    "os"
    "time"
    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
    "golang.org/x/net/html"
)

type Crawler struct {
    DB *sql.DB
}

func main() {
    db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    c := &Crawler{DB: db}

    ticker := time.NewTicker(6 * time.Hour)
    go func() {
        for range ticker.C {
            c.RunCrawl()
        }
    }()

    r := gin.Default()
    r.GET("/health", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"status": "healthy"}) })
    r.POST("/crawl", func(ctx *gin.Context) {
        go c.RunCrawl()
        ctx.JSON(200, gin.H{"message": "Crawl triggered"})
    })
    r.Run(":8082")
}

func (c *Crawler) RunCrawl() {
    sources := []string{
        "https://www.upsc.gov.in/examinations/active-exams",
        "https://ssc.gov.in/",
        "https://www.ncs.gov.in/",
    }

    for _, src := range sources {
        c.fetchAndProcess(src)
    }
}

func (c *Crawler) fetchAndProcess(url string) {
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    doc, _ := html.Parse(resp.Body)
    var links []string
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "a" {
            for _, a := range n.Attr {
                if a.Key == "href" {
                    links = append(links, a.Val)
                }
            }
        }
        for child := n.FirstChild; child != nil; child = child.NextSibling {
            f(child)
        }
    }
    f(doc)

    for _, link := range links {
        if c.DB == nil {
            continue
        }
        var id int
        err := c.DB.QueryRow("SELECT id FROM jobs_government WHERE apply_url=$1", link).Scan(&id)
        if err == sql.ErrNoRows {
            c.DB.Exec("INSERT INTO jobs_government (title, apply_url, source, created_at) VALUES ($1, $2, $3, NOW())",
                "Govt Job Notification", link, url)
        }
    }
}
