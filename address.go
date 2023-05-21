package adidasv2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	retryhandler "umbrella/internal/retry_handler"
	definederrors "umbrella/internal/utils/defined_errors"
	"umbrella/internal/utils/helpers"

	http "github.com/vimbing/fhttp"
)

func Address(config *Config) error {
	state := config.TaskStates.Address

	config.DefaultConfig.Log.SetState(state.Name)

	retry := retryhandler.Retry{
		Max:          state.Retry,
		Delay:        time.Millisecond * time.Duration(config.DefaultConfig.Delay),
		BypassErrors: []error{definederrors.ERROR_STOP_TASK, WaitError{}},
	}

	return retry.Retry(func() error {
		config.DefaultConfig.Log.Yellow("Filling address...")

		var vatCode string

		if helpers.ParseBoolean(config.DefaultConfig.Profile.Invoice) {
			vatCode = config.DefaultConfig.Profile.Nip
		} else {
			vatCode = ""
		}

		payload, err := json.Marshal(AddressPayload{
			Customer: Customer{
				Email:             config.DefaultConfig.TaskData.Email,
				ReceiveSmsUpdates: false,
			},
			ShippingAddress: ShippingAddress{
				EmailAddress: config.DefaultConfig.TaskData.Email,
				Country:      strings.ToUpper(config.Region.RegionCode),
				FirstName:    config.DefaultConfig.Profile.FirstName,
				LastName:     config.DefaultConfig.Profile.LastName,
				Address1:     config.DefaultConfig.Profile.Street,
				Zipcode:      config.DefaultConfig.Profile.PostalCode,
				City:         config.DefaultConfig.Profile.City,
				PhoneNumber:  config.DefaultConfig.Profile.Phone,
				Address2:     "",
			},
			BillingAddress: BillingAddress{
				EmailAddress: config.DefaultConfig.TaskData.Email,
				Country:      strings.ToUpper(config.Region.RegionCode),
				FirstName:    config.DefaultConfig.Profile.FirstName,
				LastName:     config.DefaultConfig.Profile.LastName,
				Address1:     config.DefaultConfig.Profile.Street,
				Zipcode:      config.DefaultConfig.Profile.PostalCode,
				City:         config.DefaultConfig.Profile.City,
				PhoneNumber:  config.DefaultConfig.Profile.Phone,
				Address2:     "",
				VatCode:      vatCode,
			},
			LegalAcceptance: LegalAcceptance{
				ScenarioName:     "guest-checkout:adidas:PL",
				AcceptanceLocale: fmt.Sprintf("%s-%s", strings.ToLower(config.Region.RegionCode), strings.ToUpper(config.Region.RegionCode)),
				Acceptances: []Acceptances{
					{
						DocumentName: "doc-mrkt-email-checkout:adidas:PL:2023227",
						Accepted:     true,
					},
					{
						DocumentName: "doc-mrkt-paidmedia:adidas:PL:202326",
						Accepted:     true,
					},
				},
			},
		})

		if err != nil {
			config.DefaultConfig.Log.Red(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			return err
		}

		req, err := http.NewRequest("PATCH", fmt.Sprintf("https://www.adidas.pl/api/chk/baskets/%s", config.Resources.BasketID), bytes.NewBuffer(payload))

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
			"referer":                {config.Region.GetRefererHeader("delivery")},
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
			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return CheckError(res, config)
		}

		config.DefaultConfig.Log.Green("Address filled successfully!")

		return nil
	})
}
