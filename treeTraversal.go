package main

import (
	"fmt"
	"math/rand"
)

type Tree struct {
	Left *Tree
	Value int
	Right *Tree
}

func Walk(tree *Tree, channel chan int) {
	if tree == nil {
		return
	}
	Walk(tree.Left, channel)
	channel <- tree.Value
	Walk(tree.Right, channel)
}

func Walker(tree *Tree) <-chan int {
	channel := make(chan int)
	go func() {
		Walk(tree, channel)
		close(channel)
	}()
	return channel
}

func Compare(originalTree, newTree *Tree) bool {
	originalChannel, newChannel := Walker(originalTree), Walker(newTree)
	for {
		originalValue, originalError := <-originalChannel
		newValue, newError := <-newChannel

		if !originalError || !newError {
			if originalError == newError {
				fmt.Println("Error states the same")
			} else {
				fmt.Println("One error occurred")
			}
			return originalError == newError
		}
		if originalValue != newValue {
			fmt.Println("Different values")
			break
		}
	}
	return false
}

func New(nodesNumber, valueMultiplier int) *Tree {
	var tree *Tree
	for _, value := range rand.Perm(nodesNumber) {
		tree = insert(tree, (1+value)*valueMultiplier)
	}
	return tree
}

func insert(tree *Tree, value int) *Tree {
	if tree == nil {
		return &Tree{nil, value, nil}
	}
	if value < tree.Value {
		tree.Left = insert(tree.Left, value)
		return tree
	}
	tree.Right = insert(tree.Right, value)
	return tree
}

func main() {
	originalTree := New(5, 1)
	fmt.Println(Compare(originalTree, New(5, 1)), "\n-- Same Contents")
	fmt.Println(Compare(originalTree, New(4, 1)), "\n-- Differing Sizes")
	fmt.Println(Compare(originalTree, New(5, 2)), "\n-- Differing Values")
	fmt.Println(Compare(originalTree, New(6, 2)), "\n-- Dissimilar")
}