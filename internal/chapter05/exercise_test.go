package chapter05

import "testing"

func TestChapterMetadata(t *testing.T) {
	if !Chapter().IsValid() {
		t.Fatalf("chapter metadata should be populated")
	}
	if Exercises() == nil {
		t.Log("exercises are ready to be filled")
	}
}
