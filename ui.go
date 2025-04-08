package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"go.bug.st/serial"
)

// idea is the main window will help with the setup of the port
// The other window will spawn when the port is open, will have a close button
// and will have a text area to display the logs

// SerialSettings represents the UI elements for serial port configuration
type SerialSettings struct {
	portSelect *widget.Select
	baudRate   *widget.Select
	baudEntry  *widget.Entry
	dataBits   *widget.Select
	stopBits   *widget.Select
	parityBits *widget.Select
	serialConf SerialConf
	connectBtn *widget.Button
	console    *widget.TextGrid
}

type FileUI struct {
	fileBtn      *widget.Button
	selectedFile *widget.Label
	fileURI      fyne.URI
}

func (f *FileUI) showFilePicker(w *fyne.Window) {
	dialog.ShowFileSave(func(uri fyne.URIWriteCloser, err error) {
		savedFile := "No file selected"
		if err != nil {
			dialog.ShowError(err, *w)
			return
		}
		if uri == nil {
			return
		}
		savedFile = uri.URI().Path()
		f.fileURI = uri.URI()
		f.selectedFile.SetText(savedFile)
	}, (*w))
}

func NewFileUI(w *fyne.Window) *FileUI {
	f := &FileUI{}
	f.fileBtn = widget.NewButton("File", func() { f.showFilePicker(w) })
	f.selectedFile = widget.NewLabel("No file selected")
	f.fileURI = nil
	return f
}

// NewSerialSettings creates and initializes the serial settings UI components
func NewSerialSettings(ports []string) *SerialSettings {
	s := &SerialSettings{}

	// Default values for serial configuration
	s.serialConf = SerialConf{
		BaudRate: 9600,
		DataBits: 8,
		StopBits: serial.OneStopBit,
		Parity:   serial.NoParity,
	}

	// Port selection
	s.portSelect = widget.NewSelect(ports, func(value string) {
		s.serialConf.Port = value
	})
	if len(ports) > 0 {
		s.portSelect.SetSelected(ports[0])
		s.serialConf.Port = ports[0]
	}

	// Custom baud rate entry
	s.baudEntry = widget.NewEntry()
	s.baudEntry.SetPlaceHolder("Enter custom baud rate")
	s.baudEntry.Disable() // Disabled by default
	s.baudEntry.OnChanged = func(value string) {
		if baud, err := strconv.Atoi(value); err == nil {
			s.serialConf.BaudRate = baud
		}
	}

	// Baud rate options
	baudRates := []string{"9600", "19200", "38400", "57600", "115200", "Custom"}
	s.baudRate = widget.NewSelect(baudRates, func(value string) {
		if value == "Custom" {
			s.baudEntry.Enable()
		} else {
			s.baudEntry.Disable()
			switch value {
			case "9600":
				s.serialConf.BaudRate = 9600
			case "19200":
				s.serialConf.BaudRate = 19200
			case "38400":
				s.serialConf.BaudRate = 38400
			case "57600":
				s.serialConf.BaudRate = 57600
			case "115200":
				s.serialConf.BaudRate = 115200
			}
		}
	})
	s.baudRate.SetSelected("9600")

	// Data bits options
	dataBitsOptions := []string{"5", "6", "7", "8"}
	s.dataBits = widget.NewSelect(dataBitsOptions, func(value string) {
		switch value {
		case "5":
			s.serialConf.DataBits = 5
		case "6":
			s.serialConf.DataBits = 6
		case "7":
			s.serialConf.DataBits = 7
		case "8":
			s.serialConf.DataBits = 8
		}
	})
	s.dataBits.SetSelected("8")

	// Stop bits options
	stopBitsOptions := []string{"1", "1.5", "2"}
	s.stopBits = widget.NewSelect(stopBitsOptions, func(value string) {
		switch value {
		case "1":
			s.serialConf.StopBits = serial.OneStopBit
		case "1.5":
			s.serialConf.StopBits = serial.OnePointFiveStopBits
		case "2":
			s.serialConf.StopBits = serial.TwoStopBits
		}
	})
	s.stopBits.SetSelected("1")

	// Parity options
	parityOptions := []string{"None", "Odd", "Even", "Mark", "Space"}
	s.parityBits = widget.NewSelect(parityOptions, func(value string) {
		switch value {
		case "None":
			s.serialConf.Parity = serial.NoParity
		case "Odd":
			s.serialConf.Parity = serial.OddParity
		case "Even":
			s.serialConf.Parity = serial.EvenParity
		case "Mark":
			s.serialConf.Parity = serial.MarkParity
		case "Space":
			s.serialConf.Parity = serial.SpaceParity
		}
	})
	s.parityBits.SetSelected("None")

	s.console = widget.NewTextGrid()

	// Connect button
	s.connectBtn = widget.NewButton("Connect", func() {
		go CreateLogUI(s.serialConf) // TODO: need to pass the serialConf to the logUI
	})

	return s
}

func CreateUI(ports []string) {
	a := app.New()
	w := a.NewWindow("Serial Monitor")
	fileUI := NewFileUI(&w)
	list := widget.NewList(
		func() int {
			return len(ports)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(ports[i])
		})

	// Create serial settings
	settings := NewSerialSettings(ports)

	// Create a container for all serial settings widgets
	settingsContainer := container.NewVBox(
		widget.NewLabel("Port:"),
		settings.portSelect,
		widget.NewLabel("Baud Rate:"),
		container.NewGridWithColumns(2, settings.baudRate, settings.baudEntry),
		widget.NewLabel("Data Bits:"),
		settings.dataBits,
		widget.NewLabel("Stop Bits:"),
		settings.stopBits,
		widget.NewLabel("Parity:"),
		settings.parityBits,
		fileUI.selectedFile,
		fileUI.fileBtn,
		settings.connectBtn,
	)

	w.SetContent(container.NewVBox(
		list,
		settingsContainer,
	))
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

func CreateLogUI(s SerialConf) {
	w := app.New().NewWindow("Serial Monitor")
	port, err := s.OpenPort()
	go func() {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		buffer, err := ReadPort(port)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		file, err := OpenFile()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		file.Write(buffer)
	}()

}
