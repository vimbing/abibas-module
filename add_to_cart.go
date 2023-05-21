package adidasv2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	httpclient "umbrella/internal/http_client"
	retryhandler "umbrella/internal/retry_handler"
	definederrors "umbrella/internal/utils/defined_errors"

	http "github.com/vimbing/fhttp"
)

func AddToCart(config *Config, variant *Variant) error {
	state := config.TaskStates.AddToCart

	config.DefaultConfig.Log.SetState(state.Name)

	retry := retryhandler.Retry{
		Max:          state.Retry,
		Delay:        time.Millisecond * time.Duration(config.DefaultConfig.Delay),
		BypassErrors: []error{definederrors.ERROR_STOP_TASK, WaitError{}},
	}

	return retry.Retry(func() error {
		payload, err := json.Marshal(AddToCartPayload{
			{
				ProductID:            strings.ToUpper(variant.Pid),
				Quantity:             config.DefaultConfig.TaskData.Quantity.GetQuantity(),
				ProductVariationSku:  strings.ToUpper(variant.SizePid),
				ProductID0:           strings.ToUpper(variant.SizePid),
				Size:                 variant.SizeValue,
				DisplaySize:          variant.SizeValue,
				SpecialLaunchProduct: true,
			},
		})

		config.DefaultConfig.Log.Yellow("Adding to cart...")

		if err != nil {
			config.DefaultConfig.Log.Red(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			return err
		}

		req, err := http.NewRequest("POST", "https://www.adidas.pl/api/chk/baskets/-/items?orderType=multiShipping", bytes.NewBuffer(payload))

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
			"glassversion":           {"38bbcaa"},
			"origin":                 {config.Region.GetOriginHeader()},
			"referer":                {config.Region.GetRefererHeader(fmt.Sprintf("-/%s.html", strings.ToUpper(variant.Pid)))},
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

		if res.StatusCode != 200 {
			err := CheckError(res, config)

			if err != nil {
				if res.StatusCode == 400 {
					return WaitError{}
				}

				config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
				return err
			}

			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return definederrors.ERROR_PLACEHOLDER
		}

		err = CheckError(res, config)

		if err != nil {
			return err
		}

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			return err
		}

		var addToCartResponse AddToCartResponse
		json.Unmarshal([]byte(body), &addToCartResponse)

		config.Resources.BasketID = addToCartResponse.BasketID
		config.Cosmetics = GetCosmetics(config, &addToCartResponse)

		config.DefaultConfig.Log.Green("Item successfully added to cart!")

		return nil
	})
}
