package adidasv2

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"umbrella/internal/utils/helpers"

	http "github.com/vimbing/fhttp"
)

func AddIstanaHeaders(headers http.Header) http.Header {
	istanaValue := helpers.RandomString(16, true, false, true)

	istanaHeaders := map[string][]string{
		"x-instana-l": {fmt.Sprintf("1,correlationType=web;correlationId=%s", istanaValue)},
		"x-instana-s": {istanaValue},
		"x-instana-t": {istanaValue},
	}

	for k, v := range istanaHeaders {
		headers[k] = v
	}

	return headers
}

func (r *Region) GetAuthorityHeader() string {
	return fmt.Sprintf("www.adidas.%s", strings.ToLower(r.RegionCode))
}

func (r *Region) GetOriginHeader() string {
	return fmt.Sprintf("https://www.adidas.%s", strings.ToLower(r.RegionCode))
}

func (r *Region) GetRefererHeader(path string) string {
	return fmt.Sprintf("https://www.adidas.%s/%s", strings.ToLower(r.RegionCode), path)
}

func GetCacheBypassPayload() string {
	return fmt.Sprintf("pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7,la;q=0.6,de;q=0.5,%s", helpers.RandomString(helpers.RandomInt(0, 45), true, false, true))
}

func HandleModeRetryError(err error, config *Config) {
	if errors.Is(err, WaitError{}) {
		config.DefaultConfig.Log.Yellow("Waiting 25 secs...")
		time.Sleep(time.Second * 25)
	} else if errors.Is(err, LongWaitError{}) {
		config.DefaultConfig.Log.Yellow("Waiting 60 secs...")
		time.Sleep(time.Minute * 1)
	}
}
