go get github.com/PuerkitoBio/goquery
go get github.com/gocolly/colly/v2
go get jaytaylor.com/html2text
go clean -modcache
go mod download
go mod verify
go mod tidy
