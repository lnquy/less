# less
Simple Go serverless website on Amazon Web Services (AWS).   
Demo: [https://6epko5iya8.execute-api.ap-southeast-1.amazonaws.com/dev](https://6epko5iya8.execute-api.ap-southeast-1.amazonaws.com/dev).  

Website's frontend is written in `VueJS`, which stored on a public `AWS S3` bucket.  
There're two simple `AWS Lambdas` written in Go (deployed via `apex` with NodeJS shim runtime):

- `less_crawler`: Crawls [Github trending](https://github.com/trending) repositories, parses and persists data to `DynamoDB`. This function is triggered by a `CloudWatch` cron job.
- `less_caterer`: Receives HTTP request from client and lookup on `DynamoDB` for trending repositories by day.

`API Gateway` serves two APIs, one to GET the `index.html` page on S3 bucket (forward/proxy request to the URL of `index.html` file on S3). The other allows client to lookup Gihub trending repositories by day (which calling the `less_caterer` lambda).

## Architecture

This repo contains the code for `Lambda` and `S3` parts in the image below.  
![Architecture](https://github.com/lnquy/less/blob/master/images/less-arch.jpg)

## Deploy on AWS

#### Local development

You have to install [Go SDK](https://golang.org/dl/), [glide](https://github.com/Masterminds/glide), [AWS CLI](http://docs.aws.amazon.com/cli/latest/userguide/installing.html), [apex](http://apex.run/) and [configure the AWS credential](http://apex.run/#aws-credentials) to deploy your functions on AWS Lambda.  

- Clone the repository to your local `$GOPATH` and install Go dependencies:

  ```shell
  $ go get github.com/lnquy/less
  $ cd $GOPATH/src/github.com/lnquy/less
  $ glide install
  ```

- The AWS region is `ap-southeast-1`, you have to take a look on the `functions` code and change the `awsRegion` value to deploy on another region.

#### Lambdas

- Create a role for the Lambdas on AWS IAM which have access to `DynamoDB` service.

- Change the role to match your IAM Lambda role: [project.json#L4](https://github.com/lnquy/less/blob/e84fa6f1d12d83e3ff7f0bd92fbc96e044a122f7/project.json#L4).

- Deploy to AWS:

  ```shell
  $ apex deploy  
  ```

#### Frontend and S3

- Create a public S3 bucket for static web hosting, get the URL of the bucket.

- Change the `publicPath` to your S3 bucket URL: [webpack.config.js#L6](https://github.com/lnquy/less/blob/e84fa6f1d12d83e3ff7f0bd92fbc96e044a122f7/frontend/webpack.config.js#L6)

- Change the `GetCaterer` to your `POST /api/v1/trending/` API URL: [main.js#L22](https://github.com/lnquy/less/blob/e84fa6f1d12d83e3ff7f0bd92fbc96e044a122f7/frontend/src/main.js#L22)

- Build frontend:

  ```
  $ cd frontend
  $ npm install   // or yarn
  $ npm run build
  ```

- Upload all `frontend/dist` files to the root of S3 bucket, make sure all files has public read permission.

- Note the HTTP URL to the `index.html` file.

#### CloudWatch

- Create a `24_hours` interval scheduled job.
- Apply that job as the trigger of the `less_crawler` lambda.

#### DynamoDB

- Create a table with name of `less-crawler-dev` or anything you want, just make sure to change the `dynamoTable` value in Go code, too.

- Primary partition key: `date` (String).

- Primary sort key: `sort` (Number).

#### API Gateway

- Create a root `/` GET API to forward/proxy the HTTP request to the URL of `index.html` file on S3 bucket.
- Create a `/api/v1/trending` POST API which called the `less_caterer` lambda.
- You may have to allow the CORS permission on APIs, too.
- Deploy the APIs to a stage (E.g: `dev`).


Open browser and follow the link to the `GET /` API, you now have a simple serverless website up and running on AWS. Congrats :)

![less](https://github.com/lnquy/less/blob/master/images/less-demo.jpg)

## License

This project is under the MIT License. See the [LICENSE](https://github.com/lnquy/less/blob/master/LICENSE) file for the full license text.
