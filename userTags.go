package main

import (
	"fmt"
	"strings"
)

var followedInSteps = make(map[string]int)
var trustTreeParent = make(map[string]string)

func clearFlagCache() {
	followedInSteps = make(map[string]int)
	trustTreeParent = make(map[string]string)
}

var maxDepth = 4

func generateFollowersList() {
	// Does a BFS of followedInSteps
	queue := []string{k.Username}
	printInfo("Generating Tree of Trust...")
	lastDepth := 1
	for len(queue) > 0 {
		head := queue[0]
		queue = queue[1:]
		depth := followedInSteps[head] + 1
		if depth > maxDepth {
			continue
		}
		if depth > lastDepth {
			printInfo(fmt.Sprintf("Trust generated at Level #%d", depth-1))
			lastDepth = depth
		}

		bytes, _ := k.Exec("list-following", head)
		bigString := string(bytes)
		following := strings.Split(bigString, "\n")
		for _, user := range following {
			if followedInSteps[user] == 0 && user != k.Username {
				followedInSteps[user] = depth
				trustTreeParent[user] = head
				queue = append(queue, user)
			}
		}
	}
	printInfo(fmt.Sprintf("Trust-level estabilished for %d users", len(followedInSteps)))
}

func getUserFlags(username string) StyledString {
	tags := ""
	followDepth := followedInSteps[username]
	if followDepth == 1 {
		tags += fmt.Sprintf(" %s", config.Formatting.IconFollowingUser)
	} else if followDepth > 1 {
		tags += fmt.Sprintf(" %s%d", config.Formatting.IconIndirectFollowUser, followDepth-1)
	}
	return config.Colors.Message.SenderTags.stylize(tags)
}
