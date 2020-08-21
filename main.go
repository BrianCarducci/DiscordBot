package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/BrianCarducci/DiscordBot/bot_error"
	"github.com/BrianCarducci/DiscordBot/constants"
	"github.com/BrianCarducci/DiscordBot/constants/commands"


	"github.com/BrianCarducci/DiscordBot/utils"

	"github.com/bwmarrin/discordgo"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	setupBot()
}

func getAWSSecrets(secretKeys []string) (map[string]string, error) {
	secretName := "JeffBot"
	region := "us-east-1"

	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New(),
                                  aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
				case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				return nil, errors.New(secretsmanager.ErrCodeDecryptionFailure + ": " + aerr.Error())

				case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				return nil, errors.New(secretsmanager.ErrCodeInternalServiceError + ": " + aerr.Error())

				case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				return nil, errors.New(secretsmanager.ErrCodeInvalidParameterException + ": " + aerr.Error())

				case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				return nil, errors.New(secretsmanager.ErrCodeInvalidRequestException + ": " + aerr.Error())

				case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				return nil, errors.New(secretsmanager.ErrCodeResourceNotFoundException + ": " + aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return nil, err
		}
	}

	var secretStringifiedMap string
	if result.SecretString != nil {
		secretStringifiedMap = *result.SecretString
	} else {
		return nil, errors.New("The secret value is nil")
	}
	
	secretMap := map[string]string{}
	if err := json.Unmarshal([]byte(secretStringifiedMap), &secretMap); err != nil {
		return nil, err
	}

	for _, secKey := range(secretKeys) {
		_, ok := secretMap[secKey]
		if !ok {
			// TODO: Maybe loop through them all and return an error with all that don't exist, but the first is fine for now.
			return nil, errors.New("Secret " + secKey + " does not exist in Secrets Manager.")
		}
	}

	return secretMap, nil
}

func getAPITokens(secretNames []string, isLocal *bool) (map[string]string, error) {
	if *isLocal {
		tokenVals := map[string]string{}
		for _, secName := range(secretNames) {
			secVal := strings.TrimSpace(os.Getenv(secName))
			if len(secVal) == 0 {
				return nil, errors.New("Token " + secName + " is not set as an environment variable.")
			}
			tokenVals[secName] = secVal
		}
		return tokenVals, nil
	}

	return getAWSSecrets(secretNames)
}

func setupBot() {
	// We can keep this here in case we want to add any other flags in the future
	isLocalPtr := flag.Bool("local", false, "If set, the bot will be run locally (using environment variables). Else, it will use AWS Secrets Manager")
	flag.Parse()

	// Ideally make a map or something for a token's env variable name and value..
	envNames := []string{"DISCORD_TOKEN", "GOOGLE_TOKEN"}
	apiTokens, err := getAPITokens(envNames, isLocalPtr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	commands.GeoLocator.Token = apiTokens["GOOGLE_TOKEN"]

	discord, err := discordgo.New("Bot " + apiTokens["DISCORD_TOKEN"])

	// Remove the tokens from memory as soon as possible. Not sure if this helps but we'll do it for now I suppose.
	// apiTokens = nil
	if err != nil {
		fmt.Println("Could not instantiate bot. Error: " + err.Error())
		discord.Close()
		os.Exit(1)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection: " + err.Error())
		discord.Close()
		os.Exit(1)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println(constants.BotName + " is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	err := utils.RunCommand(s, m)
	if _, ok := err.(*bot_error.BotError); ok {
		// Might want to check if the code is 0. Hopefully, casting doesn't make a new BotError with code 0..
		return
	} else if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
}
