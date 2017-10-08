package main

import (
	"encoding/json"
	"github.com/apex/go-apex"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"time"
	"sort"
)

var (
	awsRegion   = "ap-southeast-1"
	dynamoTable = "less-crawler-dev"
)

type (
	Caterer struct {
		db *dynamodb.DynamoDB
	}

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

func (c *Caterer) Handle(raw json.RawMessage, ctx *apex.Context) (interface{}, error) {
	reqData := &struct {
		Date string `json:"date" omitempty`
	}{}
	if err := json.Unmarshal(raw, reqData); err != nil || reqData.Date == "" {
		reqData.Date = time.Now().Format("2006-01-02")
	}

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

	res := make([]*Repo, 0)

	for _, item := range resp.Items {
		repo := &Repo{}
		err := dynamodbattribute.UnmarshalMap(item, repo)
		if err != nil {
			continue
		}
		res = append(res, repo)
	}

	sort.Sort(BySortIndex(res))
	b, err := json.Marshal(res)
	if err != nil {
		return err.Error(), err
	}
	return string(b), nil
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	db := dynamodb.New(sess, &aws.Config{
		Region:                 aws.String(awsRegion),
		DisableParamValidation: aws.Bool(true),
	})
	apex.Handle(&Caterer{
		db: db,
	})
}

type BySortIndex []*Repo

func (r BySortIndex) Len() int           { return len(r) }
func (r BySortIndex) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r BySortIndex) Less(i, j int) bool { return r[i].Sort < r[j].Sort }
