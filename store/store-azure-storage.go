package store

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
)

// Store files on Azure Blob Storage
type AzureStorage struct {
	storageAccountName string
	storageContainer   string
	storagePipeline    pipeline.Pipeline
	storageURL         string
}

// Actions (like in store modules interfaces)
func (f *AzureStorage) Init(connection string) error {
	r := regexp.MustCompile("^(azureblob|azure):([a-z0-9][a-z0-9-]{2,62})$")
	match := r.FindStringSubmatch(connection)
	if match == nil || len(match) != 3 {
		return errors.New("invalid connection string for blob storage")
	}
	f.storageContainer = match[2]

	// Get storage account name and key from the environment
	name := os.Getenv("AZURE_STORAGE_ACCOUNT")
	key := os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	if name == "" || key == "" {
		return errors.New("environmental variables AZURE_STORAGE_ACCOUNT and AZURE_STORAGE_ACCESS_KEY are not defined")
	}
	f.storageAccountName = name

	// Endpoint storage
	f.storageURL = fmt.Sprintf("https://%s.blob.core.windows.net/%s", f.storageAccountName, f.storageContainer)

	// Azure Storage auth
	credential, err := azblob.NewSharedKeyCredential(f.storageAccountName, key)
	if err != nil {
		return err
	}
	f.storagePipeline = azblob.NewPipeline(credential, azblob.PipelineOptions{
		Retry: azblob.RetryOptions{
			MaxTries: 5,
		},
	})

	return nil

}

func (f *AzureStorage) Set(name string, in io.Reader, tag interface{}) (tagOut interface{}, err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}

	// Create blob URL
	u, err := url.Parse(f.storageURL + "/" + name)
	if err != nil {
		return
	}
	blockBlobURL := azblob.NewBlockBlobURL(*u, f.storagePipeline)
	var accessConditions azblob.BlobAccessConditions
	if tag == nil {
		// If no blob at that path yet will success
		accessConditions = azblob.BlobAccessConditions{
			ModifiedAccessConditions: azblob.ModifiedAccessConditions{
				IfNoneMatch: "*",
			},
		}
	} else {
		// Uploads success only if file hasn't been modified since last download
		accessConditions = azblob.BlobAccessConditions{
			ModifiedAccessConditions: azblob.ModifiedAccessConditions{
				IfMatch: *tag.(*azblob.ETag),
			},
		}
	}

	resp, err := azblob.UploadStreamToBlockBlob(context.TODO(), in, blockBlobURL, azblob.UploadStreamToBlockBlobOptions{
		BufferSize:       3 * 1024 * 1024,
		MaxBuffers:       2,
		AccessConditions: accessConditions,
	})
	if err != nil {
		if stgErr, ok := err.(azblob.StorageError); !ok {
			return nil, fmt.Error("network error - in uploading file: %s", err.Error())

		} else {
			return nil, fmt.Error("storage azure failed - uploading the file: %s", stgErr.Response().Status)
		}
	}

	// Get the ETag
	tagObj := resp.ETag()
	tagOut = &tagObj

	return tagOut, nil

}

func (f *AzureStorage) Get(name string, out io.Writer) (found bool, tag interface{}, err error) {
	if name == "" {
		err = errors.New("Empty name")
		return
	}

	found = true

	// Create blob URL
	u, err := url.Parse(f.storageURL + "/" + name)
	if err != nil {
		return
	}
	blockBlobURL := azblob.NewBlockBlobURL(*u, f.storagePipeline)

	// Download file
	resp, err := blockBlobURL.Download(context.TODO(), 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	if err != nil {
		if stgErr, ok := err.(azblob.StorageError); !ok {
			err = fmt.Errorf("network error while downloading the file: %s", err.Error())
		} else {
			// Blob checker
			if stgErr.Response().StatusCode == http.StatusNotFound {
				err = nil
				found = false
				return
			}
			err = fmt.Errorf("azure Storage error while downloading the file: %s", stgErr.Response().Status)
		}
		return
	}
	body := resp.Body(azblob.RetryReaderOptions{
		MaxRetryRequests: 5,
	})
	defer body.Close()

	// Empty file exist
	if resp.ContentLength() == 0 {
		body.Close()
		found = true
		return
	}

	// Copy response
	_, err = io.Copy(out, body)
	if err != nil {
		return
	}

	// Get ETag
	tagObj := resp.ETag()
	tag = &tagObj

	return

}

func (f *AzureStorage) Delete(name string, tag interface{}) (err error) {
	if name == "" {
		err = errors.New("empty name")
		return
	}

	// Create blob URL
	u, err := url.Parse(f.storageURL + "/" + name)
	if err != nil {
		return
	}
	blockBlobURL := azblob.NewBlockBlobURL(*u, f.storagePipeline)

	var accessConditions azblob.BlobAccessConditions
	if tag != nil {
		accessConditions = azblob.BlobAccessConditions{
			ModifiedAccessConditions: azblob.ModifiedAccessConditions{
				IfMatch: *tag.(*azblob.ETag),
			},
		}
	}

	// Delete blob
	_, err = blockBlobURL.Delete(context.TODO(), azblob.DeleteSnapshotsOptionInclude, accessConditions)
	return

}
