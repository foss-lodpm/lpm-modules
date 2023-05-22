package main

import "C"
import (
	"log"
	"lpm_builder/pkg/builder"
	common "lpm_builder/pkg/common"
	template "lpm_builder/pkg/template"
	"os"
	"path/filepath"
	"time"
	"unsafe"

	"github.com/urfave/cli/v2"
)

var MODULE_VERSION string = "1.0.0"

func main() {}

func handle_cli(args []string) {
	app := &cli.App{
		Name:     "lpm-builder",
		Version:  MODULE_VERSION,
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
		Action: func(c *cli.Context) error {
			templateDir := c.String("build")
			templateDir, err := filepath.Abs(templateDir)
			common.FailOnError(err)

			builder.StartBuilding(templateDir)

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "build",
				Required: false,
				Aliases:  []string{"b"},
				Usage:    "Takes path of the template directory.",
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "gen",
				Aliases:     []string{"g"},
				Usage:       "generates the build files that will be used for creating lod packages.",
				UsageText:   "lpm -m builder --generate(or -g) --name(or -n) {package_name}",
				Description: "Mostly used to bootstrap providing package build files/templates.",
				Action: func(c *cli.Context) error {
					template_name := c.String("name")
					out_path := c.String("out")

					if out_path != "" {
						err := os.MkdirAll(out_path, os.ModePerm)
						common.FailOnError(err, "Failed on creating directory")

						out_path += "/"
					}

					template.CreateDefault(c, template_name, out_path)

					return nil
				},
				Flags: []cli.Flag{
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

	err := app.Run(args)
	common.FailOnError(err, "")

}

//export lpm_entrypoint
func lpm_entrypoint(config_path_ptr *C.char, db_path_ptr *C.char, argc C.int, argv **C.char) {
	// 	config_path := C.GoString(config_path_ptr)
	// 	db_path := C.GoString(db_path_ptr)
	//
	// 	_ = config_path
	// 	_ = db_path

	// initialize logger
	common.Logger = log.New(os.Stdout, "[lpm-builder]: ", log.Ldate|log.Ltime)

	var args []string
	for i := 0; i < int(argc); i++ {
		argPtr := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + uintptr(i)*unsafe.Sizeof(*argv)))
		args = append(args, C.GoString(*argPtr))
	}

	handle_cli(args[2:])
}
