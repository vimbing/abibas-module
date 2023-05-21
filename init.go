package adidasv2

import (
	profilesreader "umbrella/internal/file_readers/profiles_reader"
	"umbrella/internal/file_readers/proxy"
	tasksreader "umbrella/internal/file_readers/tasks_reader"
	"umbrella/internal/modules"
	retrymenager "umbrella/internal/retry_menager"
	globaltypes "umbrella/internal/utils/global_types"

	"github.com/vimbing/cclient"
	"github.com/vimbing/fhttp/cookiejar"
	tls "github.com/vimbing/utls"
)

func Init(proxy *proxy.Proxy, profile *profilesreader.Profile, taskData *tasksreader.TaskData, id int) (Config, error) {
	defaultConfig := modules.Init()

	err := defaultConfig.SetDefaultConfig(proxy, profile, taskData, id)

	if err != nil {
		return Config{}, err
	}

	taskStates := TaskStates{
		Login: globaltypes.TaskState{
			Name:  "LOGIN",
			Retry: retrymenager.GetStateRetry(taskData.Website, "Session"),
		},
		CartClear: globaltypes.TaskState{
			Name:  "CART CLEANING",
			Retry: retrymenager.GetStateRetry(taskData.Website, "AddToCart"),
		},
		AddToCart: globaltypes.TaskState{
			Name:  "ADD TO CART",
			Retry: retrymenager.GetStateRetry(taskData.Website, "AddToCart"),
		},
		Address: globaltypes.TaskState{
			Name:  "ADDRESS",
			Retry: retrymenager.GetStateRetry(taskData.Website, "AddToCart"),
		},
		Monitor: globaltypes.TaskState{
			Name:  "MONITOR",
			Retry: retrymenager.GetStateRetry(taskData.Website, "AddToCart"),
		},
		Order: globaltypes.TaskState{
			Name:  "ORDER",
			Retry: retrymenager.GetStateRetry(taskData.Website, "AddToCart"),
		},
		Akamai: globaltypes.TaskState{
			Name:  "AKAMAI",
			Retry: retrymenager.GetStateRetry(taskData.Website, "AddToCart"),
		},
	}

	config := Config{
		DefaultConfig: defaultConfig,
		TaskStates:    taskStates,
		Resources:     Resources{},
	}

	config.Region = &Region{
		RegionCode: "pl",
	}

	config.DefaultConfig.Network.Client, _ = cclient.NewClient(tls.HelloChrome_100, "", true, 6)
	config.DefaultConfig.Network.Client.Jar, _ = cookiejar.New(nil)

	return config, nil
}
