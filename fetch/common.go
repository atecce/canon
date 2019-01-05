package fetch

import (
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

const domain = "https://www.gutenberg.org/"

var authorCollector *colly.Collector

func init() {
	authorCollector = colly.NewCollector()

	authorCollector.OnRequest(func(r *colly.Request) {
		logrus.WithFields(logrus.Fields{
			"path":   r.URL.Path,
			"method": r.Method,
		}).Info("crawling")
	})
}
