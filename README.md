goS3
=====

goS3 is a Go package for interacting with Amazon S3. It provides a simple interface for uploading files and generating presigned URLs for downloading files.

Installation
------------

To install goS3, run the following command:

```go
go get github.com/einnovationlabs/goS3
```

Usage
-----

### Creating a Client

To create a new client, you need to provide your AWS credentials and the name of the S3 bucket you want to interact with. Here's an example:

```go
package main

import (
	"context"
	"fmt"

	"github.com/your-username/goS3"
)

func main() {
	cfg := goS3.Config{
		AccessKey:  "your-access-key",
		SecretKey: "your-secret-key",
		Region:     "your-region",
		BucketName: "your-bucket-name",
	}

	ctx := context.Background()
	client, err := goS3.New(ctx, cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
}
```

### Uploading a File

To upload a file to your S3 bucket, use the `Upload` method:

```go
fileContent := []byte("Hello, World!")
key := "path/to/your/file.txt"
if err := client.Upload(ctx, key, fileContent); err != nil {
	fmt.Println(err)
	return
}
fmt.Println("File uploaded successfully.")
```

### Generating a Presigned URL

To generate a presigned URL for downloading a file, use the `GeneratePresignedURL` method:

```go
key := "path/to/your/file.txt"
expiry := 5 * time.Minute // URL will be valid for 5 minutes
url, err := client.GeneratePresignedURL(ctx, key, expiry)
if err != nil {
	fmt.Println(err)
	return
}
fmt.Println("Presigned URL:", url)
```

License
-------

goS3 is released under the MIT License.
