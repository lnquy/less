# less
Simple Go serverless website on Amazone Web Services (AWS).

## Architecture

![Architecture](https://github.com/lnquy/less/blob/master/images/less-arch.jpg)

This repo (`less`) is the AWS Lambda part in the architacture image above.
All other parts must be setup and configure manually via AWS Console.

##  Deploy Go functions on AWS Lambda

You have to install Go SDK, AWS CLI, apex and configure the credential to able to deploy your functions on AWS Lambda.

```shell
$ go get github.com/lnquy/less
$ cd $GOPATH/src/github.com/lnquy/less
$ apex deploy
```

## License

This project is under the MIT License. See the [LICENSE](https://github.com/lnquy/less/blob/master/LICENSE) file for the full license text.