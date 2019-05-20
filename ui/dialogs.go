package main

import (
	ui "github.com/VladimirMarkelov/clui"
	term "github.com/nsf/termbox-go"
	"strings"
)

// OnClick sets the callback that is called when one clicks button
// with mouse or pressing space on keyboard while the button is active
func (b Button) OnClick(fn func(ui.Event)) {
	b.onClick = fn
}
type QuestionDialog struct {
	View    *ui.Window
	result  int
	onClose func()
}
type Button struct {
	ui.BaseControl
	shadowColor term.Attribute
	bgActive    term.Attribute
	pressed     int32
	onClick     func(ui.Event)
}

func CreateQuestionDialog(title, question string, buttons []string, defaultButton int) *QuestionDialog {
	dlg := new(QuestionDialog)
	if len(buttons) == 0 {
		buttons = []string{"OK"}
	}

	cw, ch := term.Size()
	ui.DrawFrame(0,0,cw,40,ui.BorderThin)
	lines := strings.Split(question, "\n")
	maxLen :=0
	for _,v := range lines{
		thisLen := len(v)*2
		if thisLen>maxLen{
			maxLen = thisLen
		}
	}
	dlg.View = ui.AddWindow(cw/2-12, ch/2-8, maxLen-25, 3, title)

	ui.WindowManager().BeginUpdate()
	defer ui.WindowManager().EndUpdate()
	dlg.View.SetConstraints(30, 3)
	dlg.View.SetModal(true)
	dlg.View.SetPack(ui.Vertical)
	currentPosX,currentPosY := dlg.View.Pos()
	dlg.View.SetPos(currentPosX-20, currentPosY)
	ui.CreateFrame(dlg.View, 1, 1, ui.BorderNone, ui.Fixed)
	dlg.View.SetAlign(ui.AlignCenter)
	fbtn := ui.CreateFrame(dlg.View, 1, 1,ui.BorderNone, 1)
	fbtn.SetAlign(ui.AlignCenter)
	ui.CreateFrame(fbtn, 1, 1, ui.BorderNone, ui.AutoSize)

	lb := ui.CreateLabel(fbtn, 10, 3, question, 1)

	lb.SetMultiline(true)
	ui.CreateFrame(fbtn, 1, 1,ui.BorderNone, ui.Fixed)

	ui.CreateFrame(dlg.View, 1, 1, ui.BorderNone, ui.Fixed)
	frm1 := ui.CreateFrame(dlg.View, 16, 4, ui.BorderNone, ui.Fixed)
	ui.CreateFrame(frm1, 1, 1, ui.BorderNone, 1)

	bText := buttons[0]
	btn1 := ui.CreateButton(frm1, ui.AutoSize, ui.AutoSize, bText, ui.Fixed)
	btn1.OnClick(func(ev ui.Event) {
		dlg.result = ui.DialogButton1

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
			dlg.result = ui.DialogButton2
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
			dlg.result = ui.DialogButton3
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
		if dlg.result == ui.DialogAlive {
			dlg.result = ui.DialogClosed
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
