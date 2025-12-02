use std::io;
fn main() {
    let mut input_text = String::new();
    io::stdin()
        .read_line(&mut input_text)
        .expect("Failed to read line");
    let number: i32 = input_text
        .trim()
        .parse()
        .expect("Invalid input: Please enter a valid number");
    if number % 2 == 0 {
        println!("Yes");
    } else {
        println!("No");
    }
}
