package adidasv2

import (
	"fmt"
	"time"
	retryhandler "umbrella/internal/retry_handler"
	definederrors "umbrella/internal/utils/defined_errors"
)

func CartCleaner(config *Config) error {
	state := config.TaskStates.CartClear

	config.DefaultConfig.Log.SetState(state.Name)

	retry := retryhandler.Retry{
		Max:          state.Retry,
		Delay:        time.Millisecond * time.Duration(config.DefaultConfig.Delay),
		BypassErrors: []error{definederrors.ERROR_STOP_TASK, WaitError{}},
	}

	return retry.Retry(func() error {
		items, err := GetBaskets(config)

		if err != nil {
			return err
		}

		err = CancelPayment(config)

		if err != nil {
			return err
		}

		if len(items) == 0 {
			return nil
		}

		config.DefaultConfig.Log.Yellow(fmt.Sprintf("Found %d items in cart, removing...", len(items)))

		for _, item := range items {
			err = DeleteFromCart(config, item)

			if err != nil {
				return err
			}
		}

		return nil
	})
}
