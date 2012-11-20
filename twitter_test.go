// testing examples from twitter documentation:
// https://dev.twitter.com/docs/auth/authorizing-request
// https://dev.twitter.com/docs/auth/creating-signature
package twitter

import (
	"testing"
)

var exampleParameters = map[string]string{
	"status":                 "Hello Ladies + Gentlemen, a signed OAuth request!",
	"include_entities":       "true",
	"oauth_consumer_key":     "xvz1evFS4wEEPTGEFPHBog",
	"oauth_nonce":            "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
	"oauth_signature_method": "HMAC-SHA1",
	"oauth_timestamp":        "1318622958",
	"oauth_token":            "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
	"oauth_version":          "1.0",
}

var (
	exampleConsumerSecret   = `kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw`
	exampleOAuthTokenSecret = `LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE`
)

func TestOAuthSignatureParameterStringWithoutOAuthToken(t *testing.T) {

	exampleParametersExceptOAuthToken := map[string]string{}
	for k, v := range exampleParameters {
		if k != "oauth_token" {
			exampleParametersExceptOAuthToken[k] = v 
		}
	}

	expectedParameterString := `include_entities=true&oauth_consumer_key=xvz1evFS4wEEPTGEFPHBog&oauth_nonce=kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg&oauth_signature_method=HMAC-SHA1&oauth_timestamp=1318622958&oauth_version=1.0&status=Hello%20Ladies%20%2B%20Gentlemen%2C%20a%20signed%20OAuth%20request%21`
	parameterString := oAuthSignatureParameterString(exampleParametersExceptOAuthToken)

	if parameterString != expectedParameterString {
		t.Fatalf("parameter string encoding with missing oauth_token failed: \nGOT:\n%s\nEXPECTED:\n%s\n", parameterString, expectedParameterString)
	}
}

func TestOAuthSignatureParameterString(t *testing.T) {

	expectedParameterString := `include_entities=true&oauth_consumer_key=xvz1evFS4wEEPTGEFPHBog&oauth_nonce=kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg&oauth_signature_method=HMAC-SHA1&oauth_timestamp=1318622958&oauth_token=370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb&oauth_version=1.0&status=Hello%20Ladies%20%2B%20Gentlemen%2C%20a%20signed%20OAuth%20request%21`
	parameterString := oAuthSignatureParameterString(exampleParameters)

	if parameterString != expectedParameterString {
		t.Fatalf("parameter string encoding failed: \nGOT:\n%s\nEXPECTED:\n%s\n", parameterString, expectedParameterString)
	}
}

func TestOAuthSignatureBaseString(t *testing.T) {

	expectedBaseString := `POST&https%3A%2F%2Fapi.twitter.com%2F1%2Fstatuses%2Fupdate.json&include_entities%3Dtrue%26oauth_consumer_key%3Dxvz1evFS4wEEPTGEFPHBog%26oauth_nonce%3DkYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1318622958%26oauth_token%3D370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb%26oauth_version%3D1.0%26status%3DHello%2520Ladies%2520%252B%2520Gentlemen%252C%2520a%2520signed%2520OAuth%2520request%2521`
	baseString := OAuthSignatureBaseString(`POST`, `https://api.twitter.com/1/statuses/update.json`, exampleParameters)

	if baseString != expectedBaseString {
		t.Fatalf("base string encoding failed: \nGOT:\n%s\nEXPECTED:\n%s\n", baseString, expectedBaseString)
	}
}

func TestOAuthSignature(t *testing.T) {

	// NOTE: the oauth signature should be base64 encoded
	expectedSignature := `tnnArxj06cWHq44gCs1OSKk/jLY=`
	signature := OAuthSignature(
		OAuthSignatureBaseString(`POST`, `https://api.twitter.com/1/statuses/update.json`, exampleParameters),
		exampleConsumerSecret,
		exampleOAuthTokenSecret)

	if signature != expectedSignature {
		t.Fatalf("signature failed: \nGOT:\n%s\nEXPECTED:\n%s\n", signature, expectedSignature)
	}
}

func TestOAuthHeader(t *testing.T) {

	expectedHeader := `OAuth oauth_consumer_key="xvz1evFS4wEEPTGEFPHBog", oauth_nonce="kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg", oauth_signature="tnnArxj06cWHq44gCs1OSKk%2FjLY%3D", oauth_signature_method="HMAC-SHA1", oauth_timestamp="1318622958", oauth_token="370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb", oauth_version="1.0"`
	header := OAuthHeader(`POST`, `https://api.twitter.com/1/statuses/update.json`, 
		exampleParameters, exampleConsumerSecret, exampleOAuthTokenSecret)

	if header != expectedHeader {
		t.Fatalf("failed to create valid authorization header: \nGOT:\n%s\nEXPECTED:\n%s\n", header, expectedHeader)
	}
}

