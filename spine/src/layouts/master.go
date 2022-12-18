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
package layouts

import (
	"path"
	"spine/src/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func StartMaster() {
	// Creating application and the master window.
	spine := app.New()
	masterWindow := spine.NewWindow("Spine")
	masterWindow.SetMaster()
	masterWindow.Resize(fyne.NewSize(800, 500))

	inputBox := widget.NewMultiLineEntry()
	inputBox.SetPlaceHolder("Simple, Lightweight, Fast...")

	inputBox.OnChanged = func(text string) {
	}

	openFile := fyne.NewMenuItem("Open File", func() {
		dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
			openErr := utils.OpenFile(inputBox, file)
			if openErr != nil {
				dialog.ShowError(err, masterWindow)
			}
			masterWindow.SetTitle(path.Base(file.URI().Path()) + "- Spine")
		}, masterWindow).Show()
	})

	saveFile := fyne.NewMenuItem("Save File", func() {
		utils.SaveFile(inputBox, masterWindow)
	})

	// Adding the Menu Items to the Master Window.
	masterWindow.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File", openFile, saveFile),
	))
	// Setting the content and rendering to the screen.
	masterWindow.SetContent(inputBox)
	masterWindow.ShowAndRun()
}
