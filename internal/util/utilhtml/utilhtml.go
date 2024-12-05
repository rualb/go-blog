package utilhtml

import (
	"github.com/microcosm-cc/bluemonday"
)

// `<script>alert('XSS')</script><p onclick="steal()">Hello, <b>World</b>!</p>`

var policyWYSIWYG *bluemonday.Policy

func getPolicy() *bluemonday.Policy {

	if policyWYSIWYG == nil {
		p := bluemonday.UGCPolicy()

		p.AllowStyles("text-align").Matching(bluemonday.CellAlign).Globally() // .OnElements("p")

		policyWYSIWYG = p
	}
	return policyWYSIWYG
}

func Sanitize(unsafeHTML string) (string, error) {

	policy := getPolicy()

	safeHTML := policy.Sanitize(unsafeHTML)

	return safeHTML, nil
}
