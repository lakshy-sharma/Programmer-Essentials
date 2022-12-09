/*
The main file is responsible for controlling the flow of the program.
This file is responsible for starting the editor when the program is started.
*/

// Importing the modules.
mod editor;
mod terminal;
mod document;
mod row;

// Importing the functions.
use editor::Editor;
pub use row::Row;
pub use terminal::Terminal;
pub use editor::Position;
pub use document::Document;

fn main() {
    Editor::default().run();
}
