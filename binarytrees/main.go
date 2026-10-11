package main

import (
	"fmt"
	"math/rand"
	"time"

	"go-homework/internal/binarytree"
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	trees := []*binarytree.Tree{
		binarytree.New(1, rng),
		binarytree.New(1, rng),
		binarytree.New(1, rng),
		binarytree.New(1, rng),
		binarytree.New(2, rng),
	}

	for i, tree := range trees {
		fmt.Printf("tree %d: %v (root=%d)\n", i+1, binarytree.Values(tree), tree.Value)
	}
	fmt.Printf("first four contain the same values: %t\n", binarytree.Same(trees[0], trees[1], trees[2], trees[3]))
	fmt.Printf("tree 5 contains a different value set: %t\n", !binarytree.Same(trees[0], trees[1], trees[2], trees[4]))
}
