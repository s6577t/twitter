package twitter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"strings"
)

type Client struct {
	ConsumerKey string
	ConsumerSecret string
	OAuthToken string
	OAuthTokenSecret string
}

func (tc *Client) RequestParameters () (params map[string]string, err error) {
	
	var nonce string

	if nonce, err = NOnce(); err != nil {
		return
	}

	params = map[string]string{
		"oauth_consumer_key": tc.ConsumerKey,
		"oauth_nonce": nonce,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_version": "1.0",
	}

	return
}

func (tc *Client) AuthorizedRequestParameters () (params map[string]string, err error) {
	
	if params, err = tc.RequestParameters(); err != nil {
		return
	}

	params["oauth_token"] = tc.OAuthToken

	return
}

func (tc *Client) Search (q string) (json string, err error) {

	var params map[string]string

	if params, err = tc.AuthorizedRequestParameters(); err != nil {
		return
	}

	params["q"] = q

	baseUrl := "https://api.twitter.com/1.1/search/tweets.json"
		
	url := func () string {
		
		queryParameters := url.Values{}
		queryParameters.Add("q", q)
		
		return fmt.Sprintf("%s?%s", baseUrl, queryParameters.Encode())
	}()

	client := &http.Client{}

	// a new request, not one that is parsed, can't fail
	request, _ := http.NewRequest("GET", url, nil)
	
	request.Header.Add("Authorization", OAuthHeader("GET", baseUrl, params, tc.ConsumerSecret, tc.OAuthTokenSecret))
	
	var response *http.Response
	
	if response, err = client.Do(request); err != nil {
		return
	}

	defer response.Body.Close()
	
	var responseBodyBytes []byte

	if responseBodyBytes, err = ioutil.ReadAll(response.Body); err != nil {
		return
	} 

	json = string(responseBodyBytes)

	return
}

func (tc *Client) RequestToken (callbackUrl string) (requestToken *RequestToken, err error) {
	
	var params map[string]string

	if params, err = tc.RequestParameters(); err != nil {
		return
	}

	params["oauth_callback"] = callbackUrl

	baseUrl := `https://api.twitter.com/oauth/request_token`
		
	url := func () string {
		
		queryParameters := url.Values{}
		queryParameters.Add("oauth_callback", callbackUrl)
		
		return fmt.Sprintf("%s?%s", baseUrl, queryParameters.Encode())
	}()

	client := &http.Client{}

	// a new request, not one that is parsed, can't fail
	request, _ := http.NewRequest("POST", url, nil)
	
	request.Header.Add("Authorization", OAuthHeader("POST", baseUrl, params, tc.ConsumerSecret, ""))
	
	var response *http.Response
	
	if response, err = client.Do(request); err != nil {
		return
	}

	defer response.Body.Close()
	
	requestToken, err = parseRequestToken(response.Body)
	
	return
}

func (tc *Client) AccessToken (requestToken *RequestToken, callbackParams *OAuthCallbackParameters) (accessToken *AccessToken, err error) {
	
	// copy the receiver and assign the tokens from the request token, we'll be using authorized parameters
	tc = &*tc
	tc.OAuthToken = requestToken.OAuthToken
	tc.OAuthTokenSecret = requestToken.OAuthTokenSecret

	var params map[string]string

	if params, err = tc.AuthorizedRequestParameters(); err != nil {
		return
	}

	params["oauth_verifier"] = callbackParams.OAuthVerifier

	url := `https://api.twitter.com/oauth/access_token`
	
	client := &http.Client{}

	// a new request, not one that is parsed, can't fail
	request, _ := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("oauth_verifier=%s", PercentEncode(callbackParams.OAuthVerifier))))
	
	request.Header.Add("Authorization", OAuthHeader("POST", url, params, tc.ConsumerSecret, tc.OAuthTokenSecret))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var response *http.Response
	
	if response, err = client.Do(request); err != nil {
		return
	}

	defer response.Body.Close()
	
	accessToken, err = parseAccessToken(response.Body)

	return
}

