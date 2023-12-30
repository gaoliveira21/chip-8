package graphics

type Graphics struct {
	display [][]byte
	Width   int
	Height  int
}

func newDisplay(h int, w int) [][]byte {
	d := make([][]byte, h)

	for i := 0; i < h; i++ {
		d[i] = make([]byte, w)
	}

	return d
}

func NewGraphics() *Graphics {
	w := 0x40
	h := 0x20

	return &Graphics{
		Width:   w,
		Height:  h,
		display: newDisplay(h, w),
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

func (g *Graphics) EnableHighResolutionMode() {
	g.Width = 0x80
	g.Height = 0x40

	g.display = newDisplay(g.Height, g.Width)
}

func (g *Graphics) DisableHighResolutionMode() {
	g.Width = 0x40
	g.Height = 0x20

	g.display = newDisplay(g.Height, g.Width)
}

func (g *Graphics) ScrollDown(shift uint8) {
	s := int(shift)

	for i := g.Height - 1; i >= s; i-- {
		g.display[i] = g.display[i-s]
	}

	for i := 0; i < s; i++ {
		g.display[i] = make([]byte, g.Width)
	}
}

func (g *Graphics) ScrollRight() {
	s := 4

	for i := 0; i < g.Height; i++ {
		for j := g.Width - 1; j >= 0; j-- {
			if j < s {
				g.display[i][j] = 0x0
			} else {
				g.display[i][j] = g.display[i][j-s]
			}
		}
	}
}

func (g *Graphics) ScrollLeft() {
	s := 4

	for i := 0; i < g.Height; i++ {
		for j := 0; j < g.Width; j++ {
			if j < (g.Width - s) {
				g.display[i][j] = g.display[i][j+s]
			} else {
				g.display[i][j] = 0x0
			}
		}
	}
}
