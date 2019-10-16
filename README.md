# kbtui
Keybase TUI written in Go using [@dxb](https://keybase.io/dxb)'s 
Keybase [bot framework](https://godoc.org/samhofi.us/x/keybase).
It started as a joke, then a bash script, and now here it is!

For support or suggestions check out the [kbtui team](https://keybase.io/team/kbtui)

## Features
* Dark Mode (or rather mode based on Terminal Theme)
* Read and reply to messages
* Feed view to see mentions
* Stream view to see all incoming messages
* List view to show activity
* Chat view to interact with the current channel
* Marks unread messages in the List view
* Reactions to messages
* Auto #general teams when not given a channel
* Pretty format headers in List view from window size
* Message editing
* Twitter-style feed reading public messages
* Message replies

## Todo
* Track multiple conversations at once


### Building and Running
Easiest Way:
```
go get -u github.com/rudi9719/kbtui
```
Or you can do the following:
```
go get ./
go run build.go
go run build.go {build, buildBeta... etc}
./kbtui
```

You may see an error with `go get ./` about PATHs, that may be safely ignored.

If you see an error about a missing dependancy during a build, you'll want to resolve that.


Occasionally when @dxb updates his API it will be necessary to run 
`go get -u ./`
