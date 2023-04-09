// go build -o liblpm_builder.so -buildmode=c-shared lpm_builder.go

package main

import "C"
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"unsafe"

	"github.com/urfave/cli/v2"
)

func main() {}

type Version struct {
	ReadableFormat string `json:"readable_format"`
	Major          int    `json:"major"`
	Minor          int    `json:"minor"`
	Patch          int    `json:"patch"`
	Tag            string `json:"tag"`
}

type Dependency struct {
	Name    string  `json:"name"`
	Version Version `json:"version"`
}

type Template struct {
	Name                string       `json:"name"`
	Description         string       `json:"description"`
	Maintainer          string       `json:"maintainer"`
	SourcePkg           string       `json:"source_pkg"`
	Repository          string       `json:"repository"`
	Homepage            string       `json:"homepage"`
	Arch                string       `json:"arch"`
	Kind                string       `json:"kind"`
	Tags                []string     `json:"tags"`
	Version             Version      `json:"version"`
	License             string       `json:"license"`
	RuntimeDependencies []Dependency `json:"runtime_dependencies"`
	RuntimeSuggestions  []Dependency `json:"runtime_suggestions"`
	BuildDependencies   []Dependency `json:"build_dependencies"`
}

func handle_cli(args []string) {
	app := &cli.App{
		Name:     "lpm-builder",
		Version:  "v0.0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Onur Özkan",
				Email: "contact@onurozkan.dev",
			},
		},
		Copyright: "(c) 2023 Free and Open Source Lod Package Manager",
		HelpName:  "lpm-builder",
		Usage:     "Builder module for Lod Package Manager",
		UsageText: "lpm -m builder [global options] command [command options] [arguments...]",
		Commands: []*cli.Command{
			{
				Name:        "gen",
				Aliases:     []string{"g"},
				Usage:       "generates the build files that will be used for creating lod packages.",
				UsageText:   "lpm -m builder --generate(or -g) --name(or -n) {package_name}",
				Description: "Mostly used to bootstrap providing package build files/templates.",
				Action: func(c *cli.Context) error {
					is_default := c.Bool("default")
					is_interactive := c.Bool("interactive")

					template_name := c.String("name")
					out_path := c.String("out")

					if is_default && is_interactive {
						cli.ShowSubcommandHelp(c)
					} else if is_default {
						create_default_template(c, template_name, out_path)
					} else if is_interactive {
						fmt.Println("execute function for interactive mode")
					} else {
						cli.ShowSubcommandHelp(c)
					}

					return nil
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "default",
						Aliases: []string{"d"},
						Usage:   "Generates the template with default/empty values.",
					},
					&cli.BoolFlag{
						Name:    "interactive",
						Aliases: []string{"i"},
						Usage:   "Generates the template interactively via the user prompts.",
					},
					&cli.StringFlag{
						Name:     "name",
						Required: true,
						Aliases:  []string{"n"},
						Usage:    "Name of the template.",
					},
					&cli.StringFlag{
						Name:    "out",
						Aliases: []string{"o"},
						Usage:   "Output directory path.",
					},
				},
			},
		},
	}

	app.Run(args)

}

func create_default_template(c *cli.Context, template_name string, out_path string) {
	if out_path != "" {
		err := os.MkdirAll(out_path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}

		out_path += "/"
	}

	template := &Template{
		Name:                template_name,
		Description:         "",
		Maintainer:          "",
		SourcePkg:           "",
		Repository:          "",
		Homepage:            "",
		Arch:                "",
		Kind:                "",
		Tags:                []string{},
		Version:             Version{},
		License:             "",
		RuntimeDependencies: []Dependency{},
		RuntimeSuggestions:  []Dependency{},
		BuildDependencies:   []Dependency{},
	}

	template_json, err := json.MarshalIndent(template, "", " ")
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile(out_path+"template.json", template_json, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

//export lpm_entrypoint
func lpm_entrypoint(config_path_ptr *C.char, db_path_ptr *C.char, argc C.int, argv **C.char) {
	// 	config_path := C.GoString(config_path_ptr)
	// 	db_path := C.GoString(db_path_ptr)
	//
	// 	_ = config_path
	// 	_ = db_path

	var args []string
	for i := 0; i < int(argc); i++ {
		argPtr := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + uintptr(i)*unsafe.Sizeof(*argv)))
		args = append(args, C.GoString(*argPtr))
	}

	handle_cli(args[2:])
}
