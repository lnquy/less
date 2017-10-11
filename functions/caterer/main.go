package main

import (
	"encoding/json"
	"github.com/apex/go-apex"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"sort"
	"time"
)

var (
	awsRegion   = "ap-southeast-1"   // AWS region to deploy Lambdas and other services
	dynamoTable = "less-crawler-dev" // DynamoDB table name to read from
)

type (
	// Caterer serves the POST /trending API which lookup from DynamoDB for trending repositories by day.
	Caterer struct {
		db *dynamodb.DynamoDB // DynamoDB client
	}

	// Repo represents a Github trending repository.
	Repo struct {
		Date        string `json:"date"`
		Name        string `json:"name"`
		Url         string `json:"url"`
		Description string `json:"description"`
		Language    string `json:"language"`
		Stars       string `json:"stars"`
		Forks       string `json:"forks"`
		TodayStars  string `json:"today_stars"`
		Sort        int    `json:"sort"`
	}
)

// Handle implemented from apex.Handler which will be called when the lambda started up.
// Parse request to get the day value, it not specify, current server time will be used.
// Lookup from database for trending repository by day.
func (c *Caterer) Handle(raw json.RawMessage, ctx *apex.Context) (interface{}, error) {
	// Parse request body
	reqData := &struct {
		Date string `json:"date" omitempty`
	}{}
	if err := json.Unmarshal(raw, reqData); err != nil || reqData.Date == "" {
		reqData.Date = time.Now().Format("2006-01-02")
	}

	// Build lookup query then query
	query := &dynamodb.QueryInput{
		TableName: aws.String(dynamoTable),
		KeyConditions: map[string]*dynamodb.Condition{
			"date": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{S: aws.String(reqData.Date)},
				},
			},
		},
	}
	resp, err := c.db.Query(query)
	if err != nil {
		return err.Error(), err
	}

	// Map data returned from DynamoDB to actual Repo instances
	res := make([]*Repo, 0)
	for _, item := range resp.Items {
		repo := &Repo{}
		err := dynamodbattribute.UnmarshalMap(item, repo)
		if err != nil {
			continue
		}
		res = append(res, repo)
	}

	sort.Sort(BySortIndex(res)) // Re-sort by sort index (optional)
	b, err := json.Marshal(res)
	if err != nil {
		return err.Error(), err
	}
	return string(b), nil
}

func main() {
	// New DynamoDB client
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	db := dynamodb.New(sess, &aws.Config{
		Region:                 aws.String(awsRegion),
		DisableParamValidation: aws.Bool(true),
	})

	// Lambda handler
	apex.Handle(&Caterer{
		db: db,
	})
}

// BySortIndex implements sort interface which allows to sort a list of Repos by sort index.
type BySortIndex []*Repo

func (r BySortIndex) Len() int           { return len(r) }
func (r BySortIndex) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r BySortIndex) Less(i, j int) bool { return r[i].Sort < r[j].Sort }
