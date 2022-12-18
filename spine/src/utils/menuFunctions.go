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
package utils

import (
	"io/ioutil"
	"spine/src/global"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func OpenFile(inputBox *widget.Entry, file fyne.URIReadCloser) error {
	// Check if the user has provided a valid file path.
	if file != nil {
		// Set the current path.
		global.CurrentPath = file.URI().Path()
		// Read the file text.
		content, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		// Set text content in the input box.
		inputBox.SetText(string(content))
	}
	return nil
}
