use wasm_bindgen::prelude::*;
extern crate js_sys;

#[wasm_bindgen]
extern {
    #[wasm_bindgen(js_namespace = console)]
    fn log(s: &str);
}

#[wasm_bindgen(start)]
pub fn main() {
    // Called by Javascript
    log("Hello World!!!")
}