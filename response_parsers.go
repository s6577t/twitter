package twitter

import (
	"io"
	"strconv"
	"net/url"
	"io/ioutil"
)

func parseForValues (reader io.Reader) (values url.Values, err error) {
	
	var bytes []byte
	
	if bytes, err = ioutil.ReadAll(reader); err != nil {
		return nil, err
	}

	return url.ParseQuery(string(bytes))
}

func parseRequestToken (reader io.Reader) (t *RequestToken, err error) {
	
	var values url.Values
	
	if values, err = parseForValues(reader); err != nil {
		return
	}

	var callbackConfirmed bool

	if callbackConfirmed, err = strconv.ParseBool(values["oauth_callback_confirmed"][0]); err != nil {
		return
	}

	return &RequestToken{
		OAuthToken: values["oauth_token"][0],
		OAuthTokenSecret: values["oauth_token_secret"][0],
		OAuthCallbackConfirmed: callbackConfirmed,
	}, nil
}

func parseAccessToken (reader io.Reader) (t *AccessToken, err error) {
	
	var values url.Values
	
	if values, err = parseForValues(reader); err != nil {
		return
	}

	return &AccessToken{
		OAuthToken: values["oauth_token"][0],
		OAuthTokenSecret: values["oauth_token_secret"][0],
		ScreenName: values["screen_name"][0],
		TwitterUserId: values["user_id"][0],
	}, nil
}