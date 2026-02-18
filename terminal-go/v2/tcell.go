package main

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func InitTCellScreen() (s tcell.Screen, err error, finalizeScreen func()) {
	s, err = tcell.NewScreen()
	if err != nil {
		return nil, err, nil
	}
	if err = s.Init(); err != nil {
		return nil, err, nil
	}

	// Set default text style
	defStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	s.SetStyle(defStyle)

	// Clear screen
	s.Clear()

	finalizeScreen = func() {
		// You have to catch panics in a defer, clean up, and re-raise them - otherwise your application can die without leaving any diagnostic trace.
		// https://github.com/gdamore/tcell/blob/main/TUTORIAL.md
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}

	return s, nil, finalizeScreen
}

func DrainTCellEvents(s tcell.Screen, inputBuffer []InputAction) (buffer []InputAction, isExit bool) {
	inputBuffer = inputBuffer[:0]
	for {
		// Update screen
		s.Show()

		select {
		case ev := <-s.EventQ():
			// Process event
			switch ev := ev.(type) {

			case *tcell.EventResize:
				s.Sync()

			case *tcell.EventKey:
				key, str := ev.Key(), ev.Str()
				// log.Printf("key event: %v, %v", key, str)

				if key == tcell.KeyEscape || key == tcell.KeyCtrlC {
					inputBuffer = inputBuffer[:0]
					inputBuffer = append(inputBuffer, Exit)
					return inputBuffer, true
				}

				switch {
				case key == tcell.KeyUp, str == "w":
					inputBuffer = append(inputBuffer, Up)
				case key == tcell.KeyDown, str == "s":
					inputBuffer = append(inputBuffer, Down)
				case key == tcell.KeyLeft, str == "a":
					inputBuffer = append(inputBuffer, Left)
				case key == tcell.KeyRight, str == "d":
					inputBuffer = append(inputBuffer, Right)
				case str == "r":
					inputBuffer = append(inputBuffer, Restart)
				}
			}

		default:
			// no more events to process, terminate loop and return
			return inputBuffer, false
		}
	}
}
