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
					//fmt.Printf("\033[1A")
					line_buffer = line_buffer_history[history_index]
					history_index++
				}
			case DOWN:
				{
					//fmt.Printf("\033[1B")
					line_buffer = line_buffer_history[history_index]
					history_index++

				}
			case RIGHT:
				{
					if cursor.x < len(line_buffer) {
						fmt.Printf("\033[1C")
						cursor.x = cursor.x + 1
					}
				}
			case LEFT:
				{
					if cursor.x > 0 {
						fmt.Printf("\033[1D")
						cursor.x = cursor.x - 1
					}
				}

			}
		} else if is_enter(read_data) {
			line_buffer_history = append(line_buffer_history, line_buffer)
			line_buffer = []byte{}
			fmt.Print("\n")
			fmt.Print("Enter text: ")
			cursor.x = 0

		} else if is_delete(read_data) {
			if len(line_buffer) > 0 && cursor.x > 0 {
				fmt.Printf("\033[s")
				fmt.Printf("\033[2K")
				fmt.Printf("\033[F")
				fmt.Printf("\033[1B")
				fmt.Print("Enter text: ")

				//for i := 0; i < len(line_buffer)-cursor.x; i++ {
				//fmt.Printf("\033[1C")
				//}

				//for _, _ = range line_buffer {
				//fmt.Printf("\b")
				/*}*/
				line_buffer = append(line_buffer[:cursor.x-1], line_buffer[cursor.x:]...)
				//line_buffer = append(line_buffer[:cursor.x], append(read_data, line_buffer[cursor.x:]...)...)
				//line_buffer = append(line_buffer, read_data...)
				fmt.Printf("%s", line_buffer)
				cursor.x = cursor.x - 1

				fmt.Printf("\033[u")
				fmt.Printf("\033[1D")
			}

		} else {
			fmt.Printf("\033[s")
			fmt.Printf("\033[2K")
			fmt.Printf("\033[F")
			fmt.Printf("\033[1B")
			fmt.Print("Enter text: ")

			//for i := 0; i < len(line_buffer)-cursor.x; i++ {
			//fmt.Printf("\033[1C")
			//}

			//for _, _ = range line_buffer {
			//fmt.Printf("\b")
			/*}*/

			line_buffer = append(line_buffer[:cursor.x], append(read_data, line_buffer[cursor.x:]...)...)
			//line_buffer = append(line_buffer, read_data...)
			fmt.Printf("%s", line_buffer)
			cursor.x = cursor.x + 1

			fmt.Printf("\033[u")
			fmt.Printf("\033[1C")
		}
	}
}
