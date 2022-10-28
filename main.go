package main

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	maxPartSize        = int64(5 * 1024 * 1024)
	maxRetries         = 3
	awsAccessKeyID     = os.Getenv("S3_ACCESS_ID")
	awsSecretAccessKey = os.Getenv("S3_SECRET_KEY")
	awsBucketRegion    = os.Getenv("S3_REGION")
	awsBucketName      = os.Getenv("S3_BUCKET_NAME")
	awsEndpoint        = os.Getenv("S3_ENDPOINT")
	awsUrlPrefix		   = os.Getenv("S3_URL_PREFIX")
	httpUserName       = os.Getenv("HTTP_USERNAME")
	httpPassWord       = os.Getenv("HTTP_PASSWORD")
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		httpUserName: httpPassWord, // user:foo password:bar
	}))

	authorized.POST("upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		dir, _ := c.GetQuery("dir")
		if dir == "" {
			t := time.Now()
			dir = fmt.Sprintf("%d/%d/%d", t.Year(), t.Month(), t.Day())
		}
		p := path.Join(dir, file.Filename)
		url, err := upload(*file, p)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "err": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "url": url, "path": p})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8081")
}

func upload(file multipart.FileHeader, path string) (string, error) {
	ctx := context.Background()
	endpoint := awsEndpoint
	accessKeyID := awsAccessKeyID
	secretAccessKey := awsSecretAccessKey
	useSSL := true

	fileSize := file.Size
	f, err := file.Open()
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called mymusic.
	bucketName := awsBucketName
	location := awsBucketRegion

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	contentType := "application/octet-stream"

	// Upload the file with FPutObject
	info, err := minioClient.PutObject(ctx, bucketName, path, f, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", path, info.Size)
	publicUrl := fmt.Sprintf("%s/%s", awsUrlPrefix, path)
	return publicUrl, nil
}
