package readline

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

var prompt string

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

func is_tab(input []byte) bool {
	if len(input) != 1 {
		return false
	}

	if input[0] == 9 {
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
	fmt.Print(prompt)
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

func MoveCursorEndOfLine(length int) {
	fmt.Printf("\033[%dG", length)

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
func Readline(_prompt string, callback func(input []byte)) {

	prompt = _prompt

	cursor := Cursor{x: 0}
	tree := prefix_tree{root: &prefix_node{children: map[byte]*prefix_node{}}}

	enableRaw()
	var b []byte = make([]byte, 100)
	line_buffer := []byte{}
	line_buffer_history := [][]byte{}
	history_index := -1
	//SetForgroundColorWithRBG(1, 240, 22)
	//SetBackgroundColorWithRBG(233, 1, 1)
	fmt.Print(prompt)

	cur_line_buffer := []byte{}

	for {

		byte_count, _ := os.Stdin.Read(b)
		read_data := b[:byte_count]
		if ok, arrow := is_arrow(read_data); ok {
			switch arrow {
			case UP:
				{
					if len(line_buffer_history) > 0 {
						//fmt.Printf("\033[1A")
						if history_index == -1 {
							cur_line_buffer = line_buffer
						}
						if history_index < len(line_buffer_history)-1 {
							history_index++
							line_buffer = line_buffer_history[history_index]
							Render(line_buffer)
							length := len(line_buffer) + len(prompt) + 1
							MoveCursorEndOfLine(length)
							cursor.x = len(line_buffer)
						}

					}
				}
			case DOWN:
				{
					if len(line_buffer_history) > 0 {
						if history_index > 0 {
							history_index--
							line_buffer = line_buffer_history[history_index]
							Render(line_buffer)
							length := len(line_buffer) + len(prompt) + 1
							MoveCursorEndOfLine(length)
							cursor.x = len(line_buffer)
						} else if history_index == 0 {
							history_index--
							line_buffer = cur_line_buffer
							Render(line_buffer)
							length := len(line_buffer) + len(prompt) + 1
							MoveCursorEndOfLine(length)
							cursor.x = len(line_buffer)
						}

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
		} else if is_tab(read_data) {
			search_root := find_search_root(tree.root, line_buffer)
			prefixs := search_prefix_tree(search_root)
			fmt.Printf("\n")
			for _, prefix := range prefixs {
				fmt.Printf("%s", string(line_buffer[:len(line_buffer)-1])+string(prefix))
			}
			fmt.Printf("\n")

		} else if is_enter(read_data) {
			history_index = -1

			if len(line_buffer) > 0 {
				line_buffer_history = append([][]byte{line_buffer}, line_buffer_history...)
			}
			add_word_to_prefix_tree(tree.root, append(line_buffer, byte(10)))

			fmt.Print("\n")
			callback(line_buffer)

			line_buffer = []byte{}
			fmt.Print("\n")
			fmt.Print(prompt)
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
