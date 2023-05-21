package adidasv2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"umbrella/internal/antibots/captcha/captchaai"
	settingsreader "umbrella/internal/file_readers/settings_reader"
	httpclient "umbrella/internal/http_client"
	retryhandler "umbrella/internal/retry_handler"
	definederrors "umbrella/internal/utils/defined_errors"

	http "github.com/vimbing/fhttp"
)

func solveCaptcha(config *Config) (string, error) {
	if len(settingsreader.GetCaptchaAiKey()) < 5 {
		config.DefaultConfig.Log.Yellow("You need to have captchaai key set!")
		return "", definederrors.ERROR_STOP_TASK
	}

	captchaTask := captchaai.RecaptchaV3Task{
		Config: captchaai.Config{
			ApiKey: settingsreader.GetCaptchaAiKey(),
		},
		RecaptchaV3Payload: captchaai.RecaptchaV3Payload{
			WebsiteURL: "https://www.adidas.pl/api/account/login",
			Sitekey:    "6LdQquAaAAAAALeU6cp88M5ByhWDANC1-ei8xfMW",
			PageAction: "login",
			Proxy:      config.DefaultConfig.Proxy.String,
		},
	}

	config.DefaultConfig.Log.Yellow("Solving captcha...")

	capChan, err := captchaai.Solve(captchaai.Config{ApiKey: settingsreader.GetCaptchaAiKey()}, captchaTask)

	if err != nil {
		return "", err
	}

	token := <-capChan

	if len(token) < 5 {
		config.DefaultConfig.Log.Yellow("Captcha solving error, check your captchaai key and balance!")
		return "", definederrors.ERROR_PLACEHOLDER
	}

	return token, nil
}

func Login(config *Config) error {
	state := config.TaskStates.Login

	config.DefaultConfig.Log.SetState(state.Name)

	retry := retryhandler.Retry{
		Max:          state.Retry,
		Delay:        time.Millisecond * time.Duration(config.DefaultConfig.Delay),
		BypassErrors: []error{definederrors.ERROR_STOP_TASK, WaitError{}},
	}

	return retry.Retry(func() error {
		token, err := solveCaptcha(config)

		if err != nil {
			return err
		}

		payload, err := json.Marshal(LoginPayload{
			User:            config.DefaultConfig.TaskData.Email,
			Password:        config.DefaultConfig.TaskData.Password,
			PersistentLogin: true,
			Recaptcha:       token,
		})

		config.DefaultConfig.Log.Yellow("Logging in...")

		if err != nil {
			config.DefaultConfig.Log.Red(definederrors.MESSAGE_JSON_MARSHALING_ERROR)
			return err
		}

		req, err := http.NewRequest("POST", "https://www.adidas.pl/api/account/login", bytes.NewBuffer(payload))

		if err != nil {
			config.DefaultConfig.Log.Red(definederrors.MESSAGE_REQUEST_CREATING_ERROR)
			return err
		}

		req.Header = AddIstanaHeaders(http.Header{
			"authority":          {config.Region.GetAuthorityHeader()},
			"accept":             {"*/*"},
			"accept-language":    {GetCacheBypassPayload()},
			"content-type":       {"application/json"},
			"dnt":                {"1"},
			"origin":             {config.Region.GetOriginHeader()},
			"referer":            {config.Region.GetRefererHeader("account-login")},
			"sec-ch-ua":          {"\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\""},
			"sec-ch-ua-mobile":   {"?0"},
			"sec-ch-ua-platform": {"\"Windows\""},
			"sec-fetch-dest":     {"empty"},
			"sec-fetch-mode":     {"cors"},
			"sec-fetch-site":     {"same-origin"},
			"user-agent":         {USER_AGENT},
			"x-client-tag":       {"AccountPortal"},
		})

		res, err := config.DefaultConfig.Network.Client.Do(req)

		if err != nil {
			config.DefaultConfig.Log.Red(definederrors.MESSAGE_REQUEST_SENDING_ERROR)
			return err
		}

		b, _ := httpclient.GetBodyString(res, &config.DefaultConfig.Log)
		fmt.Println(b)

		defer res.Body.Close()

		config.Resources.SessionTimeout = time.Now().Add(time.Hour * 2)

		if res.StatusCode != 200 {
			if res.StatusCode == 400 {
				config.DefaultConfig.Log.RedDelay("Error while logging in, probably wrong credentials! [stoping task]")
				return definederrors.ERROR_STOP_TASK
			}

			config.DefaultConfig.Log.StatusCodeErrorDelay(res.Status)
			return definederrors.ERROR_PLACEHOLDER
		}

		config.DefaultConfig.Log.Green("Succesfully logged in!")

		return nil
	})
}
