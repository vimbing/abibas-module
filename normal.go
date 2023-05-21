package adidasv2

import (
	"strings"
	profilesreader "umbrella/internal/file_readers/profiles_reader"
	"umbrella/internal/file_readers/proxy"
	tasksreader "umbrella/internal/file_readers/tasks_reader"
	quicktaskshandler "umbrella/internal/quicktasks_handler"
	clititle "umbrella/internal/utils/cli_title"
	definederrors "umbrella/internal/utils/defined_errors"
	waithandler "umbrella/internal/utils/wait_handler"
)

func Normal(proxy *proxy.Proxy, profile *profilesreader.Profile, task *tasksreader.TaskData, id int) {
	defer clititle.DecreaseRunning()

	config, err := Init(proxy, profile, task, id)

	if err != nil {
		return
	}
	go quicktaskshandler.RegisterTaskToQuicktaskHandler(config.DefaultConfig.TaskData)

	for {
		if err != nil {
			if err.Error() == definederrors.IDENTIFIER_STOP_TASK {
				break
			}

			HandleModeRetryError(err, &config)
		}

		SolveAkamai(&config, false)

		err = Login(&config)

		if err != nil {
			continue
		}

		SolveAkamai(&config, true)

		err = CartCleaner(&config)

		if err != nil {
			continue
		}

		err = waithandler.HandleUserWait(&config.DefaultConfig, config.Resources.SessionTimeout)

		if err != nil {
			continue
		}

		variant, _, err := Monitor(&config)

		if err != nil {
			continue
		}

		config.DefaultConfig.CreateStartTimestamp()

		SolveAkamai(&config, true)

		err = AddToCart(&config, &variant)

		if err != nil {
			continue
		}

		err = Address(&config)

		if err != nil {
			continue
		}

		SolveAkamai(&config, true)

		err = Order(&config)

		if err != nil {
			continue
		}

		if strings.ToLower(task.ResetAfterSuccess) != "true" {
			return
		}
	}
}
