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
	// Create blob URL
	u, err := url.Parse(f.storageURL + "/" + name)
	if err != nil {
		return
	}
	blockBlobURL := azblob.NewBlockBlobURL(*u, f.storagePipeline)

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

}
