package readline

import (
	"fmt"
	"testing"
)

func TestPrefix(t *testing.T) {

	tree := prefix_tree{root: &prefix_node{children: map[byte]*prefix_node{}}}

	line_buffer := []byte("test_data")
	add_word_to_prefix_tree(tree.root, line_buffer)

	line_buffer = []byte("dog")
	add_word_to_prefix_tree(tree.root, line_buffer)

	line_buffer = []byte("tek")
	add_word_to_prefix_tree(tree.root, line_buffer)

	search_root := find_search_root(tree.root, []byte("te"))

	fmt.Printf("%+v", search_root)
	prefixs := search_prefix_tree(search_root)
	for _, prefix := range prefixs {
		fmt.Printf("\n%s", string(prefix))
	}

}
