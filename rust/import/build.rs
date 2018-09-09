extern crate cc;

fn main() {
    cc::Build::new()
        .file("./vender/v8/index.c")
        .compile("libdouble.a");
}