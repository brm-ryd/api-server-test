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
func (f *AzureStorage) Get(name string, out io.Writer) (found bool, tag interface{}, err error) {

}
