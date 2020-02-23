package main

import (
  "archive/zip"
  "bytes"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "net/url"
  "os"
  "blackhatgo/BING/metadata"
  "github.com/PuerkitoBio/goquery"
)

func handler(i int, s *goquery.Selection) {
  url, ok := s.Find("a").Attr("href")
  if !ok {
    fmt.Println("handler error: searching for links")
    return
  }

  fmt.Printf("%d: %s\n", i, url)
  res, err := http.Get(url)
  if err != nil {
    fmt.Println("error http.Get failed")
    return
  }

  buf, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println("error, ioutil.ReadAll failed")
    return
  }
  defer res.Body.Close()

  r, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
  if err != nil {
    fmt.Println("zip.NewReader failed")
    return
  }

  cp, ap, err := metadata.NewProperties(r)
  if err != nil {
    fmt.Println("metadata new props part failed")
    return
  }

  log.Printf(
    "%25s %25s - %s %s\n",
    cp.Creator,
    cp.LastModifiedBy,
    ap.Application,
    ap.GetMajorVersion())
}

func main() {
  if len(os.Args) != 3 {
    log.Fatalln("Missing argument: <domain> <file extension>")
  }
  domain  := os.Args[1]
  filetype := os.Args[2]

  q := fmt.Sprintf(
    "site:%s && filetype:%s && instreamset:(url title):%s",
    domain,
    filetype,
    filetype)
  search := fmt.Sprintf("http://www.bing.com/search?q=%s", url.QueryEscape(q))
  doc, err := goquery.NewDocument(search)
  if err != nil {
    log.Panicln(err)
  }

  s := "html body div#b_content main ol#b_results li.b_algo h2"
  doc.Find(s).Each(handler)
  fmt.Println(search)
}
