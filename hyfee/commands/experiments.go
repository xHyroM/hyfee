package commands

import (
	"hyros_coffee/handler"
	"hyros_coffee/hyfee"
	"hyros_coffee/utils"
	"reflect"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/oauth2"
	"github.com/disgoorg/snowflake/v2"
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
							Autocomplete: true,
						},
						discord.ApplicationCommandOptionInt{
							Name: "member_count",
							Description: "Member count (not required, will be done automatically if widgets)",
							Required: false,
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

func getExperimentsHandler(_ *hyfee.Bot) handler.CommandHandler {
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
			embed.AddField("Populations", experiment.FormatPopulations(), false)
		}

		if len(experiment.Rollout.Overrides) > 0 {
			embed.AddField("Overrides", experiment.FormatOverrides(), false)
		}

		if len(experiment.Rollout.OverridesFormatted) > 0 {
			embed.AddField("Overrides Formatted", experiment.FormatOverridesFormatted(), false)
		}

		return event.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{embed.Build()},
		})
	}
}

func eligibleExperimentsHandler(bot *hyfee.Bot) handler.CommandHandler {
	return func(event *events.ApplicationCommandInteractionCreate) error {
		experimentId := event.SlashCommandInteractionData().String("experiment")
		memberCount, isMemberCountProvided := event.SlashCommandInteractionData().OptInt("member_count")
		guildId, _ := snowflake.Parse(event.SlashCommandInteractionData().String("guild"))

		user, err := bot.Database.Get(event.User().ID.String())
		guild := discord.OAuth2Guild{
			Features: []discord.GuildFeature{},
		}

		if err == nil {
			guild = utils.GetGuild(bot.OAuth2Client, bot.Database, event.User().ID, guildId, oauth2.Session{
				AccessToken: user.AccessToken,
				RefreshToken: user.RefreshToken,
				Scopes: user.Scopes,
				TokenType: user.TokenType,
				Expiration: user.Expiration,
			})
		}	

		if err != nil {
			return event.CreateMessage(discord.MessageCreate{
				Content: "Error: Failed to parse guild id",
				Flags: discord.MessageFlagEphemeral,
			})
		}

		body := discord.Guild{
			ID: guildId,
			Features: guild.Features,
		}

		if isMemberCountProvided {
			body.MemberCount = memberCount
		}

		eligible, err := utils.IsExperimentEligible(experimentId, body)
		if err != nil {
			return event.CreateMessage(discord.MessageCreate{
				Content: "Error: " + err.Error(),
				Flags: discord.MessageFlagEphemeral,
			})
		}

		isFilters := len(eligible.Filters) > 0

		embed := discord.NewEmbedBuilder().
			SetTitle("Experiment Eligiblity Check").
			AddField("Experiment Id", experimentId, true).
			AddField("Guild Id", guildId.String(), true)

		if isFilters {
			embed.AddField("Eligible", "Possible", true)
		} else {
			embed.AddField("Eligible", utils.FormatBool(eligible.Eligible), true)
		}

		if !reflect.ValueOf(eligible.Bucket).IsZero() {
			embed.AddField("Bucket", eligible.Bucket.Format(eligible.Bucket.Id), false)
		}

		if isFilters {
			filters := []string{}

			for _, filter := range eligible.Filters {
				filters = append(filters, filter.Format())
			}

			embed.AddField("Filters (Guild must comply)", strings.Join(filters, " and ") + "\n", false)
		}

		if isFilters {
			embed.SetColor(0xfcba03)
		} else {
			if eligible.Eligible {
				embed.SetColor(0x42f554)
			} else {
				embed.SetColor(0xeb4034)
			}
		}

		return event.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{embed.Build()},
		})
	}
}

func autocompleteHandler(bot *hyfee.Bot) handler.AutocompleteHandler {
	return func(e *events.AutocompleteInteractionCreate) error {
		focusedOption := utils.GetFocusedOption(e)
		
		if focusedOption.Name == "experiment" {
			return autocompleteExperiments(e)
		}

		return autocompleteGuilds(bot, e)
	}
}

func autocompleteGuilds(bot *hyfee.Bot, e *events.AutocompleteInteractionCreate) error {
	option := e.Data.String("guild")

	user, err := bot.Database.Get(e.User().ID.String())
	if err != nil {
		return err
	}

	guilds := utils.GetGuilds(bot.OAuth2Client, bot.Database, e.User().ID, oauth2.Session{
		AccessToken: user.AccessToken,
		RefreshToken: user.RefreshToken,
		Scopes: user.Scopes,
		TokenType: user.TokenType,
		Expiration: user.Expiration,
	})

	result := []discord.AutocompleteChoice{}

	for _, guild := range guilds {
		if option == "" || strings.Contains(strings.ToLower(guild.Name), strings.ToLower(option)) {
			result = append(result, discord.AutocompleteChoiceString{
				Name: utils.IfThenElse(
					len(guild.Name) > 25,
					func() string {
						return guild.Name[:23]  + ".."
					},
					func() string {
						return guild.Name
					}),
				Value: guild.ID.String(),
			})
		}

		if len(result) == 25 {
			break
		}
	}

	return e.AutocompleteResult(result)
}

func autocompleteExperiments(e *events.AutocompleteInteractionCreate) error {
	option := e.Data.String("experiment")
	subcommand := e.Data.SubCommandName

	var experiments []utils.ExperimentKey
	if *subcommand == "eligible" {
		experiments = utils.GetExperimentKeys("&kind=guild&has_rollout=true")
	} else {
		experiments = utils.GetExperimentKeys("")
	}

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

	return e.AutocompleteResult(result)
}