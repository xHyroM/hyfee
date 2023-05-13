package commands

import (
	"hyros_coffee/handler"
	"hyros_coffee/hyfee"
	"hyros_coffee/utils"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func Experiments(bot *hyfee.Bot) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name: "experiments",
			Description: "Discord experiments",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name: "get",
					Description: "Get information about an experiment",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name: "experiment",
							Description: "The experiment id",
							Required: true,
							Autocomplete: true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name: "eligible",
					Description: "Check if a guild is eligible for an experiment",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name: "experiment",
							Description: "The experiment id",
							Required: true,
							Autocomplete: true,
						},
						discord.ApplicationCommandOptionString{
							Name: "guild",
							Description: "The guild id",
							Required: true,
						},
					},
				},
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"get": getExperimentsHandler(bot),
			"eligible": eligibleExperimentsHandler(bot),
		},
		AutocompleteHandlers: map[string]handler.AutocompleteHandler{
			"get": autocompleteHandler(bot),
			"eligible": autocompleteHandler(bot),
		},
	}
}

func getExperimentsHandler(bot *hyfee.Bot) handler.CommandHandler {
	return func(event *events.ApplicationCommandInteractionCreate) error {
		id := event.SlashCommandInteractionData().String("experiment")

		experiment, err := utils.GetExperiment(id)
		if err != nil {
			return event.CreateMessage(discord.MessageCreate{
				Content: "Error: " + err.Error(),
				Flags: discord.MessageFlagEphemeral,
			})
		}

		embed := discord.NewEmbedBuilder().
			SetTitle("Experiment Details").
			AddField("Name", experiment.FormatName(), true).
			AddField("Kind", experiment.Data.Kind, true).
			AddField("Hash", strconv.Itoa(experiment.Data.Hash), true)

		if len(experiment.Data.Description) > 0 {
			embed.AddField("Treatments", experiment.FormatDescription(), false)
		}

		if len(experiment.Rollout.Populations) > 0 {
			embed.AddField("Populations", experiment.FormatPopulation(), false)
		}

		if len(experiment.Rollout.Overrides) > 0 {
			embed.AddField("Overrides", experiment.FormatOverrides(), false)
		}

		return event.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{embed.Build()},
		})
	}
}

func eligibleExperimentsHandler(bot *hyfee.Bot) handler.CommandHandler {
	return func(e *events.ApplicationCommandInteractionCreate) error {
		id := e.SlashCommandInteractionData().String("id")
		guild := e.SlashCommandInteractionData().String("guild")

		return e.CreateMessage(discord.MessageCreate{
			Content: "You said: " + id + " " + guild,
			Flags: discord.MessageFlagEphemeral,
		})
	}
}

func autocompleteHandler(bot *hyfee.Bot) handler.AutocompleteHandler {
	return func(e *events.AutocompleteInteractionCreate) error {
		option := e.Data.String("experiment")

		experiments := utils.GetExperimentKeys()

		result := []discord.AutocompleteChoice{}

		for _, experiment := range experiments {
			if option == "" || strings.Contains(strings.ToLower(experiment.Label), strings.ToLower(option)) {
				result = append(result, discord.AutocompleteChoiceString{
					Name: utils.IfThenElse(
						len(experiment.Label) > 25,
						func() string {
							return experiment.Label[:23]  + ".."
						},
						func() string {
							return experiment.Label
						}),
					Value: experiment.Id,
				})
			}

			if len(result) == 25 {
				break
			}
		}

		return e.Result(result)
	}
}