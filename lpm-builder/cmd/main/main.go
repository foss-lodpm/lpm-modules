package main

import "C"
import (
	common "lpm_builder/pkg/common"
	template "lpm_builder/pkg/template"
	"os"
	"time"
	"unsafe"

	"github.com/urfave/cli/v2"
)

func main() {}

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
						common.FailOnError(err, "Failed on creating directory")

						out_path += "/"
					}

					if is_default && is_interactive {
						cli.ShowSubcommandHelp(c)
					} else if is_default {
						template.CreateDefault(c, template_name, out_path)
					} else if is_interactive {
						template.CreateInteractively(c, template_name, out_path)
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
