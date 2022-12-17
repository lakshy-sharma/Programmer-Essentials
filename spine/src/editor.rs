/*
Copyright © [2022] [Lakshy Sharma] <lakshy.sharma@protonmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
use crate::Document;
use crate::Terminal;
use crate::Row;
use std::env;
use termion::event::Key;
use termion::color;
use std::time::Duration;
use std::time::Instant;

const STATUS_FG_COLOR: color::Rgb = color::Rgb(63,63,63);
const STATUS_BG_COLOR: color::Rgb = color::Rgb(239,239,239);
const VERSION: &str = env!("CARGO_PKG_VERSION");

// This structure holds all the required variables for controlling the editor flow.
pub struct Editor {
    should_quit: bool,                      // Controls whether the editor  should keep running or quit.
    terminal: Terminal,                     // Used to denote a terminal object which is responsible for controlling the low level terminal control stuff.
    document: Document,                     // Denotes a document type object which is used for controlling the flow of a document being viewed by the user.
    status_message: StatusMessage,          // Denotes the status message appearing below the editor.
    cursor_position: Position,              // Controls the position of the cursor inside the editor.
    offset: Position,                       // Controls the 
}

#[derive(Default)]
pub struct Position {
    pub x: usize,
    pub y: usize,
}

struct StatusMessage {
    text: String,
    time: Instant,
}

impl StatusMessage {
    fn from(message: String) -> Self {
        Self {
            time: Instant::now(),
            text: message,
        }
    }
}

impl Editor {

    pub fn default() -> Self {
        // This function starts the editor in default mode and sets the default values for variables in the Editor struct given above.

        let args: Vec<String> = env::args().collect();
        let mut initial_status = String::from("HELP: Ctrl-S = save | Ctrl-Q = quit");

        // Check if a document has been provided or display default page.
        let document = if args.len() > 1 {
            let file_name = &args[1];
            let doc = Document::open(&file_name);
            if doc.is_ok() {
                doc.unwrap()
            } else {
                initial_status = format!("Err: Could not open the file: {}", file_name);
                Document::default()
            }
        }
        else {
            Document::default()
        };

        Self{ 
            should_quit: false,
            terminal: Terminal::default().expect("Failed to initialize terminal."),
            cursor_position: Position::default(),
            document,
            offset: Position::default(),
            status_message: StatusMessage::from(initial_status),
        }
    }

    pub fn run(&mut self) {
        // Inifinte loop to read keypresses.
        loop {
            if let Err(error) = self.refresh_screen() {
                die(error);
            }
            if self.should_quit {
                break;
            }
            if let Err(error) = self.process_keypress() {
                die(error);
            }
        }
    }
    
    fn refresh_screen(&self) -> Result<(), std::io::Error> {
        // To avoid flicker effect we hide cursor before we refresh the screen.
        Terminal::cursor_hide();
        Terminal::cursor_position(&Position::default());

        if self.should_quit {
            Terminal::clear_screen();
        }
        else {
            self.draw_rows();
            self.draw_status_bar();
            self.draw_message_bar();
            Terminal::cursor_position(&Position {
                x: self.cursor_position.x.saturating_sub(self.offset.x),
                y: self.cursor_position.y.saturating_sub(self.offset.y),
            });
        }
        Terminal::cursor_show();
        Terminal::flush()
    }

    fn draw_welcome_message(&self) {
        // This function creates a entry message for the editor.

        let mut welcome_message = format!("Spine - Simple, Lightweight, Fast...");
        let mut version_message = format!("Version - {}", VERSION);

        let screen_width = self.terminal.size().width as usize;
        let welcome_len = welcome_message.len();
        let welcome_padding = screen_width.saturating_sub(welcome_len) / 2;
        let welcome_spaces = " ".repeat(welcome_padding.saturating_sub(1));

        let version_len = version_message.len();
        let version_padding = screen_width.saturating_sub(version_len) / 2;
        let version_spaces = " ".repeat(version_padding.saturating_sub(1));

        welcome_message = format!("~{}{}", welcome_spaces, welcome_message);
        welcome_message.truncate(screen_width);

        version_message = format!("~{}{}", version_spaces, version_message);
        version_message.truncate(screen_width);

        println!("{}\r", welcome_message);
        println!("{}\r",version_message)
    }

    pub fn draw_row(&self, row: &Row) {
        // This function prints rows 
     
        let width = self.terminal.size().width as usize;
        let start = self.offset.x;
        let end = self.offset.x + width;
        let row = row.render(start,end);
        println!("{}\r", row)
    }

    fn draw_rows(&self){
        let height = self.terminal.size().height;
        for terminal_row in 0..height {
            Terminal::clear_current_line();
            if let Some(row) = self.document.row(terminal_row as usize + self.offset.y ) {
                self.draw_row(row);
            } 
            else if self.document.is_empty() && terminal_row == height / 3 {
                self.draw_welcome_message();
            }
            else {
                println!("~\r")
            }
        }
    }

    fn draw_status_bar(&self) {
        let mut status;
        let width = self.terminal.size().width as usize;
        let mut file_name = "[No Name]".to_string();
        if let Some(name) = &self.document.file_name {
            file_name = name.clone();
            file_name.truncate(20);
        }
        status = format!("{} - {} lines", file_name, self.document.len());

        // Show current cursor location in the file on the status bar.
        let line_indicator = format!(
            "{}/{}", 
            self.cursor_position.y.saturating_add(1), 
            self.cursor_position.x.saturating_add(1)
        );
        let len = status.len() + line_indicator.len();
        if width > len {
            status.push_str(&" ".repeat(width - len));
        }
        status = format!("{}{}", status, line_indicator);
        status.truncate(width);
        Terminal::set_bg_color(STATUS_BG_COLOR);
        Terminal::set_fg_color(STATUS_FG_COLOR);
        println!("{}\r", status);
        Terminal::reset_fg_color();
        Terminal::reset_bg_color();
    }

    fn draw_message_bar(&self) {
        Terminal::clear_current_line();
        let message = &self.status_message;
        if Instant::now() - message.time < Duration::new(5,0) {
            let mut text = message.text.clone();
            text.truncate(self.terminal.size().width as usize);
            print!("{}", text);
        }
    }

    fn process_keypress(&mut self) -> Result<(), std::io::Error> {
        // Read the pressed key from external infinte loop function.
        let pressed_key = Terminal::read_key()?;
 
        // Quit if the pressed key is Ctrl+Q
        match pressed_key {
            Key::Ctrl('q') => self.should_quit = true,
            Key::Ctrl('s') => {
                if self.document.file_name.is_none() {
                    self.document.file_name = Some(self.prompt("Save as: ")?);
                }
                if self.document.save().is_ok() {
                    self.status_message = StatusMessage::from("File saved successfully.".to_string());                    
                }
                else {
                    self.status_message = StatusMessage::from("Error writing file!".to_string());
                }
            }
            Key::Char(c) => {
                self.document.insert(&self.cursor_position,c);
                self.move_cursor(Key::Right);
            }
            Key::Delete => self.document.delete(&self.cursor_position),
            Key::Backspace => {
                if self.cursor_position.x > 0 || self.cursor_position.y > 0 {
                    self.move_cursor(Key::Left);
                    self.document.delete(&self.cursor_position);
                }
            }
            Key::Up 
            | Key::Down 
            | Key::Left 
            | Key::Right 
            | Key::PageUp 
            | Key::PageDown 
            | Key::End
            | Key::Home => self.move_cursor(pressed_key),
            _ => (),
        }
        self.scroll();
        Ok(())
    }

    fn scroll(&mut self) {
        let Position {x , y} = self.cursor_position;
        let width = self.terminal.size().width as usize;
        let height = self.terminal.size().height as usize;
        let mut offset = &mut self.offset;
        if y < offset.y {
            offset.y = y;
        }
        else if y>=offset.y.saturating_add(height) {
            offset.y = y.saturating_sub(height).saturating_add(1);
        }

        if x < offset.x {
            offset.x = x;
        }
        else if x >= offset.x.saturating_add(width) {
            offset.x = x.saturating_sub(width).saturating_add(1);
        }
    }

    fn move_cursor(&mut self, key: Key) {
        let Position { mut y, mut x} = self.cursor_position;
        let terminal_height = self.terminal.size().height as usize;
        let size = self.terminal.size();
        let height = self.document.len();
        let mut width = if let Some(row) = self.document.row(y) {
            row.len()
        }
        else {
            0
        };

        match key {
            Key::Up => y = y.saturating_sub(1),
            Key::Down => {
                // If the cursor is within document height limit then perform a saturating add.
                if y < height {
                    y = y.saturating_add(1);
                }
            },
            Key::Left => {
                // When user presses left key after reaching start of a row go to an upper row.
                if x > 0 {
                    x -= 1;
                } else if y >0 {
                    y -= 1;
                    if let Some(row) = self.document.row(y) {
                        x = row.len();
                    } else {
                        x = 0;
                    }
                }
            }
            Key::Right => {
                // If the cursor is within document width limit then perform a saturating add.
                if x < width {
                    x +=1;
                } else if y < height {
                    y += 1;
                    x = 0;
                }
            },
            Key::PageUp => {
                // Scroll up by the height of the terminal if the length of y is greater than the terminal height.
                y = if y > terminal_height {
                    y - terminal_height
                } else {
                    0
                }
            }
            Key::PageDown => {
                y = if y.saturating_add(terminal_height) < height {
                    y + terminal_height as usize
                } else {
                    height
                }
            }
            Key::Home => x = 0,
            Key::End => x = width,
            _ => (),
        }
        // Calculating the row length to snap the cursor movement inside the limits.
        width = if let Some(row) = self.document.row(y) {
            row.len()
        } else {
            0
        };

        // Limit value of x within the row width limits.
        if x > width {
            x = width;
        }

        self.cursor_position = Position { x, y }
    }

    fn prompt(&mut self, prompt: &str) -> Result<String, std::io::Error> {            
       let mut result = String::new();            
        loop {            
            self.status_message = StatusMessage::from(format!("{}{}", prompt, result));            
            self.refresh_screen()?;            
            if let Key::Char(c) = Terminal::read_key()? {            
                if c == '\n' {            
                    self.status_message = StatusMessage::from(String::new());            
                    break;            
                }            
                if !c.is_control() {            
                    result.push(c);            
                }            
            }            
        }            
        Ok(result)            
    }
}

fn die(error: std::io::Error) {
    // This function is used when we need to quit the editor while facing any error.
    Terminal::clear_screen();
    panic!("{}",error);
}
