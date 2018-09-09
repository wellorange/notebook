extern crate export;   // extern关键字用来导入外部的包
extern crate libloading;


mod io;                // 定义一个模块
use export::ee;        // use关键子用来使用导入的包


fn main() {
    ee::Hello();
   // io::Loadlib();
   // io::V8();
    io::helloc();
}