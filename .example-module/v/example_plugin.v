[export: 'lpm_entrypoint']
fn lpm_entrypoint(db_path &char, argc int, argv &&char) {
	println('db_path: ${unsafe { cstring_to_vstring(db_path) }}')

	for i in 0..argc {
		println('arg[$i] ${unsafe { cstring_to_vstring(argv[i]) }}')
	}
}
