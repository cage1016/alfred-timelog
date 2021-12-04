package main

import (
	"fmt"
	"log"
	"os"
	"unicode"

	aw "github.com/deanishe/awgo"

	"github.com/cate1016/timetrack-alfred-workflows/handler"
)

var wf *aw.Workflow

func init() {
	wf = aw.New()
}

func run() {
	args := wf.Args()
	if len(args) == 0 {
		exitWithError("please provide some input 👀")
	}

	handlers := map[string]func(*aw.Workflow, []string) (string, error){
		"add":    handler.Add,
		"auth":   handler.DoAuthorize,
		"deauth": handler.DoDeAuthorize,
		"setup":  handler.DoSetup,
	}

	h, found := handlers[args[0]]
	if !found {
		exitWithError("command not recognized 👀")
	}

	msg, err := h(wf, args[1:])
	if err != nil {
		exitWithError(err.Error())
		os.Exit(1)
	}

	if msg != "" {
		fmt.Print(msg)
	}
}

func main() {
	wf.Run(run)
}

func exitWithError(msg string) {
	fmt.Print(capitalize(msg))
	log.Print(msg)
	os.Exit(1)
}

func capitalize(msg string) string {
	return string(unicode.ToUpper(rune(msg[0]))) + msg[1:]
}
