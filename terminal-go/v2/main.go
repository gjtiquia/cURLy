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
	file, err := InitLogFile("log.txt")
	if err != nil {
		log.Panicf("%+v", errors.WithStack(err))
	}
	defer file.Close()
	defer LogPanicBeforeFileClose()

	// TODO : remove
	// panic(errors.WithStack(errors.New("hi i am panicking")))

	// tcell setup
	s, err := tcell.NewScreen()
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	if err := s.Init(); err != nil {
		log.Printf("%+v", err)
		return
	}

	// Set default text style
	defStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	s.SetStyle(defStyle)

	// Clear screen
	s.Clear()

	s.Put(0, 0, "H", defStyle)
	s.Put(1, 0, "i", defStyle)
	s.Put(2, 0, "!", defStyle)

	s.PutStr(0, 1, "Hello World!")

	quit := func() {
		s.Fini()
		os.Exit(0)
	}
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
				quit()
			}
		}
	}
}

func InitLogFile(filename string) (*os.File, error) {
	// truncate means delete contents on open, create if doesnt exist, write-only
	const fileFlags = os.O_TRUNC | os.O_CREATE | os.O_WRONLY

	// read = 4, write = 2, execute = 1; 6 = 4 + 2 (read write); 0 = octal; 666 = owner/group/others
	const filePerm = 0666

	file, err := os.OpenFile(filename, fileFlags, filePerm)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	return file, nil
}

func LogPanicBeforeFileClose() {
	if r := recover(); r != nil {
		log.Println("logging panic before file close")
		log.Panicf("%+v", r)
	}
}
