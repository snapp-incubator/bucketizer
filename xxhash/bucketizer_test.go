package xxhash

import (
	"testing"

	"github.com/snapp-incubator/bucketizer"
)

func TestBucketizerSupportedTypes(t *testing.T) {
	buckets := []bucketizer.Bucket{
		{Weight: 1},
		{Weight: 2},
		{Weight: 3},
	}
	bz := NewXXHASHBucketizer("test-seed", buckets...)

	// Define test cases for supported types.
	testCases := []struct {
		name  string
		value any
	}{
		{"string", "hello"},
		{"int", 42},
		{"int8", int8(7)},
		{"int16", int16(300)},
		{"int32", int32(12345)},
		{"int64", int64(67890)},
		{"float32", float32(3.14)},
		{"float64", 6.28},
		{"[]byte", []byte("bucket")},
	}

	// Run each test case.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			index, err := bz.Bucket(tc.value)
			if err != nil {
				t.Errorf("unexpected error for type %T: %v", tc.value, err)
			}
			if index < 0 || index >= len(buckets) {
				t.Errorf("bucket index %d out of range for type %T", index, tc.value)
			}
		})
	}
}

func TestBucketizerUnsupportedType(t *testing.T) {
	// Define sample buckets.
	buckets := []bucketizer.Bucket{
		{Weight: 1},
		{Weight: 2},
		{Weight: 3},
	}

	bz := NewXXHASHBucketizer("test-seed", buckets...)

	type custom struct {
		value int
	}

	// Expect an error when using an unsupported type.
	_, err := bz.Bucket(custom{value: 10})
	if err == nil {
		t.Errorf("expected error for unsupported type, got nil")
	}
}

// TestNewXXHASHBucketizerWeights verifies that the NewXXHASHBucketizer function
// computes the cumulative weight sum and bucket ranges correctly.
func TestNewXXHASHBucketizerWeights(t *testing.T) {
	buckets := []bucketizer.Bucket{
		{Name: "A", Weight: 1},
		{Name: "B", Weight: 2},
		{Name: "C", Weight: 3},
	}
	bz := NewXXHASHBucketizer("test", buckets...)

	// Expected weightSum is 6.
	if bz.weightSum != 6 {
		t.Errorf("expected weightSum to be 6, got %d", bz.weightSum)
	}

	// Expected bucketRanges is [0, 1, 3, 6].
	expectedRanges := []uint64{0, 1, 3, 6}
	if len(bz.bucketRanges) != len(expectedRanges) {
		t.Fatalf("expected bucketRanges length %d, got %d", len(expectedRanges), len(bz.bucketRanges))
	}
	for i, expected := range expectedRanges {
		if bz.bucketRanges[i] != expected {
			t.Errorf("bucketRanges[%d]: expected %d, got %d", i, expected, bz.bucketRanges[i])
		}
	}
}
