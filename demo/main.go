package main

import (
  "fmt"
  "twitter"
)

func main () {


  twitter := &twitter.Client{
   ConsumerKey: "...",
   ConsumerSecret: "...",
   //OAuthToken: "...",
   //OAuthTokenSecret: "...",
  }

  fmt.Println(twitter.RequestToken("https://localhost:8347/twitter/oauth_callback"))

  
}