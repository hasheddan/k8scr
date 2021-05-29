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

package main

import (
	"fmt"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type pullCmd struct {
	Image string `arg:"" required:"" description:"Image to pull."`
}

func (c *pullCmd) Run(k *kong.Context, transport http.RoundTripper) error {
	ref, err := name.ParseReference(c.Image)
	if err != nil {
		return err
	}
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain), remote.WithTransport(transport))
	if err != nil {
		return err
	}
	res, err := daemon.Write(ref.Context().Tag(ref.Identifier()), img)
	fmt.Fprintf(k.Stdout, res)
	return err
}
