package adidasv2

import (
	"encoding/json"
	"time"
	httpclient "umbrella/internal/http_client"
	retryhandler "umbrella/internal/retry_handler"
	definederrors "umbrella/internal/utils/defined_errors"

	http "github.com/vimbing/fhttp"
)

func GetBaskets(config *Config) ([]string, error) {
	state := config.TaskStates.CartClear

	config.DefaultConfig.Log.SetState(state.Name)

	retry := retryhandler.Retry{
		Max:          state.Retry,
		Delay:        time.Millisecond * time.Duration(config.DefaultConfig.Delay),
		BypassErrors: []error{definederrors.ERROR_STOP_TASK, WaitError{}},
	}

	itemsInCart := make([]string, 0)

	err := retry.Retry(func() error {
		config.DefaultConfig.Log.Yellow("Cancelling payment...")

		req, err := http.NewRequest("GET", "https://www.adidas.pl/api/chk/customer/baskets", nil)

		if err != nil {
			config.DefaultConfig.Log.Red(definederrors.MESSAGE_REQUEST_CREATING_ERROR)
			return err
		}

		req.Header = AddIstanaHeaders(http.Header{
			"authority":              {config.Region.GetAuthorityHeader()},
			"accept":                 {"*/*"},
			"accept-language":        {GetCacheBypassPayload()},
			"checkout-authorization": {"null"},
			"content-type":           {"application/json"},
			"dnt":                    {"1"},
			"referer":                {config.Region.GetRefererHeader("payment")},
			"sec-ch-ua":              {"\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\""},
			"sec-ch-ua-mobile":       {"?0"},
			"sec-ch-ua-platform":     {"\"Windows\""},
			"sec-fetch-dest":         {"empty"},
			"sec-fetch-mode":         {"cors"},
			"sec-fetch-site":         {"same-origin"},
			"user-agent":             {USER_AGENT},
		})

		res, err := config.DefaultConfig.Network.Client.Do(req)

		if err != nil {
			config.DefaultConfig.Log.Red(definederrors.MESSAGE_REQUEST_SENDING_ERROR)
			return err
		}

		defer res.Body.Close()

		if res.StatusCode != 200 && res.StatusCode != 404 {
			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return definederrors.ERROR_PLACEHOLDER
		}

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			return err
		}

		var getBasketsResponse GetBasketsResponse
		json.Unmarshal([]byte(body), &getBasketsResponse)

		config.Resources.BasketID = getBasketsResponse.BasketID

		config.DefaultConfig.Log.Green("Basket scraped successfully!")

		if len(getBasketsResponse.ShipmentList) == 0 {
			return nil
		}

		for _, sList := range getBasketsResponse.ShipmentList {
			if len(sList.ProductLineItemList) == 0 {
				continue
			}

			for _, item := range sList.ProductLineItemList {
				itemsInCart = append(itemsInCart, item.ItemID)
			}
		}

		return nil
	})

	return itemsInCart, err
}
