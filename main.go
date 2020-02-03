package main

import (
	"fmt"
	"github.com/pkg/term/termios"
	"os"
	"syscall"
)

func enableRaw() {

	var attr syscall.Termios

	if err := termios.Tcgetattr(os.Stdin.Fd(), &attr); err != nil {
	}

	//disable echo
	attr.Lflag &^= (syscall.ECHO | syscall.ICANON)

	if err := termios.Tcsetattr(os.Stdin.Fd(), termios.TCSANOW, &attr); err != nil {
	}

}

type Arrow int

const (
	UP    Arrow = 0
	DOWN  Arrow = 1
	RIGHT Arrow = 2
	LEFT  Arrow = 3
	NONE  Arrow = 4
)

func is_arrow(input []byte) (bool, Arrow) {

	if len(input) != 3 {
		return false, NONE
	}

	if input[0] == '\033' && input[1] == '[' {
		switch input[2] {

		case 'A':
			{
				return true, UP
			}
		case 'B':
			{
				return true, DOWN
			}
		case 'C':
			{
				return true, RIGHT
			}
		case 'D':
			{
				return true, LEFT
			}
		default:
			{
				return false, NONE
			}
		}
	}
	return false, NONE

}

func is_delete(input []byte) bool {
	if len(input) != 1 {
		return false
	}

	if input[0] == 127 {
		return true
	}

	return false

}

func is_enter(input []byte) bool {
	if len(input) != 1 {
		return false
	}

	if input[0] == 10 {
		return true
	}

	return false

}

func SetForgroundColorWithRBG(r, g, b int) {
	fmt.Printf("\033[38;2;%d;%d;%dm", r, g, b)
}

func SetBackgroundColorWithRBG(r, g, b int) {
	fmt.Printf("\033[48;2;%d;%d;%dm", r, g, b)
}

type Cursor struct {
	x int
}

func MoveCursorBeginningOfPrev() {
	fmt.Printf("\033[F")
}

func DeleteCurrentLine() {
	fmt.Printf("\033[2K")
}

func RestoreCursorPos() {
	fmt.Printf("\033[u")
}
func SaveCursorPos() {
	fmt.Printf("\033[s")
}

func PreRender() {
	SaveCursorPos()
	DeleteCurrentLine()
	MoveCursorBeginningOfPrev()
	fmt.Printf("\033[1B")
	fmt.Print("Enter text: ")
}

func PostRender() {
	RestoreCursorPos()
}

func CursorRight() {
	fmt.Printf("\033[1C")

}

func CursorLeft() {
	fmt.Printf("\033[1D")

}

func Render(line_buffer []byte) {
	PreRender()

	fmt.Printf("%s", line_buffer)

	PostRender()
}

func DeleteRender(line_buffer *[]byte, cursor *Cursor, read_data []byte) {
	PreRender()

	*line_buffer = append((*line_buffer)[:cursor.x-1], (*line_buffer)[cursor.x:]...)
	fmt.Printf("%s", *line_buffer)
	cursor.x = cursor.x - 1

	PostRender()
	CursorLeft()
}

func InsertRender(line_buffer *[]byte, cursor *Cursor, read_data []byte) {
	PreRender()

	(*line_buffer) = append((*line_buffer)[:cursor.x], append(read_data, (*line_buffer)[cursor.x:]...)...)
	fmt.Printf("%s", *line_buffer)
	cursor.x = cursor.x + 1

	PostRender()
	CursorRight()
}

func main() {

	cursor := Cursor{x: 0}

	enableRaw()
	var b []byte = make([]byte, 100)
	line_buffer := []byte{}
	line_buffer_history := [][]byte{}
	history_index := 0
	//SetForgroundColorWithRBG(1, 240, 22)
	//SetBackgroundColorWithRBG(233, 1, 1)
	fmt.Print("Enter text: ")
	for {

		byte_count, _ := os.Stdin.Read(b)
		read_data := b[:byte_count]
		if ok, arrow := is_arrow(read_data); ok {
			switch arrow {
			case UP:
				{
					if len(line_buffer_history) > 0 {
						//fmt.Printf("\033[1A")
						line_buffer = line_buffer_history[history_index]
						if history_index < len(line_buffer_history)-1 {
							history_index++
						}
						Render(line_buffer)
					}
				}
			case DOWN:
				{
					if len(line_buffer_history) > 0 {
						line_buffer = line_buffer_history[history_index]
						if history_index > 0 {
							history_index--
						}
						Render(line_buffer)
					}
				}
			case RIGHT:
				{
					if cursor.x < len(line_buffer) {
						CursorRight()
						cursor.x = cursor.x + 1
					}
				}
			case LEFT:
				{
					if cursor.x > 0 {
						CursorLeft()
						cursor.x = cursor.x - 1
					}
				}

			}
		} else if is_enter(read_data) {
			line_buffer_history = append([][]byte{line_buffer}, line_buffer_history...)
			line_buffer = []byte{}
			fmt.Print("\n")
			fmt.Print("Enter text: ")
			cursor.x = 0

		} else if is_delete(read_data) {
			if len(line_buffer) > 0 && cursor.x > 0 {
				DeleteRender(&line_buffer, &cursor, read_data)
			}

		} else {
			InsertRender(&line_buffer, &cursor, read_data)
		}
	}
}
