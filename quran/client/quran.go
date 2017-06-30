// Code generated by goagen v1.2.0-dirty, DO NOT EDIT.
//
// API "quran": quran Resource Client
//
// Command:
// $ goagen
// --design=github.com/ccdatatraits/quran/design
// --out=$(GOPATH)/src/github.com/ccdatatraits/quran
// --version=v1.2.0-dirty

package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ShowQuranPath computes a request path to the show action of quran.
func ShowQuranPath(suraID int, ayaID int) string {
	param0 := strconv.Itoa(suraID)
	param1 := strconv.Itoa(ayaID)

	return fmt.Sprintf("/quran/%s/%s", param0, param1)
}

// Get aya by sura & aya
func (c *Client) ShowQuran(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewShowQuranRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewShowQuranRequest create the request corresponding to the show action endpoint of the quran resource.
func (c *Client) NewShowQuranRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}
