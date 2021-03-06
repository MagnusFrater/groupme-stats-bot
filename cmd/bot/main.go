package main

import (
	"flag"
	"log"
	"os"
	"strings"

	groupmestatsbot "github.com/MagnusFrater/groupme-stats-bot"
)

func main() {
	accessToken := flag.String("accessToken", "", "GroupMe API client access token")
	botID := flag.String("botID", "", "GroupMe Bot ID")
	groupID := flag.String("groupID", "", "GroupMe Group ID")

	limit := flag.Int("limit", 5, "number of items to list")
	blacklist := flag.String("blacklist", "", "blacklist of comma-delimited User IDs")
	logOnly := flag.Bool("logOnly", false, "toggle output to log instead of sent via the GroupMe Bot")

	commands := flag.String("commands", "", "list of comma-delimited bot commands")
	flag.Parse()

	// initialize GroupMe Stats Bot
	statsBot, err := groupmestatsbot.New(*accessToken, *botID, *groupID, *limit, strings.Split(*blacklist, ","))
	if err != nil {
		flag.Usage()
		log.Panic(err)
	}

	// get commands
	commandList := strings.Fields(*commands)

	// show usage if zero commands toggled
	if len(commandList) == 0 {
		flag.Usage()
		printAvailableCommands()
		os.Exit(1)
	}

	// run commands
	ranCommand := false
	for _, command := range commandList {
		didRun, err := statsBot.Command(command, *logOnly)
		if err != nil {
			log.Panic(err)
		}

		if didRun == true {
			ranCommand = didRun
		}
	}

	// if no command was ran, they probably misspelled it
	// or just straight up didn't call a real command
	if !ranCommand {
		flag.Usage()
		log.Println("No valid command was toggled.")
		printAvailableCommands()
		os.Exit(1)
	}
}

func printAvailableCommands() {
	str := "Commands must contain at least one of: "
	availableCommands := groupmestatsbot.GetAllCommands()
	for i, command := range availableCommands {
		str += command

		// don't put comma after last command
		if i < len(availableCommands)-1 {
			str += ", "
		}
	}
	log.Println(str)
}
