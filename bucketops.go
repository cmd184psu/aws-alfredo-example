package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cmd184psu/alfredo"
)

type S3Details struct {
	Bucket      string `json:"bucket"`
	Region      string `json:"region"`
	Credentials struct {
		AccessKeyID     string `json:"accessKeyId"`
		SecretAccessKey string `json:"secretAccessKey"`
	} `json:"credentials"`
	Endpoint string `json:"endpoint"`
	Profile  string `json:"profile"`
}

func (details *S3Details) GetSession() alfredo.S3ClientSession {
	if len(details.Profile) == 0 {
		details.Profile = "default"
	}

	var s3c alfredo.S3ClientSession
	s3c = s3c.WithEndpoint(details.Endpoint).WithBucket(details.Bucket).WithRegion(details.Region)

	s3c.Credentials.Profile = details.Profile

	if len(details.Credentials.AccessKeyID) > 0 {
		s3c.Credentials.AccessKey = details.Credentials.AccessKeyID
		s3c.Credentials.SecretKey = details.Credentials.SecretAccessKey
	}

	if len(s3c.Credentials.AccessKey) == 0 {
		fmt.Printf("Loading credentials from file %s using profile %s\n", credentialsFile, details.Profile)
		if err := s3c.LoadCredentials(credentialsFile); err != nil {
			panic("!!! Can't load credentials!!!")
		}
	}
	if err := s3c.EstablishSession(); err != nil {
		panic("unable to establish session")
	}
	return s3c
}

func (details *S3Details) HeadBucket() error {
	s3c := details.GetSession()
	b, err := s3c.HeadBucket()
	if err != nil {
		return err
	}

	if b {
		fmt.Println("Bucket exists")
	} else {
		fmt.Println("Bucket does not exist")
	}
	fmt.Println("Now directly from s3 client provided by aws sdk")

	output, err := s3c.Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(details.Bucket),
	})

	if err != nil {
		return err
	}
	fmt.Println("output from head bucket: " + output.String())
	fmt.Println("raw output:")
	fmt.Println(alfredo.PrettyPrint(output))
	return nil
}

func (details *S3Details) ListBuckets() error {
	if len(details.Profile) == 0 {
		details.Profile = "default"
	}
	s3c := details.GetSession()
	s := s3c.ListBuckets()
	fmt.Println("all buckets available to this profile:")
	fmt.Println(alfredo.PrettyPrint(s))

	fmt.Println("Now directly from s3 client provided by aws sdk")

	output, err := s3c.Client.ListBuckets(&s3.ListBucketsInput{})

	if err != nil {
		return err
	}
	fmt.Println("output from list buckets: " + output.String())
	fmt.Println("raw output:")
	fmt.Println(alfredo.PrettyPrint(output))
	return nil
}
