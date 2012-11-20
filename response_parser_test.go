package twitter

import (
	"bytes"
	"testing"
	"reflect"
)

func TestParseRequestToken(t *testing.T) {

	example := bytes.NewBufferString(`oauth_token=ljmL0Viej977Esi0nAiOMdLYuzQQ4qLGE8bLpAup4&oauth_token_secret=pFXwiOtbB8DFRnhOws4BRb3Nd4FkFv7bnQRg5NQ8&oauth_callback_confirmed=true`)

	expected := &RequestToken{
		OAuthToken: `ljmL0Viej977Esi0nAiOMdLYuzQQ4qLGE8bLpAup4`,
		OAuthTokenSecret: `pFXwiOtbB8DFRnhOws4BRb3Nd4FkFv7bnQRg5NQ8`,
		OAuthCallbackConfirmed: true,
	}

	actual, _ := parseRequestToken(example)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("request token parsing failed:\nGOT:\n%s\nEXPECTED:\n%s", actual, expected)
	}
}