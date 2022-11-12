package models

import (
	"net/url"
)

// RedirectUrl ProxyTarget defines the upstream target.
type RedirectUrl struct {
	URL *url.URL
}
