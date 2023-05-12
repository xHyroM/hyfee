package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

type ExperimentsKeyCache struct {
	Keys []ExperimentKey
	Last int64
}

var experimentsKeysCache = ExperimentsKeyCache{}

type ExperimentKey struct {
	Label string `json:"label"`
	Id string `json:"id"`
}

func GetExperimentKeys() []ExperimentKey {
	if time.Now().Unix()-experimentsKeysCache.Last > 180 {
		res, e := http.Get("https://api.discord-experiments.xhyrom.dev/v2/experiments?only_keys=true")
		if e != nil {
			return []ExperimentKey{}
		}

		body := []ExperimentKey{}
		json.NewDecoder(res.Body).Decode(&body)

		experimentsKeysCache = ExperimentsKeyCache{
			Keys: body,
			Last: time.Now().Unix(),
		}

		return body
	} else {
		return experimentsKeysCache.Keys
	}
}