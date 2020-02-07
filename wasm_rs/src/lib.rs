use wasm_bindgen::prelude::*;
use js_sys::match_;

extern crate js_sys;

#[wasm_bindgen]
extern {
    #[wasm_bindgen(js_namespace = console)]
    fn log(s: &str);
}

#[wasm_bindgen(start)]
pub fn main() {

    let window = match web_sys::window() {
        Some(x) => x,
        None    => {
            log("no global `window` exists");
            return;
        }
    };

    let document = match window.document() {
        Some(x) => x,
        None    => {
            log("should have a document on window");
            return;
        }
    };

    let body = match document.body() {
        Some(x) => x,
        None    => {
            log("document should have a body");
            return;
        }
    };

    let element = match document.create_element("p") {
        Some(x) => x,
        None => {
            log("failed creating element `p` inside document");
            return;
        }
    };
    val.set_inner_html("Hello from Rust!");
    body.append_child(&val);
    log("Hello World!!!");
}