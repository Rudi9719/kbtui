// +build !rm_basic_commands allcommands tabcompletion
package main

import (
	"fmt"
	"regexp"
	"strings"

	"samhofi.us/x/keybase"
)

var (
	tabSlice []string
)

// This defines the handleTab function thats called by key bindind tab for the input control.
func handleTab() error {
	inputString, err := getInputString("Input")
	if err != nil {
		return err
	} else {
		// if you successfully get an input string, grab the last word from the string
		ss := regexp.MustCompile(`[ #]`).Split(inputString, -1)
		s := ss[len(ss)-1]
		// create a variable in which to store the result
		var resultSlice []string
		// if the word starts with a : its an emoji lookup
		if strings.HasPrefix(s, ":") {
			resultSlice = getEmojiTabCompletionSlice(s)
		} else {
			if strings.HasPrefix(s, "@") {
				// now in case the word (s) is a mention @something, lets remove it to normalize
				s = strings.Replace(s, "@", "", 1)
			}
			// now call get the list of all possible cantidates that have that as a prefix
			resultSlice = getChannelTabCompletionSlice(s)
		}
		rLen := len(resultSlice)
		lcp := longestCommonPrefix(resultSlice)
		if lcp != "" {
			originalViewTitle := getViewTitle("Input")
			newViewTitle := ""
			if rLen >= 1 && originalViewTitle != "" {
				if rLen == 1 {
					newViewTitle = originalViewTitle
				} else if rLen <= 5 {
					newViewTitle = fmt.Sprintf("%s|| %s", originalViewTitle, strings.Join(resultSlice, " "))
				} else if rLen > 5 {
					newViewTitle = fmt.Sprintf("%s|| %s +%d more", originalViewTitle, strings.Join(resultSlice[:6], " "), rLen-5)
				}
				setViewTitle("Input", newViewTitle)
				remainder := stringRemainder(s, lcp)
				writeToView("Input", remainder)
			}
		}
	}
	return nil
}

// Main tab completion functions
func getEmojiTabCompletionSlice(inputWord string) []string {
	// use the emojiSlice from emojiList.go and filter it for the input word
	resultSlice := filterStringSlice(emojiSlice, inputWord)
	return resultSlice
}
func getChannelTabCompletionSlice(inputWord string) []string {
	// use the tabSlice from above and filter it for the input word
	resultSlice := filterStringSlice(tabSlice, inputWord)
	return resultSlice
}

//Generator Functions (should be called externally when chat/list/join changes
func generateChannelTabCompletionSlice() {
	// fetch all members of the current channel and add them to the slice
	channelSlice := getCurrentChannelMembership()
	for _, m := range channelSlice {
		tabSlice = appendIfNotInSlice(tabSlice, m)
	}
}
func generateRecentTabCompletionSlice() {
	var recentSlice []string
	for _, s := range channels {
		if s.MembersType == keybase.TEAM {
			// its a team so add the topic name and channel name
			recentSlice = appendIfNotInSlice(recentSlice, s.TopicName)
			recentSlice = appendIfNotInSlice(recentSlice, s.Name)
		} else {
			//its a user, so clean the name and append
			recentSlice = appendIfNotInSlice(recentSlice, cleanChannelName(s.Name))
		}
	}
	for _, s := range recentSlice {
		tabSlice = appendIfNotInSlice(tabSlice, s)
	}
}

// Helper functions
func getCurrentChannelMembership() []string {
	var rs []string
	if channel.Name != "" {
		t := k.NewTeam(channel.Name)
		if testVar, err := t.MemberList(); err != nil {
			return rs // then this isn't a team, its a PM or there was an error in the API call
		} else {
			for _, m := range testVar.Result.Members.Owners {
				rs = append(rs, fmt.Sprintf("%+v", m.Username))
			}
			for _, m := range testVar.Result.Members.Admins {
				rs = append(rs, fmt.Sprintf("%+v", m.Username))
			}
			for _, m := range testVar.Result.Members.Writers {
				rs = append(rs, fmt.Sprintf("%+v", m.Username))
			}
			for _, m := range testVar.Result.Members.Readers {
				rs = append(rs, fmt.Sprintf("%+v", m.Username))
			}
		}
	}
	return rs
}
func filterStringSlice(ss []string, fv string) []string {
	var rs []string
	for _, s := range ss {
		if strings.HasPrefix(s, fv) {
			rs = append(rs, s)
		}
	}
	return rs
}
func longestCommonPrefix(ss []string) string {
	// cover the case where the slice has no or one members
	switch len(ss) {
	case 0:
		return ""
	case 1:
		return ss[0]
	}
	// all strings are compared by bytes here forward (TBD unicode normalization?)
	// establish min, max lenth members of the slice by iterating over the members
	min, max := ss[0], ss[0]
	for _, s := range ss[1:] {
		switch {
		case s < min:
			min = s
		case s > max:
			max = s
		}
	}
	// then iterate over the characters from min to max, as soon as chars don't match return
	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			return min[:i]
		}
	}
	// to cover the case where all members are equal, just return one
	return min
}
func stringRemainder(aStr, bStr string) string {
	var long, short string
	//figure out which string is longer
	switch {
	case len(aStr) < len(bStr):
		short = aStr
		long = bStr
	default:
		short = bStr
		long = aStr
	}
	// iterate over the strings using an external iterator so we don't lose the value
	i := 0
	for i < len(short) && i < len(long) {
		if short[i] != long[i] {
			// the strings aren't equal so don't return anything
			return ""
		}
		i++
	}
	// return whatever's left of the longer string
	return long[i:]
}
func appendIfNotInSlice(ss []string, s string) []string {
	for _, element := range ss {
		if element == s {
			return ss
		}
	}
	return append(ss, s)
}
