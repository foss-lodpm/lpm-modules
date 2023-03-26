use std::ffi::CStr;

#[no_mangle]
pub extern "C" fn lpm_entrypoint(
    config_path: *const std::os::raw::c_char,
    db_path: *const std::os::raw::c_char,
    argc: std::os::raw::c_uint,
    argv: *const std::os::raw::c_void,
) {
    let config_path = unsafe { CStr::from_ptr(config_path).to_str().unwrap() };
    let db_path = unsafe { CStr::from_ptr(db_path).to_str().unwrap() };

    let argv = argv as *const *const std::os::raw::c_char;
    let args = unsafe { std::slice::from_raw_parts(argv, argc as usize) }
        .iter()
        .map(|&arg| unsafe { CStr::from_ptr(arg) })
        .map(|cstr| cstr.to_string_lossy().to_string())
        .collect::<Vec<String>>();

    println!("config_path: {}", config_path);
    println!("db_path: {}", db_path);
    println!("args: {:?}", args);
}
