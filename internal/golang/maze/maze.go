package main

import (
	"fmt"
	"os"
)

/*
	广度优先搜索走迷宫
		循环创建二维数组 slice
		使用 slice 来实现队列
		用 Fscanf 读取文件
		对 point 进行抽象
*/

const filename = "maze/maze.in"

func readMaze(filename string) [][]int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var row, col int
	fmt.Fscanf(file, "%d %d", &row, &col)

	maze := make([][]int, 6)
	for i := range maze {
		maze[i] = make([]int, col)
		for j := range maze[i] {
			fmt.Fscanf(file, "%d", &maze[i][j])
		}
	}

	return maze
}

func walk(maze [][]int, start, end point) [][]int {
	steps := make([][]int, 6)
	for i := range steps {
		steps[i] = make([]int, 5)
	}

	Q := []point{start}

	for len(Q) > 0 {
		cur := Q[0]
		if cur == end { // Meet end quit.
			break
		}

		Q = Q[1:]

		// Traverse dirs to move different direction.
		for _, dir := range dirs {
			next := cur.add(dir)

			// Meet wall continue.
			val, ok := next.at(maze)
			if !ok || val == 1 {
				continue
			}

			// The road is not zero continue.
			val, ok = next.at(steps)
			if !ok || val != 0 {
				continue
			}

			// Meet start continue.
			if next == start {
				continue
			}

			curSteps, _ := cur.at(steps)
			steps[next.i][next.j] = curSteps + 1

			Q = append(Q, next)
		}
	}

	return steps
}

func printMaze(src [][]int) {
	for _, row := range src {
		for _, val := range row {
			fmt.Printf("%3d ", val)
		}
		fmt.Println()
	}
}

type point struct {
	i, j int
}

// abstract four directions
var dirs = [4]point{
	{-1, 0}, // left
	{0, -1}, // down
	{1, 0},  // right
	{0, 1},  // up
}

// Move point.
func (p point) add(r point) point {
	return point{p.i + r.i, p.j + r.j}
}

// Judge whether or not have value about the point.
func (p point) at(grid [][]int) (int, bool) {
	// verify params
	if p.i < 0 || p.i >= len(grid) {
		return 0, false
	}
	if p.j < 0 || p.j >= len(grid[p.i]) {
		return 0, false
	}

	return grid[p.i][p.j], true
}

func main() {
	// get maze
	maze := readMaze(filename)

	// walk maze
	steps := walk(maze, point{0, 0}, point{len(maze) - 1, len(maze[0]) - 1})

	// print steps
	printMaze(steps)
}
