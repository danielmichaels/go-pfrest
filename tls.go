package pfrest

import (
	"crypto/tls"
	"net/http"
)

func tlsClient(skipVerify bool) *http.Client {
	if !skipVerify {
		return http.DefaultClient
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //nolint:gosec // user-configured for self-signed pfSense certs
			},
		},
	}
}
