package readline

type prefix_node struct {
	children map[byte]*prefix_node
	value    byte
}

type prefix_tree struct {
	root *prefix_node
}

func add_word_to_prefix_tree(root *prefix_node, word []byte) {

	if len(word) > 0 {

		if child, ok := root.children[word[0]]; ok {
			add_word_to_prefix_tree(child, word[1:])
		} else {
			new_root := prefix_node{value: word[0], children: map[byte]*prefix_node{}}
			root.children[word[0]] = &new_root
			add_word_to_prefix_tree(&new_root, word[1:])
		}
	}
}

func search_prefix_tree(root *prefix_node) [][]byte {

	prefixs := [][]byte{}

	if len(root.children) > 0 {
		for _, child := range root.children {
			prefix := search_prefix_tree(child)
			prefixs = append(prefixs, prefix...)
		}
		if root.value != 0 {
			for i := 0; i < len(prefixs); i++ {
				prefixs[i] = append([]byte{root.value}, prefixs[i]...)
			}
		}
	} else {
		prefixs = [][]byte{[]byte{root.value}}
	}

	return prefixs

}

func find_search_root(root *prefix_node, prefix []byte) *prefix_node {

	if len(prefix) > 0 {
		if next_node, ok := root.children[prefix[0]]; ok {
			return find_search_root(next_node, prefix[1:])
		}
	} else {
		return root
	}

	return root

}
