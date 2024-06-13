package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: <executable> <branch name>")
		os.Exit(1)
	}

	// Fetch credentials and configuration from environment variables
	accessKey := os.Getenv("R2_ACCESS_KEY")
	secretKey := os.Getenv("R2_SECRET_KEY")
	s3URL := os.Getenv("R2_ENDPOINT")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	branchName := os.Args[1]

	// Ensure all required environment variables are set
	if accessKey == "" || secretKey == "" || s3URL == "" || bucketName == "" {
		fmt.Println("Please set R2_ACCESS_KEY, R2_SECRET_KEY, R2_ENDPOINT, and R2_BUCKET_NAME environment variables")
		os.Exit(1)
	}

	// Load AWS configuration with credentials
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
		config.WithRegion("auto"), // Ensure this region is appropriate or set it via environment variable if needed
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL: s3URL,
					}, nil
				},
			),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration, %v", err)
		os.Exit(1)
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	// Replace '/' with '_' in the branch name
	safeBranchName := strings.ReplaceAll(branchName, "/", "_")

	// Construct the key for the snapshot file
	key := fmt.Sprintf("elys-snapshot-%s.tar.lz4", safeBranchName)

	// Delete the file
	_, err = client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to delete file %q from bucket %q, %v", key, bucketName, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully deleted %q from %q\n", key, bucketName)
}
