package main

import (
	"fmt"
	"math/rand"
)

// MapItem 定义迷宫元素枚举
type MapItem int

const (
	MapWall MapItem = iota
	MapGround
	MapEntrance
	MapExit
	MapOutside
	MapStep
)

// Point 定义点结构体
type Point struct {
	X, Y int
}

// Size 定义大小结构体
type Size struct {
	Width, Height int
}

// Rect 定义矩形结构体
type Rect struct {
	X, Y, Width, Height int
}

// Maze 迷宫结构
type Maze struct {
	Map       [][]MapItem // 迷宫地图
	MapSize   Size        // 迷宫尺寸
	PlayerPos Point       // 游戏者位置
}

// MakeMaze 生成迷宫（注：宽高必须是奇数）
func (m *Maze) MakeMaze(width, height int) {
	if width%2 != 1 {
		width += 1
	}

	if height%2 != 1 {
		height += 1
	}

	// 记录迷宫尺寸
	m.MapSize = Size{Width: width, Height: height}

	// 分配迷宫内存
	m.Map = make([][]MapItem, width+2)
	for x := 0; x < width+2; x++ {
		m.Map[x] = make([]MapItem, height+2)
		// 初始化为墙
		for y := 0; y < height+2; y++ {
			m.Map[x][y] = MapWall
		}
	}

	// 定义边界
	for x := 0; x <= width+1; x++ {
		m.Map[x][0] = MapGround
		m.Map[x][height+1] = MapGround
	}

	for y := 1; y <= height; y++ {
		m.Map[0][y] = MapGround
		m.Map[width+1][y] = MapGround
	}

	// 定义入口和出口
	m.Map[1][2] = MapEntrance
	m.Map[width][height-1] = MapExit

	// 设置玩家初始位置
	m.PlayerPos = Point{X: 1, Y: 2}

	// 从任意点开始遍历生成迷宫
	x := ((rand.Intn(width-1) & 0xfffe) + 2)
	y := ((rand.Intn(height-1) & 0xfffe) + 2)
	m.TravelMaze(x, y)

	// 将边界标记为迷宫外
	for x := 0; x <= width+1; x++ {
		m.Map[x][0] = MapOutside
		m.Map[x][height+1] = MapOutside
	}

	for y := 1; y <= height; y++ {
		m.Map[0][y] = MapOutside
		m.Map[width+1][y] = MapOutside
	}
}

// TravelMaze 生成迷宫：遍历 (x, y) 四周
func (m *Maze) TravelMaze(x, y int) {
	// 定义遍历方向
	directions := [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	// 将遍历方向乱序
	for i := 0; i < 4; i++ {
		n := rand.Intn(4)
		directions[i], directions[n] = directions[n], directions[i]
	}

	// 尝试周围四个方向
	m.Map[x][y] = MapGround
	for i := 0; i < 4; i++ {
		dx, dy := directions[i][0], directions[i][1]
		if m.Map[x+2*dx][y+2*dy] == MapWall {
			m.Map[x+dx][y+dy] = MapGround
			m.TravelMaze(x+dx*2, y+dy*2) // 递归
		}
	}
}

// MapRoute 实现迷宫求解，使用深度优先搜索算法
func (m *Maze) MapRoute() {
	// 定义方向：上、右、下、左
	directions := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	// 找到入口位置
	var start Point
	found := false
	for i := 0; i < m.MapSize.Width+2 && !found; i++ {
		for j := 0; j < m.MapSize.Height+2; j++ {
			if m.Map[i][j] == MapEntrance {
				start = Point{X: i, Y: j}
				found = true
				break
			}
		}
	}

	if !found {
		fmt.Println("未找到迷宫入口")
		return
	}

	// 创建访问标记数组
	visited := make([][]bool, m.MapSize.Width+2)
	for i := range visited {
		visited[i] = make([]bool, m.MapSize.Height+2)
	}

	// 记录路径的栈
	path := []Point{start}
	visited[start.X][start.Y] = true

	// 深度优先搜索函数
	var dfs func(x, y int) bool
	dfs = func(x, y int) bool {
		// 如果到达出口，返回成功
		if m.Map[x][y] == MapExit {
			return true
		}

		// 尝试四个方向
		for _, dir := range directions {
			nx, ny := x+dir[0], y+dir[1]

			// 检查是否可以移动到(nx, ny)
			if nx >= 0 && nx < m.MapSize.Width+2 &&
				ny >= 0 && ny < m.MapSize.Height+2 &&
				(m.Map[nx][ny] == MapGround || m.Map[nx][ny] == MapExit) &&
				!visited[nx][ny] {

				visited[nx][ny] = true
				path = append(path, Point{X: nx, Y: ny})

				// 递归搜索
				if dfs(nx, ny) {
					return true
				}

				// 回溯
				path = path[:len(path)-1]
			}
		}

		return false
	}

	// 开始搜索
	if dfs(start.X, start.Y) {
		// 标记路径
		for i := 1; i < len(path)-1; i++ { // 跳过入口和出口
			point := path[i]
			m.Map[point.X][point.Y] = MapStep
		}
		fmt.Println("迷宫路径已找到并标记")
	} else {
		fmt.Println("未找到可行的迷宫路径")
	}
}

func (m *Maze) PrintMaze() {
	for j := 0; j < m.MapSize.Height+2; j++ {
		for i := 0; i < m.MapSize.Width+2; i++ {
			if m.Map[i][j] == MapWall {
				fmt.Print("\x1b[43m  \x1b[0m")
			} else if m.Map[i][j] == MapGround {
				fmt.Print("  ")
			} else if m.Map[i][j] == MapOutside {
				fmt.Print("\x1b[44m  \x1b[0m")
			} else if m.Map[i][j] == MapEntrance {
				fmt.Print("\x1b[42m  \x1b[0m")
			} else if m.Map[i][j] == MapExit {
				fmt.Print("\x1b[41m  \x1b[0m")
			} else if m.Map[i][j] == MapStep {
				fmt.Print("\x1b[45m  \x1b[0m")
			}
		}
		fmt.Println()
	}
}
