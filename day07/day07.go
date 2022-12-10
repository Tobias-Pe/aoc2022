package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type data interface {
	getSize() int
	getName() string
}

type file struct {
	size int
	name string
}

func (file file) getSize() int {
	return file.size
}

func (file file) getName() string {
	return file.name
}

type directory struct {
	parent  *directory
	name    string
	content []data
}

func (directory directory) listContent() string {
	names := ""
	for _, data := range directory.content {
		names += data.getName()
		names += ", "
	}
	return names
}

func newDirectory(parent *directory, name string) *directory {
	return &directory{
		parent:  parent,
		name:    name,
		content: []data{},
	}
}

func (directory directory) getSize() int {
	sum := 0
	for _, data := range directory.content {
		sum += data.getSize()
	}
	return sum
}

func (directory directory) getName() string {
	return directory.name
}

func readFile() {
	content, err := ioutil.ReadFile("./input.txt")
	if err != nil {
		fmt.Printf("Error on reading file: %s", err.Error())
	}
	lines := string(content)
	topLevel = newDirectory(topLevel, "/")
	currentDir := topLevel
	for _, line := range strings.Split(lines, "\n") {
		if strings.HasPrefix(line, "$ cd") {
			targetForDir := strings.TrimLeft(line, "$ cd")
			currentDir = changeDirectory(currentDir, targetForDir)
		} else if !strings.HasPrefix(line, "$") && len(line) > 1 {
			args := strings.Split(line, " ")
			populateDirectory(currentDir, args[0], args[1])
		}
	}
}

func populateDirectory(current *directory, argument string, name string) {
	if argument == "dir" {
		current.content = append(current.content, newDirectory(current, name))
	} else {
		convertedSize, err := strconv.Atoi(argument)
		if err != nil {
			fmt.Println(err)
			return
		}
		current.content = append(current.content, &file{
			size: convertedSize,
			name: name,
		})
	}
}

func changeDirectory(current *directory, target string) *directory {
	switch target {
	case "/":
        return topLevel
	case "..":
        return current.parent
	default:
		for _, data := range current.content {
			if directory, isDirectory := data.(*directory); isDirectory {
                return directory
			}
		}
    }
    fmt.Println("No Directory ", target, " found :(")
    return nil
}

func sumDirectoriesSizeWithAtMost(dir *directory, max int) int {
	sum := dir.getSize()
    if sum < max {
        fmt.Println(dir.name,sum)
		return sum
	}
	sum = 0
	for _, data := range dir.content {
		if directory, isDirectory := data.(*directory); isDirectory {
            sum += sumDirectoriesSizeWithAtMost(directory, max)
            fmt.Println(dir.name,sum)
		}
	}
	return 0
}

var topLevel *directory

func main() {
	readFile()
	fmt.Println("Part01: ", sumDirectoriesSizeWithAtMost(topLevel, 100000))
}
