package fileupload

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type fileUpload struct{}

type UploadFileInterface interface {
	UploadFile(file *multipart.FileHeader) (string, map[string]string)
}

//So what is exposed is Uploader
var FileUpload UploadFileInterface = &fileUpload{}

func (fu *fileUpload) UploadFile(file *multipart.FileHeader) (string, map[string]string) {

	errList := map[string]string{}

	f, err := file.Open()
	if err != nil {
		errList["Not_Image"] = "Please Upload a valid image"
		return "", errList
	}
	defer f.Close()

	size := file.Size
	//The image should not be more than 500KB
	fmt.Println("the size: ", size)
	if size > int64(5120000) {
		errList["Too_large"] = "Sorry, Please upload an Image of 5MB or less"
		return "", errList

	}
	//only the first 512 bytes are used to sniff the content type of a file,
	//so, so no need to read the entire bytes of a file.
	buffer := make([]byte, size)
	f.Read(buffer)
	fileType := http.DetectContentType(buffer)
	//if the image is valid
	if !strings.HasPrefix(fileType, "image") {
		errList["Not_Image"] = "Please Upload a valid image"
		return "", errList
	}
	filePath := FormatFile(file.Filename)
	fileBytes := bytes.NewReader(buffer)

	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_KEY")
	secKey := os.Getenv("AWS_SECRET")
	// endpoint := "s3.amazonaws.com"
	// ssl := true

	s, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secKey,
			""),
	})

	uploader := s3manager.NewUploader(s)

	resp, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("go-auction"),
		Key:    aws.String(filePath),
		Body:   fileBytes,
	})
	if err != nil {
		fmt.Println("the error", err)
		errList["Other_Err"] = "something went wrong"
		return "", errList
	}
	return resp.Location, nil
	// Initiate a client using DigitalOcean Spaces.
	// client, err := minio.New(endpoint, accessKey, secKey, ssl)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fileBytes := bytes.NewReader(buffer)
	// cacheControl := "max-age=31536000"
	// // make it public
	// userMetaData := map[string]string{"x-amz-acl": "public-read"}
	// n, err := client.PutObject("go-auction", filePath, fileBytes, size, minio.PutObjectOptions{ContentType: fileType, CacheControl: cacheControl, UserMetadata: userMetaData})
	// if err != nil {
	// 	fmt.Println("the error", err)
	// 	errList["Other_Err"] = "something went wrong"
	// 	return "", errList
	// }
	// fmt.Println("Successfully uploaded bytes: ", n)
	// return filePath, nil
}
