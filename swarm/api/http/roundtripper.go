// Copyright 2017 The go-xenio Authors
// Copyright 2016 The go-ethereum Authors
//
// This file is part of the go-xenio library.
//
// The go-xenio library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-xenio library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-xenio library. If not, see <http://www.gnu.org/licenses/>.

package http

import (
	"fmt"
	"net/http"

	"github.com/xenioplatform/go-xenio/log"
)

/*
http roundtripper to register for bzz url scheme
see github.com/xenioplatform/go-xenio/issues/2040
Usage:

import (
 "github.com/xenioplatform/go-xenio/common/httpclient"
 "github.com/xenioplatform/go-xenio/swarm/api/http"
)
client := httpclient.New()
// for (private) swarm proxy running locally
client.RegisterScheme("bzz", &http.RoundTripper{Port: port})
client.RegisterScheme("bzzi", &http.RoundTripper{Port: port})
client.RegisterScheme("bzzr", &http.RoundTripper{Port: port})

The port you give the Roundtripper is the port the swarm proxy is listening on.
If Host is left empty, localhost is assumed.

Using a public gateway, the above few lines gives you the leanest
bzz-scheme aware read-only http client. You really only ever need this
if you need go-native swarm access to bzz addresses.
*/

type RoundTripper struct {
	Host string
	Port string
}

func (self *RoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	host := self.Host
	if len(host) == 0 {
		host = "localhost"
	}
	url := fmt.Sprintf("http://%s:%s/%s:/%s/%s", host, self.Port, req.Proto, req.URL.Host, req.URL.Path)
	log.Info(fmt.Sprintf("roundtripper: proxying request '%s' to '%s'", req.RequestURI, url))
	reqProxy, err := http.NewRequest(req.Method, url, req.Body)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(reqProxy)
}
