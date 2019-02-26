package World

import (
	"encoding/json"
	"log"
	"os"
)

type World struct {
	Entities []Entity
}

type Entity struct {
	Points []Point
	Lines []Line
}

type Point struct {
	X float32
	Y float32
	Z float32
	ConnectedTo []int
}

type Line struct {
	P1 int
	P2 int
}

type Square struct {
	Origin Origin
	Height float32
	Width float32
	Depth float32
}

type Origin struct {
	X float32
	Y float32
	Z float32
}

type FileData struct {
	FileObjects []FileObject
}

type FileObject struct {
	Type string
	Data string
}

func NewWorld() *World {
	world := World{}
	world.Entities = []Entity{}

	return &world
}

func (w *World) Build(worldDescriptor string) error {
	log.Println("Building new world based on", worldDescriptor)
	worldFile, err := os.Open(worldDescriptor)
	if err != nil {
		log.Println("Error opening world descriptor:", err.Error())
		return err
	}

	buffer := make([]byte, 131072)

	n, err := worldFile.Read(buffer)

	if err != nil {
		log.Println("Error reading world descriptor:", err.Error())
		return err
	}

	fileData := FileData{}
	err = json.Unmarshal(buffer[:n], &fileData)
	if err != nil {
		log.Println("Error parsing world descriptor:", err.Error())
		return err
	}

	log.Println("World File contents:", fileData)

	for _, v := range fileData.FileObjects {
		err = v.ParseObject(w)
		if err != nil {
			return err
		}
	}

	log.Println("Built World state:", w)

	return nil
}

func (f *FileObject) ParseObject(world *World) error {
	switch f.Type {
	case "square":
		square := Square{}
		err := json.Unmarshal([]byte(f.Data), &square)
		if err != nil {
			log.Println("Error parsing square object data:", err.Error())
			return err
		}
		world.BuildSquare(square)
	}
	return nil
}

func (w *World) BuildSquare(data Square) {
	entity := Entity{}
	entity.Points = []Point{}
	entity.Lines = []Line{}

	//0
	entity.Points = append(entity.Points, Point{
		X: data.Origin.X,
		Y: data.Origin.Y,
		Z: data.Origin.Z,
		ConnectedTo: []int{1, 3, 4},
	})
	//1
	entity.Points = append(entity.Points, Point{
		X: data.Origin.X + data.Width,
		Y: data.Origin.Y,
		Z: data.Origin.Z,
		ConnectedTo: []int{0, 2, 5},
	})
	//2
	entity.Points = append(entity.Points, Point{
		X: data.Origin.X + data.Width,
		Y: data.Origin.Y,
		Z: data.Origin.Z + data.Depth,
		ConnectedTo: []int{1, 3, 6},
	})
	//3
	entity.Points = append(entity.Points, Point{
		X: data.Origin.X,
		Y: data.Origin.Y,
		Z: data.Origin.Z + data.Depth,
		ConnectedTo: []int{0, 2, 7},
	})
	//4
	entity.Points = append(entity.Points, Point{
		X: data.Origin.X,
		Y: data.Origin.Y + data.Height,
		Z: data.Origin.Z,
		ConnectedTo: []int{0, 5, 7},
	})
	//5
	entity.Points = append(entity.Points, Point{
		X: data.Origin.X + data.Width,
		Y: data.Origin.Y + data.Height,
		Z: data.Origin.Z,
		ConnectedTo: []int{1, 4, 6},
	})
	//6
	entity.Points = append(entity.Points, Point{
		X: data.Origin.X + data.Width,
		Y: data.Origin.Y + data.Height,
		Z: data.Origin.Z + data.Depth,
		ConnectedTo: []int{2, 5, 7},
	})
	//7
	entity.Points = append(entity.Points, Point{
		X: data.Origin.X,
		Y: data.Origin.Y + data.Height,
		Z: data.Origin.Z + data.Depth,
		ConnectedTo: []int{3, 4, 6},
	})

	for k, p := range entity.Points {
		ConnectionLoop:
		for _, c := range p.ConnectedTo {
			for _, l := range entity.Lines {
				if (l.P1 == k && l.P2 == c) || (l.P1 == c && l.P2 == k) {
					continue ConnectionLoop
				}
			}
			entity.Lines = append(entity.Lines, Line{
				P1: k,
				P2: c,
			})
		}
	}

	w.Entities = append(w.Entities, entity)
}
