extern crate libc;

use std::env;
use libloading::{Library,Symbol};

pub fn Hello(){
    println!("hello io");
}

pub fn Command(){
    let args: Vec<String> = env::args().collect();
    println!("{:?}", args);
}

pub fn Openfile(){

}




pub fn Loadlib() {
    type AddFunc = fn(isize, isize) -> isize;
    type testFunc = fn(i32) -> i32;
    
    let lib = Library::new("../../../lib/target/debug/lib.dll").unwrap();
    let func: Symbol<AddFunc> =unsafe { lib.get(b"add").unwrap()};
    println!("3 + 4 = {}", func(3,4));
    
    let test: Symbol<testFunc> = unsafe {lib.get(b"test").unwrap()};
    println!("{}",test(100));
}


pub fn v8(){

}


extern {
    fn hello(input: libc::c_int) -> libc::c_int;
}

pub fn helloc(){
  
   let input = 8;
   let output = unsafe { hello(input) };
   println!("{} * 2 = {}\n", input, output);

}