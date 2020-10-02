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
	}

	resp, err := azblob.UploadStreamToBlockBlob(context.TODO(), in, blockBlobURL, azblob.UploadStreamToBlockBlobOptions{
		AccessConditions: accessConditions,
	})
	if err != nil {
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
		MaxRetryRequests: 3,
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
