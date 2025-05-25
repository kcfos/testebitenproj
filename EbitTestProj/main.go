package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
}

var (
	G *Game
	//go:embed *
	AllAssets embed.FS

	DefaultSpriteShader   *ebiten.Shader
	WaterReflectionShader *ebiten.Shader

	Indices = [6]uint16{0, 1, 2, 2, 1, 3} //Standard Rect Double Triangle

	Prop1     Prop
	WaterProp Prop
)

func main() {
	G = &Game{}

	var err error
	DefaultSpriteShader, err = ebiten.NewShader(LoadAssetData("Assets/Default.kage"))
	if err != nil {
		panic(err)
	}
	WaterReflectionShader, err = ebiten.NewShader(LoadAssetData("Assets/WaterReflect.kage"))
	if err != nil {
		panic(err)
	}

	Prop1 = Prop{Op: &ebiten.DrawTrianglesShaderOptions{}, Shader: DefaultSpriteShader}
	gopherImg, _, err := image.Decode(bytes.NewReader(LoadAssetData("Assets/gopher.png")))
	if err != nil {
		panic(err)
	}

	Prop1.Op.Images[0] = ebiten.NewImageFromImage(gopherImg)
	Prop1.RefreshVerts(0,0)

	WaterProp = Prop{Op: &ebiten.DrawTrianglesShaderOptions{}, Shader: WaterReflectionShader}
	WaterProp.Op.Images[0] = ebiten.NewImage(600, 600)
	WaterProp.RefreshVerts(100,100)

	ebiten.SetWindowTitle("World")
	if err := ebiten.RunGame(G); err != nil {
		panic(err)
	}
}

func (G *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		fmt.Println("Creating new water image")
	WaterProp.Op.Images[0] = ebiten.NewImage(600, 600)
	}
	return nil
}

func (G *Game) Draw(screen *ebiten.Image) {
	Prop1.Draw(screen)
	Prop1.Draw(WaterProp.Op.Images[0])
	WaterProp.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

type Prop struct {
	Verts  [4]ebiten.Vertex
	Shader *ebiten.Shader
	Op     *ebiten.DrawTrianglesShaderOptions
}

func (P *Prop) Draw(targ *ebiten.Image) {
	targ.DrawTrianglesShader(P.Verts[:], Indices[:], P.Shader, P.Op)
}
func (P *Prop) RefreshVerts(DstOffsetX float32, DstOffsetY float32) {
	P.Verts[0].SrcX, P.Verts[0].SrcY = 0, 0                                                                         //TL
	P.Verts[1].SrcX, P.Verts[1].SrcY = float32(P.Op.Images[0].Bounds().Dx()), 0                                     //TR
	P.Verts[2].SrcX, P.Verts[2].SrcY = 0, float32(P.Op.Images[0].Bounds().Dy())                                     //BL
	P.Verts[3].SrcX, P.Verts[3].SrcY = float32(P.Op.Images[0].Bounds().Dx()), float32(P.Op.Images[0].Bounds().Dy()) //BR

	P.Verts[0].DstX, P.Verts[0].DstY = DstOffsetX, DstOffsetY                                                                         //TL
	P.Verts[1].DstX, P.Verts[1].DstY = DstOffsetX+float32(P.Op.Images[0].Bounds().Dx()), DstOffsetY                                     //TR
	P.Verts[2].DstX, P.Verts[2].DstY = DstOffsetX, DstOffsetY+float32(P.Op.Images[0].Bounds().Dy())                                     //BL
	P.Verts[3].DstX, P.Verts[3].DstY = DstOffsetX+float32(P.Op.Images[0].Bounds().Dx()), DstOffsetY+float32(P.Op.Images[0].Bounds().Dy()) //BR
}

func LoadAssetData(filepath string) []byte {
	data, err := AllAssets.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return data
}
