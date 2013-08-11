package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

// Command based execution, inspired by Go source.
type Command struct {
	Run func(args []string)

	Name  string
	usage string
	Short string
	Long  string

	Flag flag.FlagSet
}

func (c *Command) UsageLine() string {
	return c.Name + " " + c.usage
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine())
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

var commands = map[string]*Command{
	cmdBuild.Name:  cmdBuild,
	cmdServer.Name: cmdServer,
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
		return
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	cmdLower := strings.ToLower(args[0])
	if cmd, ok := commands[cmdLower]; ok {
		cmd.Flag.Usage = cmd.Usage
		cmd.Flag.Parse(args[1:])
		cmd.Run(cmd.Flag.Args())
		return
	}

	fmt.Fprintf(os.Stderr, "builder: unknown subcommand %q\nRun 'builder help' for usage.\n", args[0])
}

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("cmd")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

var usageTemplate = `Frost is a tool for generating static websites.

Usage:

	builder command [arguments]

The commands are:
{{range .}}{{if .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use 'builder help [command]' for more information about a command.
`

func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, commands)
}

func usage() {
	printUsage(os.Stderr)
}

var helpTemplate = `{{if .Runnable}}usage: builder {{.UsageLine}}

{{end}}{{.Long}}
`

func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		return
	}

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: builder help command\n\nToo many arguments given.\n")
		return
	}

	cmdName := strings.ToLower(args[0])
	if cmd, ok := commands[cmdName]; ok {
		tmpl(os.Stdout, helpTemplate, cmd)
		return
	}

	fmt.Fprintf(os.Stderr, "Unknown command %#q. Run 'builder help' for a list of commands.\n", cmdName)
}
