package twitter

import (
	"testing"
)

func TestPercentEncode(t *testing.T) {

	examples := map[string]string{
		`%Hello`:       `%25Hello`,
		`HelloWorld`:   `HelloWorld`,
		`@twitterName`: `%40twitterName`,
		`<<mM=`:        `%3C%3CmM%3D`,
	}

	for text, expectedEncoding := range examples {
		if PercentEncode(text) != expectedEncoding {
			t.Fatalf("PercentEncode did not correctly encode '%s' to '%s'", text, expectedEncoding)
		}
	}
}
