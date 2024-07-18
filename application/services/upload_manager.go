package services

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/storage"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (vu *VideoUpload) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {
	path := strings.Split(objectPath, os.Getenv("LOCAL_STORAGE_PATH")+"/")
	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}

	defer f.Close()

	filename := path[1]
	writer := client.Bucket(vu.OutputBucket).Object(filename).NewWriter(ctx)
	writer.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(writer, f); err != nil {
		return err
	}

	if err = writer.Close(); err != nil {
		return err
	}

	return nil
}

func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {
	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	err := vu.loadPaths()
	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	for process := 0; process < concurrency; process++ {
		go vu.uploadWorker(in, returnChannel, uploadClient, ctx)
	}

	go func() {
		for i := 0; i < len(vu.Paths); i++ {
			in <- i
		}
		close(in)
	}()

	for r := range returnChannel {
		if r != "" {
			doneUpload <- r
			break
		}
	}

	return nil
}

func (vu *VideoUpload) uploadWorker(in chan int, returnChannel chan string, uploadClient *storage.Client, ctx context.Context) {
	for i := range in {
		err := vu.UploadObject(vu.Paths[i], uploadClient, ctx)
		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[i])
			log.Printf("Failed to upload %v. Error %v", vu.Paths[i], err)
			returnChannel <- err.Error()
		}
		returnChannel <- ""
	}
	returnChannel <- "upload completed"
}

func (vu *VideoUpload) loadPaths() error {
	err := filepath.Walk(
		vu.VideoPath,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				vu.Paths = append(vu.Paths, path)
			}
			return nil
		})

	if err != nil {
		return err
	}

	return nil
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, nil, err
	}

	return client, ctx, err
}
