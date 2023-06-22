package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
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

func GetExperimentKeys(query string) []ExperimentKey {
	if time.Now().Unix()-experimentsKeysCache.Last > 180 {
		res, e := http.Get("https://api.distools.xhyrom.dev/v2/experiments?only_keys=true&also_with_unknown_ids=true"+query)
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

type Experiment struct {
	Data ExperimentData `json:"data"`
	Rollout ExperimentRollout `json:"rollout"`
}

type ExperimentData struct {
	Kind string `json:"kind"`
	Id string `json:"id"`
	Label string `json:"label"`
	Description []string `json:"description"`
	Hash int `json:"hash"`
	Buckets []int `json:"buckets"`
	ConfigKeys []string `json:"config_keys"`
}

type ExperimentRollout struct {
	Populations []ExperimentRolloutPopulation `json:"populations"`
	Revision int `json:"revision"`
	Overrides map[string][]string `json:"overrides"`
	OverridesFormatted []ExperimentRolloutPopulation `json:"overrides_formatted"`
}

type ExperimentRolloutPopulation struct {
	Buckets map[string]ExperimentRolloutPopulationBucket `json:"buckets"`
	Filters []ExperimentRolloutPopulationFilter `json:"filters"`
}

type ExperimentRolloutPopulationBucket struct {
	Rollout []ExperimentRolloutPopulationBucketRollout `json:"rollout"`
}

type ExperimentRolloutPopulationBucketRollout struct {
	Start int `json:"start"`
	End int `json:"end"`
}

type ExperimentRolloutPopulationFilter struct {
	Type string `json:"type"`
	Features []string `json:"features"`
	MinId int `json:"min_id"`
	MaxId int `json:"max_id"`
	MinCount int `json:"min_count"`
	MaxCount int `json:"max_count"`
	Ids []string `json:"ids"`
	HubTypes []int `json:"hub_types"`
	HasVanity bool `json:"has_vanity"`
	HashKey int `json:"hash_key"`
	Target int `json:"target"`
}

type Eligible struct {
	Eligible bool `json:"eligible"`
	Bucket EligibleBucket `json:"bucket"`
	Filters []ExperimentRolloutPopulationFilter `json:"filters"`
}

type EligibleBucket struct {
	ExperimentRolloutPopulationBucket
	Id string `json:"id"`
}

type ErrorResponse struct {
	Status int `json:"status"`
	Message string `json:"message"`
}

func GetExperiment(name string) (Experiment, error) {
	res, e := http.Get("https://api.distools.xhyrom.dev/v2/experiments/" + name)
	if e != nil {
		return Experiment{}, errors.New("Error getting experiment")
	}

	if res.StatusCode != 200 {
		body := ErrorResponse{}
		json.NewDecoder(res.Body).Decode(&body)

		return Experiment{}, errors.New(body.Message)
	}

	body := Experiment{}
	json.NewDecoder(res.Body).Decode(&body)

	return body, nil
}

func IsExperimentEligible(id string, guild discord.Guild) (Eligible, error) {
	payload, _ := json.Marshal(map[string]interface{}{
		"experiment_id": id,
		"guild": guild,
	})

	res, e := http.Post("https://api.distools.xhyrom.dev/v2/eligible", "application/json", bytes.NewBuffer(payload))
	if e != nil {
		return Eligible{}, errors.New("Error checking eligibility")
	}

	if res.StatusCode != 200 {
		body := ErrorResponse{}
		json.NewDecoder(res.Body).Decode(&body)

		return Eligible{}, errors.New(body.Message)
	}

	body := Eligible{}
	json.NewDecoder(res.Body).Decode(&body)

	return body, nil
}

func (experiment Experiment) FormatName() string {
	if experiment.Data.Label != "" {
		return experiment.Data.Label + " (" + experiment.Data.Id + ")"
	}

	if experiment.Data.Id != "" {
		return experiment.Data.Id
	}

	return "Unknown"
}

func (experiment Experiment) FormatDescription() string {
	formatted := ""

	for _, desc := range experiment.Data.Description {
		if desc == "Control" {
			continue
		}

		split := strings.Split(desc, ":")
		name, description := split[0], split[1]
		
		formatted += "**" + name + "**: " + description + "\n"
	}

	return formatted
}

func (experiment Experiment) FormatOverrides() string {
	overrides := []string{}

	for key, value := range experiment.Rollout.Overrides {
		var name string
		if key == "none" {
			name = "None"
		} else {
			name = "Treatment " + key
		}

		overrides = append(overrides, "**" + name + "**: " + strings.Join(value, ", "))
	}

	return strings.Join(overrides, "\n\n")
}

func (experiment Experiment) FormatPopulations() string {
	formatted := ""

	for _, population := range experiment.Rollout.Populations {
		formatted += population.Format()
		formatted += "\n"
	}

	return formatted
}

func (experiment Experiment) FormatOverridesFormatted() string {
	formatted := ""

	for _, population := range experiment.Rollout.OverridesFormatted {
		formatted += population.Format()
		formatted += "\n"
	}

	return formatted
}

func (population ExperimentRolloutPopulation) Format() string {
	formatted := ""

	filters := []string{}

	for _, filter := range population.Filters {
		filters = append(filters, filter.Format())
	}

	if len(filters) > 0 {
		formatted += "**Filters**: " + strings.Join(filters, " and ") + "\n"
	}

	for id, bucket := range population.Buckets {
		formatted += bucket.Format(id)
	}

	return formatted
}

func (bucket ExperimentRolloutPopulationBucket) Format(id string) string {
	formatted := ""

	percentage := 0
	for _, rollout := range bucket.Rollout {
		percentage += rollout.End - rollout.Start
	}

	percentage = percentage / 100

	rollouts := []string{}
	for _, rollout := range bucket.Rollout {
		rollouts = append(rollouts, strconv.Itoa(rollout.Start) + "-" + strconv.Itoa(rollout.End))
	}

	var name string
	if id == "none" {
		name = "None"
	} else {
		name = "Treatment " + id
	}

	formatted += "**" + name + "**: " + strconv.Itoa(percentage) + "% (" + strings.Join(rollouts, ", ") + ")\n"

	return formatted
}

func (filter ExperimentRolloutPopulationFilter) Format() string {
	switch filter.Type {
		case "guild_has_feature": {
			features := []string{}
			for _, feature := range filter.Features {
				features = append(features, SnakeCaseToPascalCaseWithSpaces(feature))
			}

			return "Server has feature " + strings.Join(features, " or ")
		}
		case "guild_id_range": {
			return "Server Id is in range " + strconv.Itoa(filter.MinId) + " - " + strconv.Itoa(filter.MaxId)
		}
		case "guild_member_count_range": {
			if filter.MaxCount != 0 {
				return "Server member count is in range " + strconv.Itoa(filter.MinCount) + " - " + strconv.Itoa(filter.MaxCount)
			} else {
				return "Server member count is " + strconv.Itoa(filter.MinCount) + "+"
			}
		}
		case "guild_ids": {
			return "Server Id is " + strings.Join(filter.Ids, " or ")
		}
		case "guild_hub_types": {
			hub_types := []string{}
			for _, hub_type := range filter.HubTypes {
				switch hub_type {
					case 0:
						hub_types = append(hub_types, "Default")
					case 1:
						hub_types = append(hub_types, "High School")
					case 2:
						hub_types = append(hub_types, "College")
				}
			}
			return "Server hub type is " + strings.Join(hub_types, " or ")
		}
		case "guild_has_vanity_url": {
			if filter.HasVanity {
				return "Server has a vanity url"
			}

			return "Server does not have a vanity url"
		}
		case "guild_in_range_by_hash": {
			return strconv.Itoa(filter.Target / 100) + "% of servers (hash key " + strconv.Itoa(filter.HashKey) + ", target " + strconv.Itoa(filter.Target) + ")"
		}
	}

	return filter.Type
}