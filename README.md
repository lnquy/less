# less
Simple Go serverless website on Amazone Web Services (AWS).

## Architecture

![Architecture](https://github.com/lnquy/less/blob/master/images/less-arch.jpg)

This repo (`less`) is the AWS Lambda part in the architecture image above. The role to deploy lambda functions must have access to `DynamoDB` service.  
All other parts must be setup and configure manually via AWS Console (`ap-southeast-1` region):

- S3:
  - Create a S3 static web hosting.
  - Upload all `frontend` files to the root of S3.
  - Take the HTTP URL to the index file.
- API gateway:
  - Create a root GET API `/` proxying to the HTTP URL of index file on S3.
  - Create a POST API `/api/v1/trending` which called the `less_caterer` lambda.
- Cloudwatch:
  - Create a `24_hours` scheduled job.
  - Apply that job as the trigger of the `less_crawler` lambda.
- DynamoDB:
  - Create a table with name of `less-crawler-dev`.
  - Primary partition key: `date` (String).
  - Primary sort key: `sort` (Number).

##  Deploy Go functions on AWS Lambda

You have to install [Go SDK](https://golang.org/dl/), [AWS CLI](http://docs.aws.amazon.com/cli/latest/userguide/installing.html), [apex](http://apex.run/) and [configure the AWS credential](http://apex.run/#aws-credentials) to able to deploy your functions on AWS Lambda.

```shell
$ go get github.com/lnquy/less
$ cd $GOPATH/src/github.com/lnquy/less
$ apex deploy
```

## License

This project is under the MIT License. See the [LICENSE](https://github.com/lnquy/less/blob/master/LICENSE) file for the full license text.