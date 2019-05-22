package ui

import (
	"../utils"
	"fmt"
	ui "github.com/VladimirMarkelov/clui"
	term "github.com/nsf/termbox-go"
	_ "github.com/sirupsen/logrus"
	_ "log"
	"os"
	"os/exec"
	//"strconv"
	"strings"
)

// OnClick sets the callback that is called when one clicks button
// with mouse or pressing space on keyboard while the button is active
func (b Button) OnClick(fn func(ui.Event)) {
	b.onClick = fn
}

type QuestionDialog struct {
	View    *ui.Window
	Result  int
	onClose func()
}
type Button struct {
	ui.BaseControl
	shadowColor term.Attribute
	bgActive    term.Attribute
	pressed     int32
	onClick     func(ui.Event)
}

var OptionsRight *ui.Frame
var HelpField *ui.Frame
var newValuesHolder map[string]string

func ConfigurationEditor(configMap map[string]string, configFile string) int {

	tmpMainMenuString := ""
	for k, _ := range configMap {
		tmpMainMenuString += k + "||"
	}
	mainListBoxOptions := strings.Split(tmpMainMenuString, "||")
	cw, ch := term.Size()

	view := ui.AddWindow(cw/2-45, ch/2-12, 1, 7, "ZCS LMDB Testing Configuration")
	oldValuesHolder := make(map[string]string)
	frmLeft := ui.CreateFrame(view, 8, 4, ui.BorderNone, 0)
	frmLeft.SetPack(ui.Vertical)
	//frmLeft.SetGaps(ui.KeepValue, 1)
	frmLeft.SetPaddings(1, 1)

	logBox := ui.CreateListBox(frmLeft, 20, 15, ui.Fixed)
	for _, v := range mainListBoxOptions {
		logBox.AddItem(v)
	}
	var fieldGroup *ui.Frame

	frmRight := ui.CreateFrame(view, 70, 1, ui.BorderNone, 1)
	HelpField = ui.CreateFrame(frmRight, 0, 0, ui.BorderNone, ui.Horizontal)

	OptionsRight = frmRight
	firstLoad := true
	isValid := true
	logBox.OnSelectItem(func(event ui.Event) {

		//before moving onto the next page, we need to validate that what is currently
		//on this page, is valid.

		//UpdateConfigWithNewValue(config string, oldval string, newval string)
		//newValuesHolder[label] = value
		//oldValuesHolder[label] = value
		var invalidKey string
		if !firstLoad{
			for k, v := range newValuesHolder {
				//isValidValue, notice := utils.ValidateConfigKey(k, v)
				isValid = utils.ValidateConfigKey(configFile, k, v)
				if !isValid{
					invalidKey = k
					break
				}
			}

		}
		firstLoad = false
		//if !isValid{
		//	//canBreak := false
		//	dialog := ui.CreateAlertDialog("Invalid Configuration Item", "The value set for "+invalidKey+" is not valid.", "OK")
		//	dialog.OnClose(func() {
		//		return
		//	})
		//}else{
		if isValid{
			for k, v := range newValuesHolder {
				//isValidValue, notice := utils.ValidateConfigKey(k, v)
				oldValue := k + "=" + oldValuesHolder[k]
				newValue := k + "=" + v
				utils.UpdateConfigWithNewValue(configFile, oldValue, newValue)
			}
			configMap = utils.BuildMenuFromConfig(configFile)
			tmpMainMenuString := ""
			for k, _ := range configMap {
				tmpMainMenuString += k + "||"
			}
			mainListBoxOptions = strings.Split(tmpMainMenuString, "||")
			//ui.CreateLabel(fieldGroup, ui.AutoSize, ui.AutoSize, strings.Replace(labelBuffer+label,"_"," ",-1)+":", ui.Fixed)
			//thisEditField := ui.CreateEditField(fieldGroup, ui.AutoSize, value, 0)

			newValuesHolder = make(map[string]string)
			OptionsRight.Destroy()
			subItems := strings.Split(configMap[logBox.SelectedItemText()], "||")
			frmRight := ui.CreateFrame(view, 1, 1, ui.BorderNone, ui.Vertical)
			OptionsRight = frmRight
			frmRight.SetPack(ui.Vertical)
			frmRight.SetGaps(1, 0)
			maxDeltalen := 0
			for _, v := range subItems {
				if len(v) > maxDeltalen {
					maxDeltalen = len(v)
				}
			}
			for _, v := range subItems {
				if strings.Contains(v, "=") {
					fieldGroup = ui.CreateFrame(frmRight, 70, 1, ui.BorderNone, ui.Horizontal)
					fieldGroup.SetPack(ui.Horizontal)
					fieldGroup.SetGaps(1, 1)
					label := strings.Split(v, "=")[0]

					spacingBuffer := maxDeltalen / 2
					if spacingBuffer > 20 {
						spacingBuffer = 20
					}
					labelTextLen := len(label)
					labelBuffer := ""
					for {
						if labelTextLen <= spacingBuffer {
							labelBuffer += " "
							labelTextLen++
						} else {
							break
						}
					}

					value := strings.Split(v, "=")[1]
					ui.CreateLabel(fieldGroup, ui.AutoSize, ui.AutoSize, strings.Replace(labelBuffer+label, "_", " ", -1)+":", ui.Fixed)
					thisEditField := ui.CreateEditField(fieldGroup, ui.AutoSize, value, 1)
					ui.CreateFrame(frmRight, 70, 1, ui.BorderNone, ui.Horizontal)

					/////////
					//store the current values in a map
					//this will get updated on click
					newValuesHolder[label] = value
					oldValuesHolder[label] = value

					thisEditField.OnChange(func(event ui.Event) {
						newValuesHolder[label] = thisEditField.Title()
					})
					/////////

					thisEditField.OnActive(func(active bool) {
						HelpField.Destroy()
						HelpField = ui.CreateFrame(frmRight, 70, 1, ui.BorderThin, ui.Horizontal)
						HelpField.SetPack(ui.Horizontal)
						HelpField.SetGaps(1, 1)
						lineCount, lineArray := utils.GetDescriptionTextForUpdate(configFile, label)
						TextHelpField := ui.CreateTextReader(HelpField, ui.AutoSize, lineCount, ui.Vertical)
						TextHelpField.SetBackColor(ui.ColorBlack)
						TextHelpField.SetTextColor(ui.ColorWhiteBold)
						TextHelpField.SetLineCount(lineCount)
						TextHelpField.Draw()
						TextHelpField.OnDrawLine(func(ind int) string {
							return fmt.Sprint(lineArray[ind])
						})
					})

				}
			}
		}else{
			dialog := ui.CreateAlertDialog("Invalid Configuration Item", "The value set for "+invalidKey+" is not valid.", "OK")
			dialog.OnClose(func() {
				return
			})
		}






	})
	frmRightHelp := ui.CreateFrame(frmRight, 8, 1, ui.BorderNone, ui.Fixed)
	frmRightHelp.SetPaddings(1, 1)
	frmRightHelp.SetGaps(1, ui.KeepValue)

	frmEdit := ui.CreateFrame(frmLeft, 8, 1, ui.BorderNone, ui.Fixed)
	frmEdit.SetPaddings(1, 1)
	frmEdit.SetGaps(1, ui.KeepValue)
	saveConfigBtn := ui.CreateButton(frmEdit, ui.AutoSize, 4, "Save Config", ui.Fixed)
	saveConfigBtn.SetBackColor(ui.ColorCyan)
	//Test Button is here..
	ui.CreateButton(frmEdit, ui.AutoSize, 4, "Run Test", ui.Fixed)
	ui.CreateFrame(frmEdit, 1, 1, ui.BorderNone, 1)
	btnQuit := ui.CreateButton(frmEdit, ui.AutoSize, 4, "Quit", ui.Fixed)
	btnQuit.SetBackColor(ui.ColorRed)
	ui.ActivateControl(view, logBox)
	userQuit := false
	btnQuit.OnClick(func(ev ui.Event) {
		dialog := ui.CreateConfirmationDialog("Are you sure?", "Do you wish to exit without running tests?", []string{"Yes", "No"}, 2)
		dialog.OnClose(func() {
			result := dialog.Result()
			if result == 1 {
				//RefreshScreen()
				userQuit = true
				Cleanup()
				event := ui.Event{Type: ui.EventCloseWindow, X: 1}
				ui.ProcessEvent(event)
				return
			}
		})
	})

	saveConfigBtn.OnClick(func(event ui.Event) {
		//UpdateConfigWithNewValue(config string, oldval string, newval string)
		//newValuesHolder[label] = value
		//oldValuesHolder[label] = value
		isValid := true
		for k, v := range newValuesHolder {
			//isValidValue, notice := utils.ValidateConfigKey(k, v)
			if !isValid {
				break
			}
			isValid = utils.ValidateConfigKey(configFile, k, v)
			if !isValid {
				dialog := ui.CreateAlertDialog("Invalid Configuration Item", "The value set for "+k+"is not valid.", "OK")
				dialog.OnClose(func() {
					return
				})
			} else {
				oldValue := k + "=" + oldValuesHolder[k]
				newValue := k + "=" + v
				utils.UpdateConfigWithNewValue(configFile, oldValue, newValue)
			}
		}

		configMap = utils.BuildMenuFromConfig(configFile)
		tmpMainMenuString := ""
		for k, _ := range configMap {
			tmpMainMenuString += k + "||"
		}
		mainListBoxOptions = strings.Split(tmpMainMenuString, "||")
		//ui.CreateLabel(fieldGroup, ui.AutoSize, ui.AutoSize, strings.Replace(labelBuffer+label,"_"," ",-1)+":", ui.Fixed)
		//thisEditField := ui.CreateEditField(fieldGroup, ui.AutoSize, value, 0)
	})
	btnQuit.OnClick(func(ev ui.Event) {

	})
	if userQuit {
		return 99
	}
	return 0
}

