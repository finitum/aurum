use std::net::TcpStream;
use std::io::Write;
use ureq;

fn main() {
    let mut tcps = TcpStream::connect("0.0.0.0:8000").unwrap();

    tcps.write(b"test").expect("test");

    // let res = ureq::get("http://0.0.0.0:8000").call();

    // let client = Client::new();

//     let req = Request::post("localhost:8042/login")
//         .body(r#"
// {
//     "username": "jonay2000",
//     "password": "test",
// }
//         "#)
//         .expect("couldn't add body");
//
//     /// GET /test HTTP/1.0
//     /// Host: google.com

    // req.
    // tcps.write();
}
