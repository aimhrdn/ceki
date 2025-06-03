package main

import (
	"fmt"
    "encoding/json"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/sts"
    "github.com/spf13/cobra"

)

type ValidationResult struct {
    Valid bool   `json:"valid"`
    ARN   string `json:"arn,omitempty"`
    Error string `json:"error,omitempty"`
}

type S3AccessResult struct {
    Success bool     `json:"success"`
    Buckets []string `json:"buckets,omitempty"`
    Error   string   `json:"error,omitempty"`
}

func validateAWSKeys(accessKey, secretKey, region string) ValidationResult {
    sess, err := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
        Region:      aws.String(region),
    })
    if err != nil {
        return ValidationResult{Valid: false, Error: err.Error()}
    }

    svc := sts.New(sess)
    input := &sts.GetCallerIdentityInput{}
    result, err := svc.GetCallerIdentity(input)
    if err != nil {
        return ValidationResult{Valid: false, Error: err.Error()}
    }

    return ValidationResult{Valid: true, ARN: *result.Arn}
}

func checkS3Access(accessKey, secretKey, region string) S3AccessResult {
    sess, err := session.NewSession(&aws.Config{
        Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
        Region:      aws.String(region),
    })
    if err != nil {
        return S3AccessResult{Success: false, Error: err.Error()}
    }

    svc := s3.New(sess)
    result, err := svc.ListBuckets(nil)
    if err != nil {
        return S3AccessResult{Success: false, Error: err.Error()}
    }

    var bucketNames []string
    for _, bucket := range result.Buckets {
        bucketNames = append(bucketNames, *bucket.Name)
    }

    return S3AccessResult{Success: true, Buckets: bucketNames}
}

func main() {
    var accessKey string
    var secretKey string
    var region string

    var rootCmd = &cobra.Command{
        Use:   "ceki",
        Short: "A tool to validate and check AWS keys",
    }

    var validateCmd = &cobra.Command{
        Use:   "validate",
        Short: "Validate AWS access keys",
        Run: func(cmd *cobra.Command, args []string) {
            result := validateAWSKeys(accessKey, secretKey, region)
            jsonResult, _ := json.Marshal(result)
            fmt.Println(string(jsonResult))
        },
    }

    validateCmd.Flags().StringVarP(&accessKey, "key", "k", "", "AWS Access Key ID")
    validateCmd.Flags().StringVarP(&secretKey, "secret", "s", "", "AWS Secret Access Key")
    validateCmd.Flags().StringVarP(&region, "region", "r", "us-east-1", "AWS Region")
    validateCmd.MarkFlagRequired("key")
    validateCmd.MarkFlagRequired("secret")

    var checkCmd = &cobra.Command{
        Use:   "check",
        Short: "Check AWS S3 access",
        Run: func(cmd *cobra.Command, args []string) {
            result := checkS3Access(accessKey, secretKey, region)
            jsonResult, _ := json.Marshal(result)
            fmt.Println(string(jsonResult))
        },
    }

    checkCmd.Flags().StringVarP(&accessKey, "key", "k", "", "AWS Access Key ID")
    checkCmd.Flags().StringVarP(&secretKey, "secret", "s", "", "AWS Secret Access Key")
    checkCmd.Flags().StringVarP(&region, "region", "r", "us-east-1", "AWS Region")
    checkCmd.MarkFlagRequired("key")
    checkCmd.MarkFlagRequired("secret")

    rootCmd.AddCommand(validateCmd)
    rootCmd.AddCommand(checkCmd)
    rootCmd.Execute()
}
