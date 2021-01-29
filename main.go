package main

import (
	"GolangS3API/amazonapi"
	"fmt"
	"os"
)

const uploadFolder = "upload/"
const downloadFolder = "download/"

//File stores files details
type File struct {
	Name      string `json:"FileName"`
	Size      int64  `json:"Size"`
	Timestamp string `json:"Timestamp"`
}

func uploadFetchDetails() {
	var fname string
	fmt.Print("\n\t\tFilename of document to upload? ")
	fmt.Scanf("%v\n", &fname)
	path := (uploadFolder + fname)
	fmt.Printf("\n\t\tUploading %v to S3 from %v...\n", fname, path)
	//Upload to S3
	amazonapi.UploadObject(path)
	//Show details
	details := amazonapi.ListObjects(fname)
	fmt.Printf("\t\t\tName :-\t\t%v\n", *details.Key)
	fmt.Printf("\t\t\tSize :-\t\t%v bytes\n", *details.Size)
	fmt.Printf("\t\t\tTimeStamp :-\t%v", *details.LastModified)
}
func download() {
	var fname string
	fmt.Print("\n\t\tFilename of document to download? ")
	fmt.Scanf("%v\n", &fname)
	path := (downloadFolder + fname)
	//Download from S3
	fmt.Printf("\n\t\tDownloading %v from S3 to %v...", fname, path)
	amazonapi.GetObject(fname, path)
}
func updateFetchDetails() {
	var fname string
	fmt.Print("\n\t\tFilename of document to update? ")
	fmt.Scanf("%v\n", &fname)
	path := (uploadFolder + fname)
	fmt.Printf("\n\t\tUpdating %v to S3 from %v...\n", fname, path)
	//Fetch details before update
	details := amazonapi.ListObjects(fname)
	before := File{*details.Key, *details.Size, (*details.LastModified).String()}
	//Update to S3
	amazonapi.UploadObject(path)
	//Show details
	after := amazonapi.ListObjects(fname)
	fmt.Printf("\t\t\tName :-\t\t%v\n", *after.Key)
	fmt.Printf("\t\t\tSize :-\t\t%v bytes ---> %v bytes\n", before.Size, *after.Size)
	fmt.Printf("\t\t\tTimeStamp :-\t%v ---> %v", before.Timestamp, *after.LastModified)
}
func delete() {
	var fname string
	fmt.Print("\n\t\tFilename of document to delete? ")
	fmt.Scanf("%v\n", &fname)
	fmt.Printf("\n\t\tDeleting %v from S3...", fname)
	//Delete from S3
	amazonapi.DeleteObject(fname)
}
func main() {
	amazonapi.Init()
	var option int
	for {
		fmt.Print("\n\n\nOPTIONS :-\n\n")
		fmt.Print("\t1. Upload a document and fetch details.\n")
		fmt.Print("\t2. Retrive uploaded document from s3.\n")
		fmt.Print("\t3. Retrive all data from s3(with pagination).\n")
		fmt.Print("\t4. Update already uploaded document.\n")
		fmt.Print("\t5. Delete a document.\n")
		fmt.Print("\t6. EXIT\n")
		fmt.Print("\n\tEnter your choice?")

		fmt.Scanf("%v\n", &option)
		switch option {
		case 1:
			uploadFetchDetails()
		case 2:
			download()
		case 3:
			//List all the data with pagination implemented
			amazonapi.ListObjectsWithPagination()
		case 4:
			updateFetchDetails()
		case 5:
			delete()
		case 6:
			os.Exit(0)
		}
	}
}
