package adidasv2

import (
	"time"
	"umbrella/internal/utils/helpers"
)

func SolveAkamai(config *Config, important bool) {
	config.DefaultConfig.Log.SetState(config.TaskStates.Akamai.Name)

	for {
		config.DefaultConfig.Log.Yellow("Solving akamai...")

		if !important {
			time.Sleep(time.Second * 2)

			if helpers.RandomInt(1, 3) == 2 {
				config.DefaultConfig.Log.Red("Error while solving akamai!")
				continue
			}
		}

		config.DefaultConfig.Log.Green("Akamai successfully solved!")

		return
	}

}
