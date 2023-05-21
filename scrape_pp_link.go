package adidasv2

import (
	"bytes"
	"encoding/json"
	"fmt"
	httpclient "umbrella/internal/http_client"
	definederrors "umbrella/internal/utils/defined_errors"
	"umbrella/internal/utils/helpers"

	http "github.com/vimbing/fhttp"
)

func ScrapePaypalLink(config *Config) (string, error) {
	state := config.TaskStates.Order

	config.DefaultConfig.Log.SetState(state.Name)

	for i := 0; i < state.Retry; i++ {
		config.DefaultConfig.Log.Yellow("Scraping paypal link...")

		params := helpers.CreateParams(map[string]string{
			"paymentBrand":     "PAYPAL",
			"shopperResultUrl": fmt.Sprintf("https://www.adidas.pl/payment/callback/async/PAYPAL/%s/aci", config.Resources.BasketID),
			"forceUtf8":        "&#9760",
			"shopOrigin":       "https://www.adidas.pl",
		})

		req, err := http.NewRequest("POST", fmt.Sprintf("https://eu-prod.oppwa.com/v1/checkouts/%s/payment", config.Resources.CheckoutID), bytes.NewBufferString(params))

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_REQUEST_CREATING_ERROR)
			continue
		}

		req.Header = http.Header{
			"authority":          {"eu-prod.oppwa.com"},
			"accept":             {"application/json"},
			"accept-language":    {"pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7,la;q=0.6,de;q=0.5"},
			"content-type":       {"application/x-www-form-urlencoded;charset=UTF-8"},
			"dnt":                {"1"},
			"origin":             {"https://eu-prod.oppwa.com"},
			"referer":            {"https://eu-prod.oppwa.com/v1/internalRequestIframe.html"},
			"sec-ch-ua":          {"\"Not?A_Brand\";v=\"8\",\"Chromium\";v=\"108\",\"GoogleChrome\";v=\"108\""},
			"sec-ch-ua-mobile":   {"?0"},
			"sec-ch-ua-platform": {"\"Windows\""},
			"sec-fetch-dest":     {"empty"},
			"sec-fetch-mode":     {"cors"},
			"sec-fetch-site":     {"same-origin"},
			"user-agent":         {"Screaming Frog SEO Spider/15.1"},
			"x-requested-with":   {"XMLHttpRequest"},
		}

		res, err := config.DefaultConfig.Network.Client.Do(req)

		if err != nil {
			config.DefaultConfig.Log.RedDelay(definederrors.MESSAGE_REQUEST_SENDING_ERROR)
			continue
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			if res.StatusCode == 403 || res.StatusCode == 401 || res.StatusCode == 400 {
				config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
				return "", RefreshSessionError{}
			}

			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			continue
		}

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			continue
		}

		var scrapePaypalLinkResponse ScrapePaypalLinkResponse

		json.Unmarshal([]byte(body), &scrapePaypalLinkResponse)

		if len(scrapePaypalLinkResponse.Redirect.Parameters) == 0 {
			return "", nil
		}

		url := fmt.Sprintf("https://www.paypal.com/checkoutnow?token=%s", scrapePaypalLinkResponse.Redirect.Parameters[0].Value)

		config.DefaultConfig.Log.Green(fmt.Sprintf("Paypal link successfully scraped! [%s]", url))

		return url, nil
	}

	return "", config.DefaultConfig.Log.LogReturnErrorCustomText(definederrors.ERROR_TOO_MANY_RETRYS, definederrors.MESSAGE_TOO_MANY_RETRYS)
}
