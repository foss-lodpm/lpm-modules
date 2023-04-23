package main

import "C"
import (
	"fmt"
	"unsafe"
)

//export lpm_entrypoint
func lpm_entrypoint(config_path_ptr *C.char, db_path_ptr *C.char, argc C.int, argv **C.char) {
	config_path := C.GoString(config_path_ptr)
	db_path := C.GoString(db_path_ptr)

	var args []string
	for i := 0; i < int(argc); i++ {
		argPtr := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + uintptr(i)*unsafe.Sizeof(*argv)))
		arg := C.GoString(*argPtr)
		args = append(args, arg)
	}

	fmt.Println("config_path:", config_path)
	fmt.Println("db_path:", db_path)
	fmt.Println("args:", args)
}

func main() {}
