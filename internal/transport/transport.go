// Copyright 2021 The k8src Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transport

import (
	"net/http"
	"net/url"
	"path"
)

// Wrapper implements RoundTripper by appending request path and parameters to
// its URL.
type Wrapper struct {
	Client *http.Client
	URL    *url.URL
}

// RoundTrip implements the http.RoundTripper interface.
func (w *Wrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	params, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(w.URL.String())
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(w.URL.Path, req.URL.Path)
	u.RawQuery = params.Encode()
	req.URL = u
	return w.Client.Transport.RoundTrip(req)
}
