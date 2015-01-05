package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	ctx        *Context = &Context{Fallback: "", Envs: map[string]string{}}
	inputPath  string   = ""
	outputPath string   = ""
)

type Context struct {
	Envs     map[string]string
	Input    io.Reader
	Output   io.Writer
	Fallback string
}

func (c *Context) Fetch(arg string) string {
	if c.Envs[arg] == "" {
		return c.Fallback
	}

	return c.Envs[arg]
}

func (c *Context) Execute() {
	tmpl := ""
	scanner := bufio.NewScanner(c.Input)

	for scanner.Scan() {
		if tmpl == "" {
			tmpl = scanner.Text()
		} else {
			tmpl += "\n" + scanner.Text()
		}
	}

	t := template.New("envtmpl")

	t.Parse(tmpl)

	t.Execute(c.Output, c.Envs)
}

func init() {
	flag.StringVar(&inputPath, "i", "", "Input template file path")
	flag.StringVar(&inputPath, "input", "", "Input template file path")
	flag.StringVar(&outputPath, "o", "", "Output path")
	flag.StringVar(&outputPath, "output", "", "Output path")
}

func main() {
	flag.Parse()

	if inputPath != "" {
		absInputPath, _ := filepath.Abs(inputPath)
		input, err := os.OpenFile(absInputPath, os.O_RDONLY, 0666)

		if err != nil {
			panic(err.Error())
		}
		defer input.Close()

		ctx.Input = input
	} else {
		ctx.Input = os.Stdin
	}

	if outputPath != "" {
		absOutputPath, _ := filepath.Abs(outputPath)
		output, err := os.OpenFile(absOutputPath, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			panic(err.Error())
		}
		defer output.Close()

		ctx.Output = output
	} else {
		ctx.Output = os.Stdout
	}

	for _, env := range os.Environ() {
		r := strings.Split(env, "=")

		if len(r) > 1 {
			ctx.Envs[r[0]] = r[1]
		}
	}

	ctx.Execute()
}
