package onion

import (
	"context"
	"encoding/json"
	"github.com/cretz/bine/tor"
	"github.com/go-errors/errors"
	"net/http"
	"time"
)

type checkApiResponse struct {
	IsTor bool   `json:"IsTor"`
	IP    string `json:"IP"`
}

func Check(tor *tor.Tor) (bool, error) {
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()

	dialer, err := tor.Dialer(dialCtx, nil)
	if err != nil {
		return false, err
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: dialer.DialContext,
		},
	}

	resp, err := httpClient.Get("https://check.torproject.org/api/ip")
	if err != nil {
		return false, err
	}

	res := checkApiResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return false, errors.Errorf("unable to read check api response: %v", err)
	}

	return res.IsTor, nil
}
