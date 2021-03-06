package main

import (
  "github.com/jeffail/gabs"
  "io/ioutil"
  "log"
  "math"
)

//These represent the tile order in Tiled sprite
const (
  TypeGround = 1
  TypeGoal = 2
  TypeStart = 3
  TypeBlock = 4
)

var TypeNames = map[int]string{
  TypeStart: "Start",
  TypeGoal: "Goal",
  TypeGround: "Ground",
  TypeBlock: "Block",
}

// Two dimensional map, consisting of Tiles
type Map map[int]map[int]*Tile
type Tile struct {
  X, Y int
  Type int
  M *Map
}

func (m Map) addTile(t *Tile) {
  if m[t.X] == nil {
		m[t.X] = map[int]*Tile{}
	}
  m[t.X][t.Y] = t
  t.M = &m
}
func (m Map) getTile(x, y int) *Tile {
  if (m[x] != nil) {
    return m[x][y]
  }
  return nil
}

func calcTileCoords(index, mapWidth int) (int, int) {
  x := index / mapWidth
  y := index - x * mapWidth
  return x, y
}

func parseMapFile(filePath string) *Map {
  mapData, err := ioutil.ReadFile(filePath)
  if err != nil {
    log.Fatal(err)
  }


  jsonParsed, err := gabs.ParseJSON(mapData)
  if err != nil {
    log.Fatal(err)
  }

  mapWidth := int(jsonParsed.Path("width").Data().(float64))
  // Only first layer from Tiled data is used, rest ignored
  firstMapLayer := jsonParsed.Search("layers", "data").Index(0)
  tileTypes, _ := firstMapLayer.Children()

  m := Map{}
  for index, tileType := range tileTypes {
    x, y := calcTileCoords(index, mapWidth)
    tile := Tile{X: x, Y: y, Type: int(tileType.Data().(float64))}
    m.addTile(&tile)
  }

  return &m
}

func (t *Tile) getWalkableNeighbours() []*Tile {
  var neighbours []*Tile
  // Path cannot be diagonal
  for _, offset := range [][]int {
      {-1, 0},
      {0, -1},
      {0, 1},
      {1, 0},
    } {
    if t := t.M.getTile(t.X+offset[0], t.Y+offset[1]); t != nil && t.Type == TypeGround {
      neighbours = append(neighbours, t)
    }
  }
  return neighbours
}
func (from *Tile) calcDistance(to *Tile) int {
  xDiff := math.Abs(float64(to.X - from.X))
  yDiff := math.Abs(float64(to.Y - from.Y))
  return int(xDiff + yDiff)
}
func (m *Map) getStartTile() *Tile {
  return m.findTileOfType(TypeStart)
}
func (m *Map) getGoalTile() *Tile {
  return m.findTileOfType(TypeGoal)
}

func (m *Map) findTileOfType(tileType int) *Tile {
  for _, row := range *m {
    for _, tile := range row {
      if (tile.Type == tileType) {
        return tile
      }
    }
  }
  return nil
}
