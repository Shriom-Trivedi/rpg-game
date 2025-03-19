package main

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tileset interface {
	Img(id int) *ebiten.Image
}

type UniformTilsetJSON struct {
	Path string `json:"image"`
}

type UniformTilset struct {
	img *ebiten.Image
	gid int
}

func (u *UniformTilset) Img(id int) *ebiten.Image {
	// tilesetColumns := g.tilemapImg.Bounds().Dx() / 16 // Number of tiles per row in tileset
	id -= u.gid

	srcX := (id % 40) * 16
	srcY := (id / 40) * 16

	return u.img.SubImage(
		image.Rect(
			srcX, srcY, srcX+16, srcX+16,
		),
	).(*ebiten.Image)
}

type TIleJSON struct {
	Id     int    `json: "id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
	Height int    `json:"imageheight"`
}

// For dynamic tilesets
type DynTilesetJSON struct {
	Tiles []*TIleJSON `json:"tiles"`
}

type DynTileset struct {
	imgs []*ebiten.Image
	gid  int
}

func (d DynTileset) Img(id int) *ebiten.Image {
	id -= d.gid
	
	return d.imgs[id - 3] // Subtracting 3, since buildings.json id starting from 3. TODO: fix these weird indexes 
}

func NewTileset(path string, gid int) (Tileset, error) {

	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if strings.Contains(path, "buildings") {
		// return dyn tileset
		var dynTilesetJSON DynTilesetJSON
		err = json.Unmarshal(contents, &dynTilesetJSON)
		if err != nil {
			return nil, err
		}

		dynTileset := DynTileset{}
		dynTileset.gid = gid
		dynTileset.imgs = make([]*ebiten.Image, 0)

		for _, tileJSON := range dynTilesetJSON.Tiles {

			tileJSONPath := tileJSON.Path
			tileJSONPath = filepath.Clean(tileJSONPath)
			tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = filepath.Join("assets/", tileJSONPath)

			img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
			if err != nil {
				return nil, err
			}

			dynTileset.imgs = append(dynTileset.imgs, img)
		}

		return &dynTileset, nil

	}
	// return uniform tilset
	var uniformTilsetJSON UniformTilsetJSON
	err = json.Unmarshal(contents, &uniformTilsetJSON)
	if err != nil {
		return nil, err
	}

	UniformTilset := UniformTilset{}

	tileJSONPath := uniformTilsetJSON.Path
	tileJSONPath = filepath.Clean(tileJSONPath)
	tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
	tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
	tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
	tileJSONPath = filepath.Join("assets/", tileJSONPath)

	img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
	if err != nil {
		return nil, err
	}

	UniformTilset.img = img
	UniformTilset.gid = gid

	return &UniformTilset, nil
}
