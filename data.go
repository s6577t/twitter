package twitter

type RequestToken struct {
	OAuthToken string
	OAuthTokenSecret string
	OAuthCallbackConfirmed bool
}

type OAuthCallbackParameters struct {
	OAuthToken string
	OAuthVerifier string
}

type AccessToken struct {
	OAuthToken string
	OAuthTokenSecret string
	ScreenName string
	TwitterUserId string
}