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
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/hasheddan/k8scr/internal/transport"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/net"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"
)

var _ = kong.Must(&cli)

var cli struct {
	Push pushCmd `cmd:"" description:"Push an OCI image through the API server."`
	Pull pushCmd `cmd:"" description:"Pull an OCI image through the API server."`

	Kubeconfig string `type:"existingfile" help:"Override default kubeconfig path."`
	Namespace  string `short:"n" default:"default" help:"Namespace of registry Pod."`
	Registry   string `short:"r" default:"k8scr" help:"Name of registry Pod."`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("k8scr"),
		kong.Description("Push and pull images through the Kubernetes API server."),
		kong.UsageOnError())

	// Fetch kubeconfig
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.ExplicitPath = cli.Kubeconfig
	kubeconfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{}).ClientConfig()
	ctx.FatalIfErrorf(err)

	// Build clientset for constructing URL
	clientset, err := kubernetes.NewForConfig(kubeconfig)
	ctx.FatalIfErrorf(err)
	r := clientset.CoreV1().RESTClient().Get().Namespace(cli.Namespace).Resource("pods").SubResource("proxy").Name(net.JoinSchemeNamePort("http", cli.Registry, "80"))

	// Build client for custom transport
	gv := corev1.SchemeGroupVersion
	kubeconfig.GroupVersion = &gv
	kubeconfig.APIPath = "/api"
	kubeconfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	client, err := rest.RESTClientFor(kubeconfig)
	ctx.FatalIfErrorf(err)

	// Wrap transport to rewrite paths
	w := &transport.Wrapper{
		Client: client.Client,
		URL:    r.URL(),
	}
	ctx.BindTo(w, (*http.RoundTripper)(nil))
	ctx.FatalIfErrorf(ctx.Run())
}
