package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LuigiD5555/origami/internal/lab/behaviorimage"
)

func main() {
	recipePath := flag.String("recipe", "", "behavior recipe JSON")
	outPath := flag.String("out", "", "output PNG")
	manifestPath := flag.String("manifest", "", "optional render manifest JSON (defaults to <out>.json)")
	width := flag.Int("width", behaviorimage.DefaultWidth, "canvas width")
	height := flag.Int("height", behaviorimage.DefaultHeight, "canvas height")
	flag.Parse()
	if *recipePath == "" || *outPath == "" {
		die(fmt.Errorf("-recipe and -out are required"))
	}
	raw, err := os.ReadFile(*recipePath)
	die(err)
	recipe, err := behaviorimage.DecodeRecipe(raw)
	die(err)
	pngBytes, result, err := behaviorimage.Render(recipe, *width, *height)
	die(err)
	die(os.MkdirAll(filepath.Dir(*outPath), 0o755))
	die(os.WriteFile(*outPath, pngBytes, 0o644))
	target := *manifestPath
	if target == "" {
		target = *outPath + ".json"
	}
	manifest, err := json.MarshalIndent(result, "", "  ")
	die(err)
	die(os.WriteFile(target, append(manifest, '\n'), 0o644))
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	die(enc.Encode(result))
}
func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
