// ignore
// +build allcommands walletcmd

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var walletConfirmationCode string
var walletConfirmationUser string
var walletTransactionAmnt string

func init() {
	command := Command{
		Cmd:         []string{"wallet", "confirm"},
		Description: "$user $amount / $user $confirmation - Send or confirm a wallet payment",
		Help:        "",
		Exec:        cmdWallet,
	}

	RegisterCommand(command)
}

func cmdWallet(cmd []string) {
	if len(cmd) < 3 {
		return
	}
	if cmd[0] == "wallet" {
		rand.Seed(time.Now().UnixNano())
		chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
			"abcdefghijklmnopqrstuvwxyz" +
			"0123456789")
		length := 8
		var b strings.Builder
		for i := 0; i < length; i++ {
			b.WriteRune(chars[rand.Intn(len(chars))])
		}
		walletConfirmationCode = b.String()
		walletConfirmationUser = cmd[1]
		walletTransactionAmnt = cmd[2]
		printInfo(fmt.Sprintf("To confirm sending %s to %s, type /confirm %s %s", cmd[2], cmd[1], cmd[1], walletConfirmationCode))

	} else if cmd[0] == "confirm" {
		if cmd[1] == walletConfirmationUser && cmd[2] == walletConfirmationCode {
			txWallet := k.NewWallet()
			wAPI, err := txWallet.SendXLM(walletConfirmationUser, walletTransactionAmnt, "")
			if err != nil {
				printError(fmt.Sprintf("There was an error with your wallet tx:\n\t%+v", err))
			} else {
				printInfo(fmt.Sprintf("You have sent %sXLM to %s with tx ID: %s", wAPI.Result.Amount, wAPI.Result.ToUsername, wAPI.Result.TxID))
			}

		} else {
			printError("There was an error validating your confirmation. Your wallet has been untouched.")
		}

	}

}
