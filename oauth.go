package twitter

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"sort"
	"strings"
)

var oAuthParameters = [7]string{"oauth_consumer_key", "oauth_nonce", "oauth_signature", "oauth_signature_method", "oauth_timestamp", "oauth_token", "oauth_version"}

func oAuthSignatureParameterString(parameters map[string]string) string {

	var pairs sort.StringSlice = make([]string, len(parameters))

	// create a list of pairs
	// anonymous function avoids too much variable declaration noise in this scope
	func() {
		i := 0
		for name, value := range parameters {
			// add the pair unless it's a blank oauth_token
			if name != "oauth_token" || value != "" {
				pairs[i] = fmt.Sprintf("%s=%s", PercentEncode(name), PercentEncode(value))
			}			
			i++
		}
	}()

	pairs.Sort()

	return strings.Join(pairs, "&")
}

// baseUrl MUST NOT include the query string of the final url
func OAuthSignatureBaseString(httpMethod string, baseUrl string, parameters map[string]string) string {
	return fmt.Sprintf("%s&%s&%s",
		httpMethod,
		PercentEncode(baseUrl),
		PercentEncode(
			oAuthSignatureParameterString(parameters)))
}

func OAuthSignature(baseString, consumerSecret, oAuthTokenSecret string) string {

	signingKey := fmt.Sprintf("%s&%s", PercentEncode(consumerSecret), PercentEncode(oAuthTokenSecret))

	hash := hmac.New(sha1.New, []byte(signingKey))

	io.WriteString(hash, baseString)

	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func OAuthHeader(httpMethod string, baseUrl string, parameters map[string]string, consumerSecret string, oAuthTokenSecret string) string {
	
	buffer := bytes.NewBufferString("OAuth ")

	for i, key := range oAuthParameters {
		
		if key == "oauth_token" && parameters["oauth_token"] == "" {
			continue
		}

		var value string

		if key == "oauth_signature" {
			value = OAuthSignature(OAuthSignatureBaseString(httpMethod, baseUrl, parameters), consumerSecret, oAuthTokenSecret)
		} else {
			value = parameters[key]
		}

		pair := fmt.Sprintf("%s=\"%s\"", PercentEncode(key), PercentEncode(value))

		buffer.WriteString(pair)

		if i < len(oAuthParameters) - 1 {
			buffer.WriteString(", ")
		}
	}

	return buffer.String()
}

func NOnce () (string, error) {
	
	buffer := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, buffer); err != nil {
		return "", err
	}
	
	return base64.StdEncoding.EncodeToString(buffer), nil
} 