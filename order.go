package adidasv2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	httpclient "umbrella/internal/http_client"
	retryhandler "umbrella/internal/retry_handler"
	"umbrella/internal/utils/consts"
	definederrors "umbrella/internal/utils/defined_errors"
	successhandler "umbrella/internal/utils/success_handler"
	webhookenginev2 "umbrella/internal/webhook_engine_v2"

	http "github.com/vimbing/fhttp"
)

func Order(config *Config) error {
	state := config.TaskStates.Order

	config.DefaultConfig.Log.SetState(state.Name)

	retry := retryhandler.Retry{
		Max:          state.Retry,
		Delay:        time.Millisecond * time.Duration(config.DefaultConfig.Delay),
		BypassErrors: []error{definederrors.ERROR_STOP_TASK, WaitError{}},
	}

	return retry.Retry(func() error {
		config.DefaultConfig.Log.Yellow("Creating order...")

		payload, err := json.Marshal(OrderPayload{
			BasketID:        config.Resources.BasketID,
			IsAsync:         false,
			PaymentMethodID: "PAYPAL",
			Token:           "Bearer e",
		})

		if err != nil {
			config.DefaultConfig.Log.Red(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			return err
		}

		req, err := http.NewRequest("POST", "https://www.adidas.pl/payment/hpp", bytes.NewBuffer(payload))

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
			"origin":                 {config.Region.GetOriginHeader()},
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

		config.DefaultConfig.CreateEndTimestamp()

		if res.StatusCode != 200 {
			// TODO needs handle
			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return CheckError(res, config)
		}

		body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

		if err != nil {
			return err
		}

		var orderResponse OrderResponse
		json.Unmarshal([]byte(body), &orderResponse)

		config.Resources.CheckoutID = orderResponse.CheckoutID

		if len(config.Resources.CheckoutID) == 0 {
			config.DefaultConfig.Log.Green("Error, checkout id blank!")
			return definederrors.ERROR_PLACEHOLDER
		}

		ppLink, err := ScrapePaypalLink(config)

		if err != nil || strings.Contains(ppLink, ".prod02-vm") {
			config.DefaultConfig.Log.Green("Error, while scraping pp link!")
			return err
		}

		notifierPayload := webhookenginev2.NotifierPayload{
			BotFields: webhookenginev2.BotFields{
				PrivateWebhook: config.DefaultConfig.Profile.Webhook,
				PublicWebhook:  consts.PUBLIC_WEBHOOK,
				Img:            config.Cosmetics.Image,
				Title:          "Checkout successful!",
			},
			WebhookFields: webhookenginev2.WebhookFields{
				Name:        config.Cosmetics.Name,
				Site:        config.DefaultConfig.TaskData.Website,
				Region:      config.DefaultConfig.TaskData.Region,
				Speed:       config.DefaultConfig.GetCheckoutTime(),
				Size:        config.Cosmetics.Size,
				Pid:         strings.ToUpper(config.Cosmetics.Pid),
				Price:       fmt.Sprint(config.Cosmetics.Price),
				Payment:     "PAYPAL",
				Mode:        config.DefaultConfig.TaskData.Mode,
				ProfileName: config.DefaultConfig.Profile.Name,
				TaskId:      config.DefaultConfig.TaskId,
				Email:       config.DefaultConfig.TaskData.Email,
				Password:    config.DefaultConfig.TaskData.Password,
				Proxy:       fmt.Sprintf("`%s`", config.DefaultConfig.Proxy.String),
				OrderNumber: fmt.Sprintf("[%s](%s)", "PAYPAL LINK", ppLink),
			},
		}

		successhandler.HandleSuccess(&notifierPayload)

		config.DefaultConfig.Log.Green("Order placed successfully!")

		return nil
	})
}
