package amazonapi

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3session *s3.S3

const (
	//BucketName stores the name of the bucker we want to access
	BucketName = "mygolangwebsite"
	//Region is AWS region
	Region = "ap-south-1"
)

//Init initialises the session
func Init() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(Region),
		Credentials: credentials.NewSharedCredentials("", "baidu99"),
	})))
}

func listBuckets() (resp *s3.ListBucketsOutput) {
	resp, err := s3session.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		panic(err)
	}
	return resp
}

func createBucket() (resp *s3.CreateBucketOutput) {
	resp, err := s3session.CreateBucket(&s3.CreateBucketInput{
		ACL: aws.String(s3.BucketCannedACLPrivate),
		//ACL:    aws.String(s3.BucketCannedACLPublicRead),
		Bucket: aws.String(BucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(Region),
		},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println("Bucket name already in use!")
				panic(err)
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println("Bucket exists and is owned by you!")
			default:
				panic(err)
			}
		}
	}

	return resp
}

//UploadObject uploads a new key in the specified bucket
func UploadObject(filename string) (resp *s3.PutObjectOutput) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	resp, err = s3session.PutObject(&s3.PutObjectInput{
		Body:   f,
		Bucket: aws.String(BucketName),
		Key:    aws.String(strings.Split(filename, "/")[1]),
		ACL:    aws.String(s3.BucketCannedACLPrivate),
	})

	if err != nil {
		panic(err)
	}

	return resp
}

//ListObjects lists details specific key in the specified bucket
func ListObjects(filename string) *s3.Object {
	resp, err := s3session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(BucketName),
	})
	if err != nil {
		panic(err)
	}
	var details *s3.Object
	metadata := resp.Contents
	for _, object := range metadata {
		if *object.Key == filename {
			details = object
		}
	}
	return details
}

//ListObjectsWithPagination lists all the keys in the specified bucket with Pagination
func ListObjectsWithPagination() {
	params := &s3.ListObjectsInput{
		Bucket:  aws.String(BucketName),
		MaxKeys: aws.Int64(10),
	}
	pageNum := 0

	err := s3session.ListObjectsPages(params,

		func(page *s3.ListObjectsOutput, lastPage bool) bool {

			pageNum++
			fmt.Printf("\n\n\n\n\nPAGE %v :=\n\n", pageNum)
			for _, value := range page.Contents {
				fmt.Printf("\tName :-\t\t%v\n", *value.Key)
				fmt.Printf("\tSize :-\t\t%v bytes\n", *value.Size)
				fmt.Printf("\tTimeStamp :-\t%v\n\n", *value.LastModified)
			}
			return pageNum <= 10

		})
	if err != nil {
		panic(err)
	}
}

//GetObject fetches the specified key in the specified bucket into our download folder
func GetObject(filename string, path string) {
	resp, err := s3session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(filename),
	})

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	err = ioutil.WriteFile(path, body, 0644)
	if err != nil {
		panic(err)
	}
}

//DeleteObject deletes the specified key in the specified bucket
func DeleteObject(filename string) (resp *s3.DeleteObjectOutput) {
	resp, err := s3session.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		panic(err)
	}

	return resp
}
