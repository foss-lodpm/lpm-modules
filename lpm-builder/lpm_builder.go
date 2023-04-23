// go build -o liblpm_builder.so -buildmode=c-shared lpm_builder.go

package main

import "C"
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/urfave/cli/v2"
)

func main() {}

type Version struct {
	ReadableFormat string  `json:"readable_format"`
	Major          int     `json:"major"`
	Minor          int     `json:"minor"`
	Patch          int     `json:"patch"`
	Tag            *string `json:"tag"`
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

					if out_path != "" {
						err := os.MkdirAll(out_path, os.ModePerm)
						failOnError(err, "Failed on creating directory")

						out_path += "/"
					}

					if is_default && is_interactive {
						cli.ShowSubcommandHelp(c)
					} else if is_default {
						create_default_template(c, template_name, out_path)
					} else if is_interactive {
						create_template_interactively(c, template_name, out_path)
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

func set_version_readable_format(version *Version) {
	readable_format := fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)

	if version.Tag != nil {
		readable_format += fmt.Sprintf("-%s", *version.Tag)
	}

	version.ReadableFormat = readable_format

}

func create_default_template(c *cli.Context, template_name string, out_path string) {
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

	set_version_readable_format(&template.Version)

	template_json, err := json.MarshalIndent(template, "", " ")
	failOnError(err, "Failed on serializing template to json string")

	err = ioutil.WriteFile(out_path+"template.json", template_json, 0644)
	failOnError(err, "Failed on writing to file")
}

func read_and_fill_string(prompt string, reader *bufio.Reader, to_be_filled *string) {
	fmt.Printf(prompt)
	value, err := reader.ReadString('\n')
	failOnError(err, "Failed on reading value from stdin")
	*to_be_filled = strings.TrimSpace(value)

}

func create_template_interactively(c *cli.Context, template_name string, out_path string) {
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

	reader := bufio.NewReader(os.Stdin)

	read_and_fill_string(fmt.Sprintf("Description of '%s': ", template_name), reader, &template.Description)
	read_and_fill_string(fmt.Sprintf("Maintainer of '%s': ", template_name), reader, &template.Maintainer)
	read_and_fill_string(fmt.Sprintf("Source package of '%s': ", template_name), reader, &template.SourcePkg)
	read_and_fill_string(fmt.Sprintf("Repository of '%s': ", template_name), reader, &template.Repository)
	read_and_fill_string(fmt.Sprintf("Homepage of '%s': ", template_name), reader, &template.Homepage)
	read_and_fill_string(fmt.Sprintf("Architecture of '%s': ", template_name), reader, &template.Arch)
	read_and_fill_string(fmt.Sprintf("Kind of '%s': ", template_name), reader, &template.Kind)
	read_and_fill_string(fmt.Sprintf("License of '%s': ", template_name), reader, &template.License)

	fmt.Printf("Tags of '%s': ", template_name)
	value, err := reader.ReadString('\n')
	failOnError(err, "Failed on reading value from stdin")
	template.Tags = strings.Split(strings.TrimSpace(value), " ")

	set_version_readable_format(&template.Version)

	template_json, err := json.MarshalIndent(template, "", " ")
	failOnError(err, "Failed on serializing template to json string")

	err = ioutil.WriteFile(out_path+"template.json", template_json, 0644)
	failOnError(err, "Failed on writing to file")
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

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println("Error: ", err, msg)
		os.Exit(1)

	}
}
