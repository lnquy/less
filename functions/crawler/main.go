package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/apex/go-apex"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"net/http"
	"strings"
	"time"
	"github.com/prometheus/common/log"
)

var (
	awsRegion   = "ap-southeast-1"
	dynamoTable = "less-crawler-dev"
	ghTrending  = "https://github.com/trending"
)

type (
	Crawler struct {
		client *http.Client
		db     *dynamodb.DynamoDB
		repos  []*Repo
	}

	Repo struct {
		Date        string
		Name        string
		Url         string
		Description string
		Language    string
		Stars       string
		Forks       string
		TodayStars  string
	}
)

func (c *Crawler) Handle(raw json.RawMessage, ctx *apex.Context) (interface{}, error) {
	if err := c.CrawlGithubTrending(); err != nil {
		return err.Error(), err
	}
	if _, err := c.PersistData(); err != nil {
		return err.Error(), err
	}
	return "OK", nil
}

func (c *Crawler) CrawlGithubTrending() error {
	doc, err := goquery.NewDocument(ghTrending)
	if err != nil {
		return err
	}

	// Repositories
	doc.Find("div.explore-content ol.repo-list li").Each(func(i int, s *goquery.Selection) {
		repo := &Repo{
			Date: time.Now().Format("2006-01-02"),
		}
		block := s.Find("div.d-inline-block.col-9.mb-1")
		repo.Url, _ = block.Find("a").Attr("href")
		repo.Name = strings.TrimSpace(block.Find("a").Text())

		block = s.Find("div.py-1")
		repo.Description = strings.TrimSpace(block.Find("p").Text())

		block = s.Find("div.f6")
		repo.Language = strings.TrimSpace(block.Find("span.d-inline-block.mr-3 span").Text())

		block.Find("a.muted-link svg").Each(func(j int, b *goquery.Selection) {
			if attr, ok := b.Attr("aria-label"); ok && attr == "star" {
				repo.Stars = strings.TrimSpace(b.Parent().Text())
			} else if ok && attr == "fork" {
				repo.Forks = strings.TrimSpace(b.Parent().Text())
			}
		})

		repo.TodayStars = strings.TrimSuffix(strings.TrimSpace(block.Find("span.d-inline-block.float-sm-right").Text()), " stars today")

		c.repos = append(c.repos, repo)
	})
	return nil
}

func (c *Crawler) PersistData() (*dynamodb.BatchWriteItemOutput, error) {
	items := make([]*dynamodb.WriteRequest, 0)
	counter := 0
	for _, r := range c.repos {
		// Handle null values for DynamoDB
		if r.Description == "" {
			r.Description = "null"
		}
		if r.Language == "" {
			r.Language = "null"
		}

		putReq := &dynamodb.PutRequest{
			Item: map[string]*dynamodb.AttributeValue{
				"date":        {S: aws.String(r.Date)},
				"name":        {S: aws.String(r.Name)},
				"url":         {S: aws.String(r.Url)},
				"description": {S: aws.String(r.Description)},
				"language":    {S: aws.String(r.Language)},
				"stars":       {S: aws.String(r.Stars)},
				"forks":       {S: aws.String(r.Forks)},
				"today_stars": {S: aws.String(r.TodayStars)},
				"sort":        {N: aws.String(fmt.Sprintf("%d", counter))},
			},
		}
		items = append(items, &dynamodb.WriteRequest{
			PutRequest: putReq,
		})
		counter++
	}

	in := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			dynamoTable: items,
		},
	}
	return c.db.BatchWriteItem(in)
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	db := dynamodb.New(sess, &aws.Config{
		Region: aws.String(awsRegion),
		DisableParamValidation: aws.Bool(true),
	})

	apex.Handle(&Crawler{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		db:    db,
		repos: make([]*Repo, 0),
	})
}
