# Bucketizer

A generic Go bucketizer that leverages xxhash for fast hashing and supports weighted bucketing via generics. This
package lets you distribute various types of input values—such as strings, numbers, or byte slices—into buckets defined
by custom weights.

### Overview

The bucketizer package uses a fast, non-cryptographic hash (xxhash) combined with a user-defined seed to determine the
bucket assignment for a given input.

Features

- Fast Hashing: Uses xxhash for high performance.
- Weighted Distribution: Buckets are defined with weights, so each bucket’s chance is proportional to its weight.

### Installation

To get the package, run:

```shell
go get github.com/snapp-incubator/bucketizer
```

### Usage

Here’s an example demonstrating how to create and use the bucketizer:

```go
package main

import (
	"fmt"
	"log"

	"github.com/snapp-incubator/bucketizer"
	"github.com/snapp-incubator/bucketizer/xxhash"
)

func main() {
	// Define your buckets with weights.
	buckets := []bucketizer.Bucket{
		{Weight: 1},
		{Weight: 2},
		{Weight: 3},
	}

	// Initialize the bucketizer with a seed.
	b := xxhash.NewXXHASHBucketizer("my-seed", buckets...)

	// Use the generic Bucket method to bucket different data types.
	index, err := b.Bucket("Hello World")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bucket index for 'Hello World': %d\n", index)

	// You can also bucket numeric types:
	intIndex, err := b.Bucket(123)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Bucket index for 123: %d\n", intIndex)
}

```

In this example, the bucketizer is created using a seed and a set of buckets defined by weight. The `Bucket` method
converts the input into a byte slice (using appropriate formatting for each supported type) and
computes the bucket index based on the hash.

### Contributing

Contributions are welcome! If you find any issues or have improvements, please feel free to open an issue or submit a
pull request.

### License

This project is licensed under the MIT License. See the LICENSE file for details.

