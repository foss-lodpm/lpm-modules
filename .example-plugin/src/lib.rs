use std::{env, ffi::CStr};

#[no_mangle]
pub extern "C" fn lpm_entrypoint(
    config_path: *const std::os::raw::c_char,
    db_path: *const std::os::raw::c_char,
) {
    let config_path = unsafe { CStr::from_ptr(config_path).to_str().unwrap() };
    let db_path = unsafe { CStr::from_ptr(db_path).to_str().unwrap() };
    let args: Vec<String> = env::args().collect();

    println!("config_path: {}", config_path);
    println!("db_path: {}", db_path);
    println!("args: {:?}", args);
}
