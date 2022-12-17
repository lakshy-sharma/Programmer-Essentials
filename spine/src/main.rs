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
