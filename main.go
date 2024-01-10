package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type PlayerList struct {
	Name    string
	Ping    int
	HeadURL string
}

func registerCustomFont(dc *gg.Context) error {
	fontPath := "./assets/mc.otf" // Replace with the path to your TTC (TrueType Collection) font file

	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return err
	}

	ttc, err := opentype.ParseCollection(fontBytes)
	if err != nil {
		return err
	}

	numFonts := ttc.NumFonts()
	if numFonts == 0 {
		return fmt.Errorf("font collection has no fonts")
	}

	fnt, err := ttc.Font(0) // get the first font
	if err != nil {
		return fmt.Errorf("could not open font #0: %+v", err)
	}

	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	dc.SetFontFace(face)
	return nil
}

func loadPingImage(ping int) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	pingImagePath := filepath.Join(currentDir, "assets")

	switch {
	case ping < 0:
		return filepath.Join(pingImagePath, "signal_0.png"), nil
	case ping <= 150:
		return filepath.Join(pingImagePath, "signal_5.png"), nil
	case ping <= 300:
		return filepath.Join(pingImagePath, "signal_4.png"), nil
	case ping <= 600:
		return filepath.Join(pingImagePath, "signal_3.png"), nil
	case ping <= 1000:
		return filepath.Join(pingImagePath, "signal_2.png"), nil
	default:
		return filepath.Join(pingImagePath, "signal_1.png"), nil
	}
}

func loadImageFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func loadImageFromPath(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func drawBlock(dc *gg.Context, x, z int, player PlayerList) {
	dc.SetRGB(0.827, 0.827, 0.827)
	dc.DrawRectangle(float64(x+2), float64(z), 276, 20)
	dc.Fill()

	// Register the custom font
	if err := registerCustomFont(dc); err != nil {
		log.Fatal("Error registering custom font:", err)
	}

	dc.SetRGB(0, 0, 0)

	avatar, err := loadImageFromURL(player.HeadURL)
	if err != nil {
		log.Println("Error loading avatar:", err)
		return
	}
	avatar = resize.Resize(16, 16, avatar, resize.Lanczos3)
	dc.DrawImageAnchored(avatar, int(float64(x+5)), int(float64(z+2)), 0, 0)

	pingImagePath, err := loadPingImage(player.Ping)
	if err != nil {
		log.Println("Error loading ping image path:", err)
		return
	}

	pingImage, err := loadImageFromPath(pingImagePath)
	if err != nil {
		log.Println("Error loading ping image:", err)
		return
	}
	pingImage = resize.Resize(16, 16, pingImage, resize.Lanczos3)
	dc.DrawImageAnchored(pingImage, int(float64(x+259)), int(float64(z+2)), 0, 0)

	dc.DrawStringAnchored(player.Name, float64(x+23), float64(z+16), 0, 0)
}

func renderTab(players []PlayerList) *gg.Context {

	const canvasHeight = 350
	canvasWidth := calculateCanvasWidth(len(players))

	dc := gg.NewContext(canvasWidth+2, canvasHeight)
	dc.SetRGBA(0, 0, 0, 0)
	dc.Clear()

	z := 0
	x := 0

	for _, player := range players {
		if z > 330 {
			x = x + 278
			z = 0
		}
		drawBlock(dc, x, z, player)
		z = z + 22
	}

	return dc
}

func calculateCanvasWidth(playersCount int) int {
	const blockSize = 278
	return (playersCount/16 + 1) * blockSize
}

func main() {

	players := []PlayerList{
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
		{Name: "Player1", Ping: 100, HeadURL: "https://mc-heads.net/avatar/Player1/16"},
		{Name: "Player2", Ping: 200, HeadURL: "https://mc-heads.net/avatar/Player2/16"},
	}

	dc := renderTab(players)

	// Determine the directory of the executabl
	if err := dc.SavePNG("tablist.png"); err != nil {
		fmt.Println("Error saving PNG:", err)
		return
	}

	fmt.Printf("Tablist image saved successfully")
}
