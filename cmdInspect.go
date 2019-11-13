// +build !rm_basic_commands allcommands inspectcmd

package main

import (
	"fmt"
	"regexp"
	"samhofi.us/x/keybase"
	"strconv"
	"strings"
)

func init() {
	command := Command{
		Cmd:         []string{"inspect", "id"},
		Description: "$identifier - shows info about $identifier ($identifier is either username, messageId or team)",
		Help:        "",
		Exec:        cmdInspect,
	}

	RegisterCommand(command)
}

func cmdInspect(cmd []string) {
	if len(cmd) == 2 {
		regexIsNumeric := regexp.MustCompile(`^\d+$`)
		if regexIsNumeric.MatchString(cmd[1]) {
			// Then it must be a message id
			id, _ := strconv.Atoi(cmd[1])
			go printMessage(id)

		} else {
			go printUser(strings.ReplaceAll(cmd[1], "@", ""))
		}

	} else {
		printInfo(fmt.Sprintf("To inspect something use %sid <username/messageId>", config.Basics.CmdPrefix))
	}

}
func printMessage(id int) {
	chat := k.NewChat(channel)
	messages, err := chat.ReadMessage(id)
	if err == nil {
		var response StyledString
		if messages != nil && len((*messages).Result.Messages) > 0 {
			message := (*messages).Result.Messages[0].Msg
			var apiCast keybase.ChatAPI
			apiCast.Msg = &message
			response = formatOutput(apiCast)
		} else {
			response = config.Colors.Feed.Error.stylize("message not found")
		}
		printToView("Chat", response.string())
	}
}

func formatProofs(userLookup keybase.UserAPI) StyledString {
	messageColor := config.Colors.Message
	message := basicStyle.stylize("")
	for _, proof := range userLookup.Them[0].ProofsSummary.All {
		style := config.Colors.Feed.Success
		if proof.State != 1 {
			style = config.Colors.Feed.Error
		}
		proofString := style.stylize("Proof [$NAME@$SITE]: $URL\n")
		proofString = proofString.replace("$NAME", messageColor.SenderDefault.stylize(proof.Nametag))
		proofString = proofString.replace("$SITE", messageColor.SenderDevice.stylize(proof.ProofType))
		proofString = proofString.replace("$URL", messageColor.LinkURL.stylize(proof.HumanURL))
		message = message.append(proofString)
	}
	return message.appendString("\n")
}
func formatProfile(userLookup keybase.UserAPI) StyledString {
	messageColor := config.Colors.Message
	user := userLookup.Them[0]
	profileText := messageColor.Body.stylize("Name: $FNAME\nLocation: $LOC\nBio: $BIO\n")
	profileText = profileText.replaceString("$FNAME", user.Profile.FullName)
	profileText = profileText.replaceString("$LOC", user.Profile.Location)
	profileText = profileText.replaceString("$BIO", user.Profile.Bio)

	return profileText
}

func formatFollowState(userLookup keybase.UserAPI) StyledString {
	username := userLookup.Them[0].Basics.Username
	followSteps := followedInSteps[username]
	if followSteps == 1 {
		return config.Colors.Feed.Success.stylize("<Followed!>\n\n")
	} else if followSteps > 1 {
		var steps []string
		for head := username; head != ""; head = trustTreeParent[head] {
			steps = append(steps, fmt.Sprintf("[%s]", head))
		}
		trustLine := fmt.Sprintf("Indirect follow: <%s>\n\n", strings.Join(steps, " Followed by "))
		return config.Colors.Message.Body.stylize(trustLine)
	}

	return basicStyle.stylize("")

}

func formatFollowerAndFollowedList(username string, listType string) StyledString {
	messageColor := config.Colors.Message
	response := basicStyle.stylize("")
	bytes, _ := k.Exec("list-"+listType, username)
	bigString := string(bytes)
	lines := strings.Split(bigString, "\n")
	response = response.appendString(fmt.Sprintf("%s (%d): ", listType, len(lines)-1))
	for i, user := range lines[:len(lines)-1] {
		if i != 0 {
			response = response.appendString(", ")
		}
		response = response.append(messageColor.LinkKeybase.stylize(user))
		response = response.append(getUserFlags(user))
	}
	return response.appendString("\n\n")
}

func printUser(username string) {
	messageColor := config.Colors.Message

	userLookup, _ := k.UserLookup(username)

	response := messageColor.Header.stylize("[Inspecting `$USER`]\n")
	response = response.replace("$USER", messageColor.SenderDefault.stylize(username))
	response = response.append(formatProfile(userLookup))
	response = response.append(formatFollowState(userLookup))

	response = response.append(formatProofs(userLookup))
	response = response.append(formatFollowerAndFollowedList(username, "followers"))
	response = response.append(formatFollowerAndFollowedList(username, "following"))

	printToView("Chat", response.string())
}
