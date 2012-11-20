package twitter

var percentEncodingRequired = [256]bool{
	'A': true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true,
	'a': true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true,
	'0': true, true, true, true, true, true, true, true, true, true,
	'-': true,
	'.': true,
	'_': true,
	'~': true,
}

// URL encoding process described in RFC 3986, Section 2.1, as required by
// https://dev.twitter.com/docs/auth/percent-encoding-parameters
func PercentEncode(text string) (enctext string) {

	// preallocate a slice

	n := 0

	for i := 0; i < len(text); i++ {
		if percentEncodingRequired[text[i]] {
			n += 1
		} else {
			n += 3
		}
	}

	buffer := make([]byte, n, n)

	// encode the text into the preallocated slice

	j := 0

	for i := 0; i < len(text); i++ {

		b := text[i]

		if percentEncodingRequired[b] {

			buffer[j] = b
			j += 1
		} else {

			buffer[j] = '%'
			buffer[j+1] = "0123456789ABCDEF"[b>>4]
			buffer[j+2] = "0123456789ABCDEF"[b&15]

			j += 3
		}
	}

	return string(buffer)
}
