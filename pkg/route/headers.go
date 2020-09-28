// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package route

import (
	"fmt"
	"net/http"
	"os"

	moovhttp "github.com/moov-io/base/http"
	"github.com/moov-io/customers/internal/util"
)

var (
	namespaceHeaderKey    = util.Or(os.Getenv("NAMESPACE_HEADER"), "X-Namespace")
	organizationHeaderKey = util.Or(os.Getenv("ORGANIZATION_HEADER"), "X-Organization")
)

// GetNamespace returns the X-Namespace header value and writes an error to w if the namespace header is empty
func GetNamespace(w http.ResponseWriter, r *http.Request) string {
	return getHeader(namespaceHeaderKey, w, r)
}

// GetOrganization returns the X-Organization header value and writes an error to w if the organization header is empty
func GetOrganization(w http.ResponseWriter, r *http.Request) string {
	return getHeader(organizationHeaderKey, w, r)
}

// getHeader retrieves the value associated with the header key
// An error will be written to w and an empty string returned if the header is empty.
func getHeader(headerKey string, w http.ResponseWriter, r *http.Request) string {
	h := r.Header.Get(headerKey)
	if h == "" {
		moovhttp.Problem(w, fmt.Errorf("missing %s header", headerKey))
	}
	return h
}
