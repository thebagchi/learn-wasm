use wasm_bindgen::prelude::*;

extern crate js_sys;
mod router;

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
        Err(_x) => {
            log("failed creating element `p` inside document");
            return;
        }
        Ok(x)   => x,
    };
    element.set_inner_html("Hello from Rust!");
    body.append_child(&element);
    log("Hello World!!!");
}