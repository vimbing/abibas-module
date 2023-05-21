package adidasv2

import (
	"fmt"
	"time"
	retryhandler "umbrella/internal/retry_handler"
	definederrors "umbrella/internal/utils/defined_errors"

	http "github.com/vimbing/fhttp"
)

func DeleteFromCart(config *Config, itemId string) error {
	state := config.TaskStates.CartClear

	config.DefaultConfig.Log.SetState(state.Name)

	retry := retryhandler.Retry{
		Max:          state.Retry,
		Delay:        time.Millisecond * time.Duration(config.DefaultConfig.Delay),
		BypassErrors: []error{definederrors.ERROR_STOP_TASK, WaitError{}},
	}

	return retry.Retry(func() error {
		config.DefaultConfig.Log.Yellow("Removing item from cart...")

		req, err := http.NewRequest("DELETE", fmt.Sprintf("https://www.adidas.pl/api/chk/baskets/%s/items/%s", config.Resources.BasketID, itemId), nil)

		if err != nil {
			config.DefaultConfig.Log.Red(definederrors.MESSAGE_REQUEST_CREATING_ERROR)
			return err
		}

		req.Header = AddIstanaHeaders(http.Header{
			"authority":              {config.Region.GetAuthorityHeader()},
			"accept":                 {"*/*"},
			"accept-language":        {GetCacheBypassPayload()},
			"checkout-authorization": {"Bearer e"},
			"content-type":           {"application/json"},
			"dnt":                    {"1"},
			"glassversion":           {"980f8aa"},
			"origin":                 {config.Region.GetOriginHeader()},
			"referer":                {config.Region.GetRefererHeader("cart")},
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

		if res.StatusCode != 200 && res.StatusCode != 204 {
			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return CheckError(res, config)
		}

		config.DefaultConfig.Log.Green("Item removed successfully!")

		return nil
	})
}
