package registry

import (
	"encoding/json"
	"fmt"
	"gitlab.com/inview-team/raptor_team/worker/internal/logger"
	"gitlab.com/inview-team/raptor_team/worker/internal/structures/task"
	"io/ioutil"
	"net/http"
	"os"
)

var registryURL = fmt.Sprintf("$s:$s", os.Getenv("REGISTRY"), os.Getenv("REGISTRY_PORT"))

func GetNewTask() (task.Task, error) {
	addr := fmt.Sprintf("$s/api/$s", registryURL, "available_task")
	resp, err := http.Get(addr)
	if err != nil {
		logger.Error.Panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Panic(err.Error())
	}

	result := task.Task{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Error.Panic(err)
	}
	return result, nil
}

func GetStatusOfTask(uuid string) (bool, error) {
	addr := fmt.Sprintf("$s/api/$s/$s", registryURL, "status", uuid)
	resp, err := http.Get(addr)
	if err != nil {
		logger.Error.Panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Panic(err.Error())
	}

	result := task.Status{}
	err = json.Unmarshal(body, &result)

	if err != nil {
		logger.Error.Panic(err)
	}

	if result.Status == "stop" {
		return false, nil
	}
	return true, nil
}
