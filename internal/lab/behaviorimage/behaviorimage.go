package behaviorimage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"regexp"
	"sort"
	"strings"

	"github.com/LuigiD5555/origami/internal/fixedcarrier"
	"github.com/LuigiD5555/origami/internal/lab/renderer"
)

const (
	SchemaR0        = "origami.behavior-recipe.r0"
	RenderSchemaR0  = "origami.behavior-image.r0"
	DefaultWidth    = 960
	DefaultHeight   = 640
	DefaultMaxSteps = 8
)

type Step struct {
	ID          string `json:"id"`
	Capability  string `json:"capability"`
	Instruction string `json:"instruction"`
	InputRef    string `json:"input_ref,omitempty"`
	OutputRef   string `json:"output_ref,omitempty"`
	OnSuccess   string `json:"on_success,omitempty"`
	OnUnknown   string `json:"on_unknown,omitempty"`
	Halt        bool   `json:"halt,omitempty"`
}

type Recipe struct {
	Schema         string `json:"schema"`
	ID             string `json:"id"`
	Boot           string `json:"boot"`
	OutputContract string `json:"output_contract,omitempty"`
	Steps          []Step `json:"steps"`
}

type RenderResult struct {
	Schema       string `json:"schema"`
	RecipeID     string `json:"recipe_id"`
	RecipeSHA256 string `json:"recipe_sha256"`
	PNGSHA256    string `json:"png_sha256"`
	PNGBytes     int    `json:"png_bytes"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	StepCount    int    `json:"step_count"`
}

var capabilityName = regexp.MustCompile(`^[A-Z][A-Z0-9_]{1,47}$`)

func Validate(r Recipe) error {
	if r.Schema != SchemaR0 {
		return fmt.Errorf("recipe schema must be %q", SchemaR0)
	}
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.Boot) == "" {
		return fmt.Errorf("recipe id and boot are required")
	}
	if len(r.Steps) == 0 || len(r.Steps) > DefaultMaxSteps {
		return fmt.Errorf("recipe must contain 1..%d steps", DefaultMaxSteps)
	}
	ids := map[string]int{}
	outputs := map[string]bool{}
	for i, s := range r.Steps {
		if strings.TrimSpace(s.ID) == "" || strings.TrimSpace(s.Instruction) == "" {
			return fmt.Errorf("step %d requires id and instruction", i+1)
		}
		if !capabilityName.MatchString(s.Capability) {
			return fmt.Errorf("step %q has invalid capability %q", s.ID, s.Capability)
		}
		if strings.ContainsAny(s.Instruction, "\r\n") || len([]rune(s.Instruction)) > 96 {
			return fmt.Errorf("step %q instruction must be one line <= 96 runes", s.ID)
		}
		if _, dup := ids[s.ID]; dup {
			return fmt.Errorf("duplicate step id %q", s.ID)
		}
		ids[s.ID] = i
		if s.OutputRef != "" {
			if outputs[s.OutputRef] {
				return fmt.Errorf("duplicate output_ref %q", s.OutputRef)
			}
			outputs[s.OutputRef] = true
		}
	}
	validTarget := func(target string, from int) error {
		if target == "" || target == "HALT" || target == "UNKNOWN" {
			return nil
		}
		idx, ok := ids[target]
		if !ok {
			return fmt.Errorf("transition references unknown step %q", target)
		}
		if idx <= from {
			return fmt.Errorf("R0 recipes are forward-only; step %q cannot jump to %q", r.Steps[from].ID, target)
		}
		return nil
	}
	for i, s := range r.Steps {
		if err := validTarget(s.OnSuccess, i); err != nil {
			return err
		}
		if err := validTarget(s.OnUnknown, i); err != nil {
			return err
		}
	}
	return nil
}

func CanonicalJSON(r Recipe) ([]byte, error) {
	if err := Validate(r); err != nil {
		return nil, err
	}
	// Struct fields and step order are canonical in R0; maps are deliberately absent.
	return json.Marshal(r)
}

func Render(r Recipe, width, height int) ([]byte, RenderResult, error) {
	canonical, err := CanonicalJSON(r)
	if err != nil {
		return nil, RenderResult{}, err
	}
	if width <= 0 {
		width = DefaultWidth
	}
	if height <= 0 {
		height = DefaultHeight
	}
	img := image.NewGray(image.Rect(0, 0, width, height))
	for i := range img.Pix {
		img.Pix[i] = 255
	}
	lines := visibleLines(r)
	scale := 2
	lineH := 18
	x := 18
	y := 18
	for _, line := range lines {
		if fixedcarrier.TextWidth(line, scale) > width-2*x {
			return nil, RenderResult{}, fmt.Errorf("visible recipe line too wide: %q", line)
		}
		if y+14*scale >= height-12 {
			return nil, RenderResult{}, fmt.Errorf("visible recipe exceeds canvas height")
		}
		fixedcarrier.DrawText(img, x, y, scale, line, 0)
		y += lineH
	}
	// A deterministic footer makes the image self-identifying without exposing hidden answers.
	footer := "END RECIPE | EXECUTE ONLY DECLARED STEPS"
	fixedcarrier.DrawText(img, 18, height-28, 1, footer, 64)
	pngBytes, err := renderer.EncodePNG(img)
	if err != nil {
		return nil, RenderResult{}, err
	}
	rs := sha256.Sum256(canonical)
	ps := sha256.Sum256(pngBytes)
	result := RenderResult{Schema: RenderSchemaR0, RecipeID: r.ID, RecipeSHA256: hex.EncodeToString(rs[:]), PNGSHA256: hex.EncodeToString(ps[:]), PNGBytes: len(pngBytes), Width: width, Height: height, StepCount: len(r.Steps)}
	return pngBytes, result, nil
}

func visibleLines(r Recipe) []string {
	lines := []string{
		"BOOT ORIGAMI BEHAVIOR R0",
		"RECIPE " + sanitize(r.ID),
		"MODE READ ONE STEP AT A TIME",
		"IF A STEP IS UNKNOWN FOLLOW ITS UNKNOWN ROUTE",
	}
	for i, s := range r.Steps {
		lines = append(lines, fmt.Sprintf("STEP %d %s", i+1, sanitize(s.ID)))
		lines = append(lines, "CAP "+sanitize(s.Capability))
		if s.InputRef != "" {
			lines = append(lines, "IN "+sanitize(s.InputRef))
		}
		lines = append(lines, "DO "+sanitize(s.Instruction))
		if s.OutputRef != "" {
			lines = append(lines, "OUT "+sanitize(s.OutputRef))
		}
		success := s.OnSuccess
		if success == "" {
			if s.Halt {
				success = "HALT"
			} else if i+1 < len(r.Steps) {
				success = r.Steps[i+1].ID
			} else {
				success = "HALT"
			}
		}
		unknown := s.OnUnknown
		if unknown == "" {
			unknown = "UNKNOWN"
		}
		lines = append(lines, "OK "+sanitize(success)+" | UNKNOWN "+sanitize(unknown))
	}
	if r.OutputContract != "" {
		lines = append(lines, "OUTPUT "+sanitize(r.OutputContract))
	}
	return lines
}

func sanitize(s string) string {
	s = strings.ToUpper(strings.Join(strings.Fields(s), " "))
	// The bootstrap font has a deliberately small alphabet; replace unsupported punctuation.
	repl := strings.NewReplacer("{", "[", "}", "]", "\"", "", "'", "", ";", ",", "@", "A", "%", "PCT")
	return repl.Replace(s)
}

func DecodeRecipe(data []byte) (Recipe, error) {
	var r Recipe
	if err := json.Unmarshal(data, &r); err != nil {
		return Recipe{}, err
	}
	return r, Validate(r)
}

func EqualRender(a, b []byte) bool { return bytes.Equal(a, b) }

func SortedCapabilities(r Recipe) []string {
	seen := map[string]bool{}
	for _, s := range r.Steps {
		seen[s.Capability] = true
	}
	out := make([]string, 0, len(seen))
	for c := range seen {
		out = append(out, c)
	}
	sort.Strings(out)
	return out
}

var _ = color.Gray{}
