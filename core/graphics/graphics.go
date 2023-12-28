package graphics

const width = 0x40
const height = 0x20

type Graphics struct {
	display [height][width]byte
	Width   int
	Height  int
}

func NewGraphics() *Graphics {
	return &Graphics{
		Width:  width,
		Height: height,
	}
}

func (g *Graphics) Clear() {
	for i := 0; i < g.Height; i++ {
		for j := 0; j < g.Width; j++ {
			g.display[i][j] = 0x00
		}
	}
}

func (g *Graphics) GetPixel(y int, x int) byte {
	return g.display[y][x]
}

func (g *Graphics) SetPixel(y int, x int, b byte) {
	g.display[y][x] = b
}
