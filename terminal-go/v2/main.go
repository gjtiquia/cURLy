package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/pkg/errors"
)

func main() {
	// logging setup
	err, logPanicAndCloseFile := InitLogFile("log.txt")
	if err != nil {
		log.Panicf("%+v", errors.WithStack(err))
	}
	defer logPanicAndCloseFile()

	// tcell setup
	s, err, finalizeScreen := InitTCellScreen()
	if err != nil {
		log.Panicf("%+v", err)
	}
	defer finalizeScreen()

	// Set default text style
	defStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	s.SetStyle(defStyle)

	// Clear screen
	s.Clear()

	s.Put(0, 0, "H", defStyle)
	s.Put(1, 0, "i", defStyle)
	s.Put(2, 0, "!", defStyle)

	s.PutStr(0, 1, "Hello World!")

	for {
		// Update screen
		s.Show()

		// Poll event (can be used in select statement as well)
		ev := <-s.EventQ()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			}
		}
	}
}

func InitLogFile(filename string) (err error, logPanicAndCloseFile func()) {
	// truncate means delete contents on open, create if doesnt exist, write-only
	const fileFlags = os.O_TRUNC | os.O_CREATE | os.O_WRONLY

	// read = 4, write = 2, execute = 1; 6 = 4 + 2 (read write); 0 = octal; 666 = owner/group/others
	const filePerm = 0666

	file, err := os.OpenFile(filename, fileFlags, filePerm)
	if err != nil {
		return err, nil
	}

	log.SetOutput(file)
	logPanicAndCloseFile = func() {
		defer file.Close()
		if r := recover(); r != nil {
			log.Println("logging panic before file close")
			log.Panicf("%+v", r)
		}
	}
	return nil, logPanicAndCloseFile
}

func InitTCellScreen() (s tcell.Screen, err error, finalizeScreen func()) {
	s, err = tcell.NewScreen()
	if err != nil {
		return nil, err, nil
	}
	if err = s.Init(); err != nil {
		return nil, err, nil
	}

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
