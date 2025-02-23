package xxhash

import (
	"errors"
	"fmt"

	"github.com/cespare/xxhash"
	"github.com/snapp-incubator/bucketizer"
)

type XXHASHBucketizer struct {
	Buckets      []bucketizer.Bucket
	weightSum    uint64
	bucketRanges []uint64
	seed         []byte
}

func (b XXHASHBucketizer) bucketBytes(value []byte) (int, error) {
	a := append(value, b.seed...)
	hashDigest := xxhash.Sum64(a)
	reminder := hashDigest % b.weightSum
	n := len(b.bucketRanges)
	for i := 0; i < n-1; i++ {
		if reminder >= b.bucketRanges[i] && reminder < b.bucketRanges[i+1] {
			return i, nil
		}
	}
	return 0, errors.New("invalid reminder value")
}

// Bucket is a generic method that accepts any type. It converts the input to []byte based on its underlying type:
// - For []byte, the slice is used directly.
// - For string, it is converted to []byte.
// - For integers, the value is formatted using "%d".
// - For floats, the value is formatted using "%g".
// For unsupported types, it returns an error.
func (b XXHASHBucketizer) Bucket(value interface{}) (int, error) {
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		data = []byte(fmt.Sprintf("%d", v))
	case float32, float64:
		data = []byte(fmt.Sprintf("%g", v))
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
	return b.bucketBytes(data)
}

// NewXXHASHBucketizer initializes and returns a new XXHASHBucketizer instance.
// It calculates the cumulative bucket ranges based on each bucket's weight.
func NewXXHASHBucketizer(seed string, buckets ...bucketizer.Bucket) XXHASHBucketizer {
	var sumOfWeights uint64
	bucketRanges := make([]uint64, 0, len(buckets)+1)
	bucketRanges = append(bucketRanges, 0)
	for i, bucket := range buckets {
		sumOfWeights += uint64(bucket.Weight)
		bucketRanges = append(bucketRanges, bucketRanges[i]+uint64(bucket.Weight))
	}
	return XXHASHBucketizer{
		Buckets:      buckets,
		weightSum:    sumOfWeights,
		bucketRanges: bucketRanges,
		seed:         []byte(seed),
	}
}
