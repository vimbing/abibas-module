package adidasv2

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"umbrella/internal/utils/helpers"
)

func Monitor(config *Config) (Variant, []Variant, error) {
	state := config.TaskStates.Monitor

	config.DefaultConfig.Log.SetState(state.Name)

	for {
		select {
		case cData := <-config.DefaultConfig.TaskData.MonitorModeUtils.VariantChannel:
			var variants []Variant
			cData.UnmarshalData(&variants)

			if len(variants) == 0 {
				return Variant{}, []Variant{}, errors.New("placeholder")
			}

			selectedVariant := variants[helpers.RandomInt(0, len(variants)-1)]

			config.DefaultConfig.Log.Green(fmt.Sprintf("Found info from monitor mode! [%d Variants][Selected: %s]", len(variants), selectedVariant.SizeValue))

			return selectedVariant, variants, nil
		default:
			if time.Now().After(config.Resources.SessionTimeout) && !strings.EqualFold(config.DefaultConfig.TaskData.Mode, "monitor") {
				config.DefaultConfig.Log.Yellow("Session needs refresh...")
				return Variant{}, []Variant{}, RefreshSessionError{}
			}

			helpers.Sleep(config.DefaultConfig.Delay)

			currentItemData, err := GetMulitiData(config)

			if err != nil {
				continue
			}

			if len(currentItemData) > 0 {
				if strings.ToLower(config.DefaultConfig.TaskData.Size) == "random" {
					config.DefaultConfig.Log.Green("Item is instock, proceeding to checkout!")
					return currentItemData[helpers.RandomInt(0, len(currentItemData)-1)], currentItemData, nil
				}

				var availableSizes []Variant

				for _, variant := range currentItemData {
					if helpers.SliceContainsString(config.DefaultConfig.TaskData.SizeArray, variant.SizeValue) {
						availableSizes = append(availableSizes, variant)
					}
				}

				if len(availableSizes) > 0 {
					config.DefaultConfig.Log.Green("Item instock, picking random from provided sizes!")
					return availableSizes[helpers.RandomInt(0, len(availableSizes)-1)], availableSizes, nil
				}
			}

			config.DefaultConfig.Log.Yellow("Waiting for restock...")
		}
	}
}
