package main

import "C"
import (
	"fmt"
	"os"
)

//export lpm_entrypoint
func lpm_entrypoint(config_path_ptr *C.char, db_path_ptr *C.char) {
	config_path := C.GoString(config_path_ptr)
	db_path := C.GoString(db_path_ptr)
	args := os.Args

	fmt.Println("config_path:", config_path)
	fmt.Println("db_path:", db_path)
	fmt.Println("args:", args)
}

func main() {}