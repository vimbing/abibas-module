package adidasv2

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	httpclient "umbrella/internal/http_client"
	retryhandler "umbrella/internal/retry_handler"
	definederrors "umbrella/internal/utils/defined_errors"
	"umbrella/internal/utils/helpers"

	http "github.com/vimbing/fhttp"
)

var (
	latestType = ""
)

func GetItemData(config *Config) ([]Variant, error) {
	latestType = ITEM_DATA_TYPE_FRONTEND

	state := config.TaskStates.Address

	variants := make([]Variant, 0)

	retry := retryhandler.Retry{
		Max:          state.Retry,
		Delay:        time.Millisecond * time.Duration(config.DefaultConfig.Delay),
		BypassErrors: []error{definederrors.ERROR_STOP_TASK, WaitError{}},
	}

	err := retry.Retry(func() error {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://www.adidas.pl/api/products/%s/availability", config.DefaultConfig.TaskData.Sku), nil)

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
			if res.StatusCode == 403 {
				config.DefaultConfig.Log.YellowDelay("Item blocked, monitoring...")
				return nil
			}

			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return definederrors.ERROR_PLACEHOLDER
		}

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			return err
		}

		var getItemDataResponse GetItemDataResponse
		json.Unmarshal([]byte(body), &getItemDataResponse)

		if len(getItemDataResponse.VariationList) == 0 {
			return nil
		}

		for _, v := range getItemDataResponse.VariationList {
			if v.Availability > 0 {
				variants = append(variants, Variant{
					Pid:       getItemDataResponse.ID,
					SizePid:   v.Sku,
					SizeValue: v.Size,
				})
			}
		}

		return nil
	})

	return variants, err
}

func GetItemDataBackend(config *Config) ([]Variant, error) {
	latestType = ITEM_DATA_TYPE_BACKEND

	state := config.TaskStates.Monitor

	variants := make([]Variant, 0)

	retry := retryhandler.Retry{
		Max:          state.Retry,
		Delay:        time.Millisecond * time.Duration(config.DefaultConfig.Delay),
		BypassErrors: []error{definederrors.ERROR_STOP_TASK, WaitError{}},
	}

	err := retry.Retry(func() error {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.3stripes.net/gw-api/v2/products/%s/Availability/", strings.ToUpper(config.DefaultConfig.TaskData.Sku)), nil)

		if err != nil {
			config.DefaultConfig.Log.Red(definederrors.MESSAGE_REQUEST_CREATING_ERROR)
			return err
		}

		req.Header = http.Header{
			"host":            {"api.3stripes.net"},
			"user-agent":      {fmt.Sprintf("adidas/2022.11.25.13.12CFNetwork/1390Darwin/22.0.0%s", helpers.RandomString(helpers.RandomInt(8, 25), true, false, true))},
			"x-market":        {"PL"},
			"accept-language": {GetCacheBypassPayload()},
			"accept":          {"application/hal+json"},
			"content-type":    {"application/json;charset=UTF-8"},
			"Cache-Control":   {"no-cache:max-age=0"},
			"Pragma":          {"no-cache"},
			"x-api-key":       {"m79qyapn2kbucuv96ednvh22"},
		}

		res, err := config.DefaultConfig.Network.Client.Do(req)

		if err != nil {
			config.DefaultConfig.Log.Red(definederrors.MESSAGE_REQUEST_SENDING_ERROR)
			return err
		}

		defer res.Body.Close()

		if res.StatusCode != 200 && res.StatusCode != 404 {
			if res.StatusCode == 403 {
				config.DefaultConfig.Log.YellowDelay("Item blocked, monitoring...")
				return nil
			}

			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return definederrors.ERROR_PLACEHOLDER
		}

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			return err
		}

		var getItemDataResponse GetItemDataBackendResponse
		json.Unmarshal([]byte(body), &getItemDataResponse)

		if len(getItemDataResponse.Embedded.Variations) < 1 {
			config.DefaultConfig.Log.YellowDelay("No variants for item, maybe not loaded...")
			return nil
		}

		for _, v := range getItemDataResponse.Embedded.Variations {
			if v.Orderable {
				variants = append(variants, Variant{
					SizePid:   v.VariationProductID,
					SizeValue: v.Size,
					Pid:       getItemDataResponse.ProductID,
				})
			}
		}

		return nil
	})

	return variants, err
}

func GetMulitiData(config *Config) ([]Variant, error) {
	if latestType == ITEM_DATA_TYPE_FRONTEND {
		return GetItemDataBackend(config)
	}

	return GetItemData(config)
}
