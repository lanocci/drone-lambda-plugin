package main

import (
    "fmt"
    "os"
    "github.com/aws/aws-sdk-go/service/lambda"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/kelseyhightower/envconfig"
)

type Config struct {
	FunctionName string `envconfig:"PLUGIN_FUNCTION_NAME" requeired:"true"`
	S3Bucket string `envconfig:"PLUGIN_S3_BUCKET" required:"true"`
	S3Key string `envconfig:"PLUGIN_FILE_NAME" required:"true"`
	Region string `envconfig:"PLUGIN_LAMBDA_REGION" default:"us-east-1"`
}


func main() {
	var conf Config
	envconfig.Process("", &conf)
    svc := lambda.New(session.New(&aws.Config{
        Region: aws.String(conf.Region),
    }))

    input := &lambda.UpdateFunctionCodeInput{
        FunctionName:    aws.String(conf.FunctionName),
        Publish:         aws.Bool(true),
        S3Bucket:        aws.String(conf.S3Bucket),
        S3Key:           aws.String(conf.S3Key),
    }

    result, err := svc.UpdateFunctionCode(input)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
                case lambda.ErrCodeServiceException:
                    fmt.Println(lambda.ErrCodeServiceException, aerr.Error())
                case lambda.ErrCodeResourceNotFoundException:
                    fmt.Println(lambda.ErrCodeResourceNotFoundException, aerr.Error())
                case lambda.ErrCodeInvalidParameterValueException:
                    fmt.Println(lambda.ErrCodeInvalidParameterValueException, aerr.Error())
                case lambda.ErrCodeTooManyRequestsException:
                    fmt.Println(lambda.ErrCodeTooManyRequestsException, aerr.Error())
                case lambda.ErrCodeCodeStorageExceededException:
                    fmt.Println(lambda.ErrCodeCodeStorageExceededException, aerr.Error())
                default:
                    fmt.Println(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            fmt.Println(err.Error())
        }
        os.Exit(1)
    }

    fmt.Println(result)
}
