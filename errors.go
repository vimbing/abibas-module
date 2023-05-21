package adidasv2

import (
	"strings"
	httpclient "umbrella/internal/http_client"
	definederrors "umbrella/internal/utils/defined_errors"

	http "github.com/vimbing/fhttp"
)

type OpCode string

const (
	OpCode_PaymentInProgressError OpCode = "PaymentInProgressBasket"
)

type WaitError struct {
}

func (w WaitError) Error() string {
	return "task_wait_error"
}

type LongWaitError struct {
}

func (w LongWaitError) Error() string {
	return "task_long_wait_error"
}

type RefreshSessionError struct {
}

func (r RefreshSessionError) Error() string {
	return "task_refresh"
}

type RelogError struct {
}

func (r RelogError) Error() string {
	return "relog_needed"
}

func CheckError(res *http.Response, config *Config) error {
	body, err := httpclient.GetBodyString(res, &config.DefaultConfig.Log)

	if err != nil {
		return err
	}

	config.DefaultConfig.Log.TextLog(body, res.StatusCode)

	switch b, status := body, res.StatusCode; {
	case status == 200 || status == 201:
		return nil
	case strings.Contains(b, "PaymentInProgressBasket"):
		config.DefaultConfig.Log.RedDelay("Error, there is open payment on account!")

		err = CancelPayment(config)

		if err != nil {
			return err
		}

		return definederrors.ERROR_PLACEHOLDER
	case status == 403:
		config.DefaultConfig.Log.RedDelay("Error, probably proxy blocked! [stoping task]")
		return definederrors.ERROR_STOP_TASK
	case strings.Contains(b, "ProductItemNotAvailableException"):
		config.DefaultConfig.Log.RedDelay("Error while adding to cart [ProductItemNotAvailableException]")
		return LongWaitError{}
	default:
		return definederrors.ERROR_PLACEHOLDER
	}
}