func CreateQuestionDialog(title, question string, buttons []string, defaultButton int) *QuestionDialog {
	dlg := new(QuestionDialog)
	if len(buttons) == 0 {
		buttons = []string{"OK"}
	}

	cw, ch := term.Size()
	ui.DrawFrame(0, 0, cw, 40, ui.BorderThin)
	lines := strings.Split(question, "\n")
	maxLen := 0
	for _, v := range lines {
		thisLen := len(v) * 2
		if thisLen > maxLen {
			maxLen = thisLen
		}
	}
	dlg.View = ui.AddWindow(cw/2-12, ch/2-8, maxLen-25, 3, title)

	ui.WindowManager().BeginUpdate()
	defer ui.WindowManager().EndUpdate()
	dlg.View.SetConstraints(30, 3)
	dlg.View.SetModal(true)
	dlg.View.SetPack(ui.Vertical)
	currentPosX, currentPosY := dlg.View.Pos()
	dlg.View.SetPos(currentPosX-20, currentPosY)
	ui.CreateFrame(dlg.View, 1, 1, ui.BorderNone, ui.Fixed)
	dlg.View.SetAlign(ui.AlignCenter)
	fbtn := ui.CreateFrame(dlg.View, 1, 1, ui.BorderNone, 1)
	fbtn.SetAlign(ui.AlignCenter)
	ui.CreateFrame(fbtn, 1, 1, ui.BorderNone, ui.AutoSize)

	lb := ui.CreateLabel(fbtn, 10, 3, question, 1)

	lb.SetMultiline(true)
	ui.CreateFrame(fbtn, 1, 1, ui.BorderNone, ui.Fixed)

	ui.CreateFrame(dlg.View, 1, 1, ui.BorderNone, ui.Fixed)
	frm1 := ui.CreateFrame(dlg.View, 16, 4, ui.BorderNone, ui.Fixed)
	ui.CreateFrame(frm1, 1, 1, ui.BorderNone, 1)

	bText := buttons[0]
	btn1 := ui.CreateButton(frm1, ui.AutoSize, ui.AutoSize, bText, ui.Fixed)
	btn1.OnClick(func(ev ui.Event) {
		dlg.Result = ui.DialogButton1
		ui.WindowManager().DestroyWindow(dlg.View)
		ui.WindowManager().BeginUpdate()
		closeFunc := dlg.onClose
		ui.WindowManager().EndUpdate()
		if closeFunc != nil {
			closeFunc()
		}
	})
	var btn2, btn3 *ui.Button

	if len(buttons) > 1 {
		ui.CreateFrame(frm1, 1, 1, ui.BorderNone, 1)
		btn2 = ui.CreateButton(frm1, ui.AutoSize, ui.AutoSize, buttons[1], ui.Fixed)
		btn2.OnClick(func(ev ui.Event) {
			dlg.Result = ui.DialogButton2
			ui.WindowManager().DestroyWindow(dlg.View)
			if dlg.onClose != nil {
				dlg.onClose()
			}
		})
	}
	if len(buttons) > 2 {
		ui.CreateFrame(frm1, 1, 1, ui.BorderNone, 1)
		btn3 = ui.CreateButton(frm1, ui.AutoSize, ui.AutoSize, buttons[2], ui.Fixed)
		btn3.OnClick(func(ev ui.Event) {
			dlg.Result = ui.DialogButton3
			ui.WindowManager().DestroyWindow(dlg.View)
			if dlg.onClose != nil {
				dlg.onClose()
			}
		})
	}

	ui.CreateFrame(frm1, 1, 1, ui.BorderNone, 1)

	if defaultButton == ui.DialogButton2 && len(buttons) > 1 {
		ui.ActivateControl(dlg.View, btn2)
	} else if defaultButton == ui.DialogButton3 && len(buttons) > 2 {
		ui.ActivateControl(dlg.View, btn3)
	} else {
		ui.ActivateControl(dlg.View, btn1)
	}

	dlg.View.OnClose(func(ev ui.Event) bool {
		if dlg.Result == ui.DialogAlive {
			dlg.Result = ui.DialogClosed
			if ev.X != 1 {
				ui.WindowManager().DestroyWindow(dlg.View)
			}
			if dlg.onClose != nil {
				dlg.onClose()
			}
		}
		return true
	})

	return dlg
}
func formatter(line string) {
	fmt.Println(line)
}
func Cleanup() {
	ui.DeinitLibrary()
	go ui.Stop()
	cmd := exec.Command("reset")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	output := cmd.Run()
	print(output)
	print("\033[H\033[2J")
}
