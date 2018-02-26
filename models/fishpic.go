package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var LOCAL_CACHE = envy.Get("LOCAL_CACHE", "false")

type Fishpic struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Species   string    `json:"species"`
	Location  string    `json:"location"`
}

// String is not required by pop and may be deleted
func (f Fishpic) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Fishpics is not required by pop and may be deleted
type Fishpics []Fishpic

// String is not required by pop and may be deleted
func (f Fishpics) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

func (f *Fishpic) Init() {
	newUUID, err := uuid.NewV4()
	if err != nil {
		errors.WithStack(err)
	}
	f.ID = newUUID

	currTime := time.Now().UTC()
	f.CreatedAt = currTime
	f.UpdatedAt = currTime
}

func (f *Fishpic) CreateAndSave() error {
	f.Init()
	err := f.ValidateAndSave()
	if err != nil {
		return err
	}

	return nil
}

func (f *Fishpic) ValidateAndSave() error {
	bytes, err := json.Marshal(f)
	if err != nil {
		return err
	}

	if LOCAL_CACHE == "true" {
		err = ioutil.WriteFile("gofishdb/"+f.ID.String(), bytes, 0600)
		if err != nil {
			return err
		}
	} else {
		// write data to s3
		err = Uploader("gofishdb/"+f.ID.String(), bytes)
		if err != nil {
			return err
		}
	}
	return nil
}

func Uploader(filename string, data []byte) error {
	// https://jto.nyc3.digitaloceanspaces.com
	// The session the S3 Uploader will use
	endpoint := "nyc3.digitaloceanspaces.com"
	region := "nyc3"
	myBucket := "jto"

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: &endpoint,
		Region:   &region,
	}))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(myBucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	log.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))

	return nil
}

func Downloader(filename string) ([]byte, error) {
	// https://jto.nyc3.digitaloceanspaces.com
	// The session the S3 Uploader will use
	endpoint := "nyc3.digitaloceanspaces.com"
	region := "nyc3"
	myBucket := "jto"

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: &endpoint,
		Region:   &region,
	}))

	// Create an uploader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(tmpfile, &s3.GetObjectInput{
		Bucket: aws.String(myBucket),
		Key:    aws.String(filename),
	})

	data := make([]byte, n)

	if err != nil {
		return data, fmt.Errorf("failed to download file, %v", err)
	}

	_, err = tmpfile.Read(data)
	if err != nil {
		return data, err
	}

	log.Printf("file downloaded, %d bytes\n", n)

	return data, nil
}

func Lister(prefix string) ([]string, error) {
	// https://jto.nyc3.digitaloceanspaces.com
	// The session the S3 Uploader will use
	endpoint := "nyc3.digitaloceanspaces.com"
	region := "nyc3"
	myBucket := "jto"

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: &endpoint,
		Region:   &region,
	}))

	client := s3.New(sess)

	params := s3.ListObjectsV2Input{
		Bucket: &myBucket,
		Prefix: &prefix,
	}

	// Example iterating over pages of a ListObjectsV2 operation.
	var ids []string
	err := client.ListObjectsV2Pages(&params,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			objects := page.Contents
			for _, object := range objects {
				key := *object.Key
				key = strings.Replace(key, prefix, "", 1)

				if key == "" {
					continue
				}

				ids = append(ids, key)
			}

			return !lastPage
		})

	return ids, err
}

func (f *Fishpic) Find(id string) error {
	if LOCAL_CACHE == "true" {
		data, err := ioutil.ReadFile("gofishdb/" + id)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, f)
		if err != nil {
			return err
		}

		return nil

	}

	data, err := Downloader("gofishdb/" + id)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, f)
	if err != nil {
		return err
	}

	return nil
}

func (f *Fishpic) List() (*[]string, error) {
	var ids []string

	if LOCAL_CACHE == "true" {
		files, err := ioutil.ReadDir("gofishdb/")
		if err != nil {
			return &ids, err
		}

		ids = make([]string, len(files))
		for i, file := range files {
			ids[i] = file.Name()
		}

		return &ids, err
	}

	ids, err := Lister("gofishdb/")
	if err != nil {
		return &ids, err
	}

	return &ids, nil
}
