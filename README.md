# kbtui
Keybase TUI written in Go using [@dxb](https://keybase.io/dxb)'s 
Keybase [bot framework](https://godoc.org/samhofi.us/x/keybase).
It started as a joke, then a bash script, and now here it is!


## Features
* Dark Mode (or rather mode based on Terminal Theme)
* Read and reply to messages
* Feed view to see mentions
* Stream view to see all incoming messages
* List view to show activity
* Chat view to interact with the current channel

## Todo
* Reactions to messages
* Mark unread messages in the List view
* Pretty format headers in List view
* Twitter-style feed reading public messages
* Track multiple conversations at once
* Auto #general teams when not given a channel

### Building and Running
```
go get ./
go build
./kbtui
```
Or
```
go get ./
go run main.go
```
Occasionally when @dxb updates his API it will be necessary to run 
`go get -u ./`