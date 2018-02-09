package main

import (
	"fmt"
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
)

const (
	wall   = '#'
	void   = '-'
	target = '.'
)

type cell struct {
	value    rune
	backbone bool
	router   bool
}

type file struct {
	rows         int
	columns      int
	routerRadius int
	backboneCost int
	routerCost   int
	budget       int
	grid         [][]cell
}

type gridInfo struct {
	walls   int
	voids   int
	targets int
}

func main() {
	inputFile := readFile(os.Args[1])
	simulate(&inputFile)
}

func simulate(inputFile *file) {
	stats := getGridInfo(*inputFile)
	fmt.Println("Walls: ", stats.walls, ", Voids: ", stats.voids, ", Targets: ", stats.targets)
	gridFill(inputFile)
}

func gridFill(inputFile *file){
	gridSize := (inputFile.routerRadius * 2) + 1
	routerCount := 0
	for i := gridSize; i < inputFile.rows; i += gridSize {
		for j := gridSize; j < inputFile.columns; j += gridSize {
			if addRouter(i, j, inputFile) {
				routerCount++
			}
		}
	}
	fmt.Println(routerCount, " routers placed.")
}

func addRouter(x, y int, inputFile *file) bool {
	if inputFile.grid[x][y].value == target {
		inputFile.grid[x][y].router = true
		return true
	}
	return false
}

//Returns the total number of each feature: Walls, Voids and Targets
func getGridInfo(inputFile file) gridInfo {
	gInf := gridInfo{0, 0, 0}

	var walls, voids, targets = 0, 0, 0
	for i := 0; i < inputFile.rows; i++ {
		for j := 0; j < inputFile.columns; j++ {
			//fmt.Println("row ", i, " column ", j, " rows ", len(inputFile.grid), " columns ", len(inputFile.grid[i]))
			if inputFile.grid[i][j].value == wall {
				walls++
			} else if inputFile.grid[i][j].value == void {
				voids++
			} else if inputFile.grid[i][j].value == target {
				targets++
			}
		}
	}

	gInf.targets = targets
	gInf.voids = voids
	gInf.walls = walls

	return gInf
}

//Reads in the file
func readFile(fileName string) file {
	var f, err = os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineIndex := 0

	var inputFile file

	for scanner.Scan() {
		line := strings.Fields(scanner.Text())

		switch lineIndex {
		//First Row: Row count, Column count and Router Radius
		case 0:
			fmt.Println(line)
			inputFile.rows, _ = strconv.Atoi(line[0])
			inputFile.columns, _ = strconv.Atoi(line[1])
			grid := make([][]cell, inputFile.rows)
			for i := 0; i < inputFile.rows; i++ {
				grid[i] = make([]cell, inputFile.columns)
			}
			inputFile.grid = grid
			inputFile.routerRadius, _ = strconv.Atoi(line[2])
		//Second Row: Backbone Cost, Router Cost and Budget
		case 1:
			fmt.Println(line)
			inputFile.backboneCost, _ = strconv.Atoi(line[0])
			inputFile.routerCost, _ = strconv.Atoi(line[1])
			inputFile.budget, _ = strconv.Atoi(line[2])
		//Third Row: Grid X Y coords of initial backbone connection
		case 2:
			fmt.Println(line)
			x, _ := strconv.Atoi(line[0])
			y, _ := strconv.Atoi(line[1])
			inputFile.grid[x][y].backbone = true
		//Remaining line are the floor plan
		default:
			for i, char := range line[0] {
				inputFile.grid[lineIndex-3][i].value = char
			}
		}
		lineIndex++
	}
	return inputFile
}
