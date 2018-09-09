#[no_mangle]
pub extern "C" fn add(a: isize, b: isize) -> isize {
    a + b
}

#[no_mangle]
pub extern  fn test(a: i32) -> i32 {
  return a*2
}