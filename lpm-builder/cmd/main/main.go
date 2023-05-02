package main

import "C"
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	common "lpm_builder/pkg/common"
	template "lpm_builder/pkg/template"
	"os"
	"os/exec"
	"path/filepath"
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
		Action: func(c *cli.Context) error {
			template_dir := c.String("build")

			// TODO
			// Provide custom bash functions that can be used in move_pkg script

			template_json_content, err := ioutil.ReadFile(template_dir + "/template")
			common.FailOnError(err, "Failed reading template json file")

			var template template.Template
			err = json.Unmarshal(template_json_content, &template)
			common.FailOnError(err, "Failed reading template json file")

			fmt.Printf("%#v\n++++++++++++++++\n\n", template)

			tmp_pkg_dir := filepath.Join("/tmp", template.Name)

			err = os.Setenv("BUILD_ROOT", tmp_pkg_dir)
			common.FailOnError(err, "Couldn't set $BUILD_ROOT.")

			stage0Scripts := filepath.Join(template_dir, "stage0")

			validateChecksumBash := `
			function validate_checksum {
				# Get the file path and checksum from the arguments
				file_path="$1"
				expected_checksum="$2"

				# Calculate the actual checksum of the file
				actual_checksum="$(sha256sum "$file_path" | awk '{print $1}')"

				# Compare the actual and expected checksums
				if [[ "$actual_checksum" == "$expected_checksum" ]]; then
					echo "Checksum validation successful"
					return 0
				else
					echo "Checksum validation failed"
					rm "$1"
					exit 1
				fi
			}`

			baseScript := fmt.Sprintf(
				`
				set -e

				%s
				`,
				validateChecksumBash)

			err = os.MkdirAll(tmp_pkg_dir, os.ModePerm)
			common.FailOnError(err, "Couldn't create temporary directory for building lod package.")

			tmp_src_dir := filepath.Join(tmp_pkg_dir, "src")
			err = os.MkdirAll(tmp_src_dir, os.ModePerm)
			common.FailOnError(err, "Couldn't create source directory for downloading/building package source.")

			initScript := fmt.Sprintf(`
			%s

			target_script=$(<%s)
			eval "$target_script"
			`, baseScript, stage0Scripts+"/init")

			cmd := exec.Command("bash", "-c", initScript)
			cmd.Dir = tmp_pkg_dir
			_, err = cmd.Output()
			common.FailOnError(err, "Couldn't execute init script from template directory.")

			cmd = exec.Command("/bin/bash", stage0Scripts+"/build")
			cmd.Dir = tmp_src_dir
			_, err = cmd.Output()
			common.FailOnError(err, "Couldn't execute build script from template directory.")

			cmd = exec.Command("/bin/bash", stage0Scripts+"/install_files")
			cmd.Dir = tmp_pkg_dir
			_, err = cmd.Output()
			common.FailOnError(err, "Couldn't execute install_files script from template directory.")

			// err = os.RemoveAll(tmp_pkg_dir)
			// common.FailOnError(err, "Failed cleaning the temporary files of building lod package.")

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
						err := cli.ShowSubcommandHelp(c)
						common.FailOnError(err, "")
					} else if is_default {
						template.CreateDefault(c, template_name, out_path)
					} else if is_interactive {
						template.CreateInteractively(c, template_name, out_path)
					} else {
						err := cli.ShowSubcommandHelp(c)
						common.FailOnError(err, "")
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

	var args []string
	for i := 0; i < int(argc); i++ {
		argPtr := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + uintptr(i)*unsafe.Sizeof(*argv)))
		args = append(args, C.GoString(*argPtr))
	}

	handle_cli(args[2:])
}
