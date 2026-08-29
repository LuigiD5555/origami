package microisa

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LuigiD5555/origami/internal/lab/glyphcalc"
)

type SIMDRound struct {
	Round  int      `json:"round"`
	Start  int      `json:"start"`
	Count  int      `json:"count"`
	Tuples []string `json:"tuples,omitempty"`
}

type SIMDPlan struct {
	Schema          string      `json:"schema"`
	Mode            string      `json:"mode"`
	LaneWidth       int         `json:"lane_width"`
	GlyphCount      int         `json:"glyph_count"`
	Columns         int         `json:"columns"`
	Rounds          []SIMDRound `json:"rounds"`
	SemanticMeaning string      `json:"semantic_meaning"`
	Rule            string      `json:"rule"`
}

type SIMDResponse struct {
	Schema  string       `json:"schema"`
	Model   string       `json:"model"`
	Status  string       `json:"status"`
	Rounds  []SIMDRound  `json:"rounds"`
	Metrics UsageMetrics `json:"metrics,omitempty"`
}

type SIMDCheck struct {
	Schema                  string       `json:"schema"`
	Model                   string       `json:"model"`
	LaneWidth               int          `json:"lane_width"`
	Total                   int          `json:"total"`
	ExactKnown              int          `json:"exact_known"`
	FalseKnown              int          `json:"false_known"`
	Unknown                 int          `json:"unknown"`
	Missing                 int          `json:"missing"`
	FalseExactZero          bool         `json:"false_exact_zero"`
	ExactTransportRecovered bool         `json:"exact_transport_recovered"`
	Metrics                 UsageMetrics `json:"metrics,omitempty"`
}

func PrepareSIMD(runDir, outDir string, width int) error {
	if width <= 0 || width > 64 {
		return fmt.Errorf("lane width must be 1..64")
	}
	gb, err := os.ReadFile(filepath.Join(runDir, "glyphs.json"))
	if err != nil {
		return err
	}
	var gt glyphcalc.GroundTruth
	if err = json.Unmarshal(gb, &gt); err != nil {
		return err
	}
	rb, err := os.ReadFile(filepath.Join(runDir, "render.json"))
	if err != nil {
		return err
	}
	var rm glyphcalc.RenderMeta
	if err = json.Unmarshal(rb, &rm); err != nil {
		return err
	}
	plan := SIMDPlan{Schema: "ohf.context-simd.plan.v1", Mode: "TUPLE4", LaneWidth: width, GlyphCount: len(gt.Glyphs), Columns: rm.Columns, SemanticMeaning: "NONE", Rule: "Each glyph is independent; emit shape,holes,direction,frames as four base-4 digits; ? means not uniquely readable."}
	var expected []SIMDRound
	for start, round := 0, 0; start < len(gt.Glyphs); start, round = start+width, round+1 {
		n := width
		if start+n > len(gt.Glyphs) {
			n = len(gt.Glyphs) - start
		}
		plan.Rounds = append(plan.Rounds, SIMDRound{Round: round, Start: start, Count: n})
		er := SIMDRound{Round: round, Start: start, Count: n}
		for i := 0; i < n; i++ {
			g := gt.Glyphs[start+i]
			er.Tuples = append(er.Tuples, fmt.Sprintf("%d%d%d%d", g.Shape, g.Holes, g.Direction, g.Frames))
		}
		expected = append(expected, er)
	}
	if err = os.MkdirAll(filepath.Join(outDir, "blind"), 0755); err != nil {
		return err
	}
	carrier, err := os.ReadFile(filepath.Join(runDir, "carrier.png"))
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(outDir, "blind", "carrier.png"), carrier, 0644); err != nil {
		return err
	}
	prompt := simdPrompt(plan)
	if err = os.WriteFile(filepath.Join(outDir, "blind", "MASTER_PROMPT.md"), []byte(prompt), 0644); err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(outDir, "blind", "query.txt"), []byte("Execute the supplied CONTEXT_SIMD schedule exactly. Return only the requested JSON.\n"), 0644); err != nil {
		return err
	}
	pb, _ := json.MarshalIndent(plan, "", "  ")
	pb = append(pb, '\n')
	if err = os.WriteFile(filepath.Join(outDir, "plan.json"), pb, 0644); err != nil {
		return err
	}
	exp := SIMDResponse{Schema: "ohf.context-simd.response.v1", Model: "EXPECTED_DO_NOT_SEND", Status: "KNOWN", Rounds: expected}
	eb, _ := json.MarshalIndent(exp, "", "  ")
	eb = append(eb, '\n')
	if err = os.WriteFile(filepath.Join(outDir, "expected.response.json"), eb, 0644); err != nil {
		return err
	}
	tmpl := SIMDResponse{Schema: "ohf.context-simd.response.v1", Model: "REPLACE_ME", Status: "UNKNOWN", Rounds: plan.Rounds}
	tb, _ := json.MarshalIndent(tmpl, "", "  ")
	tb = append(tb, '\n')
	return os.WriteFile(filepath.Join(outDir, "response.template.json"), tb, 0644)
}

func CheckSIMD(bundleDir, responsePath string) (SIMDCheck, error) {
	pb, err := os.ReadFile(filepath.Join(bundleDir, "plan.json"))
	if err != nil { return SIMDCheck{}, err }
	var plan SIMDPlan
	if err = json.Unmarshal(pb, &plan); err != nil { return SIMDCheck{}, err }
	eb, err := os.ReadFile(filepath.Join(bundleDir, "expected.response.json"))
	if err != nil { return SIMDCheck{}, err }
	var exp SIMDResponse
	if err = json.Unmarshal(eb, &exp); err != nil { return SIMDCheck{}, err }
	rb, err := os.ReadFile(responsePath)
	if err != nil { return SIMDCheck{}, err }
	var got SIMDResponse
	if err = json.Unmarshal(rb, &got); err != nil { return SIMDCheck{}, err }
	truth := map[int]string{}
	for _, r := range exp.Rounds { for i, t := range r.Tuples { truth[r.Start+i] = t } }
	seen := map[int]bool{}
	c := SIMDCheck{Schema: "ohf.context-simd.check.v1", Model: got.Model, LaneWidth: plan.LaneWidth, Total: plan.GlyphCount, Metrics: got.Metrics}
	var recovered []byte
	allExact := true
	for _, r := range got.Rounds {
		for i, t := range r.Tuples {
			idx := r.Start + i
			if idx < 0 || idx >= plan.GlyphCount || seen[idx] { c.Missing++; allExact = false; continue }
			seen[idx] = true
			t = strings.TrimSpace(t)
			if t == "?" || strings.EqualFold(t, "UNKNOWN") { c.Unknown++; allExact = false; continue }
			want := truth[idx]
			if t == want {
				c.ExactKnown++
				b, ok := tupleToByte(t)
				if !ok { allExact = false } else { recovered = append(recovered, b) }
			} else { c.FalseKnown++; allExact = false }
		}
	}
	for i := 0; i < plan.GlyphCount; i++ { if !seen[i] { c.Missing++; allExact = false } }
	c.FalseExactZero = c.FalseKnown == 0
	if allExact && len(recovered) == plan.GlyphCount {
		transport, err := os.ReadFile(filepath.Join(filepath.Dir(bundleDir), "transport.bin"))
		if err == nil && string(transport) == string(recovered) { c.ExactTransportRecovered = true }
	}
	out, _ := json.MarshalIndent(c, "", "  ")
	out = append(out, '\n')
	_ = os.WriteFile(filepath.Join(bundleDir, "check.json"), out, 0644)
	return c, nil
}

func tupleToByte(s string) (byte, bool) {
	if len(s) != 4 { return 0, false }
	v := make([]byte, 4)
	for i := 0; i < 4; i++ { if s[i] < '0' || s[i] > '3' { return 0, false }; v[i] = s[i]-'0' }
	return v[0] + 4*v[1] + 16*v[2] + 64*v[3], true
}

func simdPrompt(p SIMDPlan) string {
	var sched []string
	for _, r := range p.Rounds { sched = append(sched, fmt.Sprintf("ROUND %d: glyph indices %d..%d (%d independent lanes)", r.Round, r.Start, r.Start+r.Count-1, r.Count)) }
	return fmt.Sprintf(`# OHF CONTEXT_SIMD v1 — TUPLE4

Purpose: minimize serial reasoning. This is independent-lane batching, NOT a claim about hardware parallelism inside the model.

You receive one carrier image. Number glyphs row-major from 0. The grid has %d columns and %d glyphs total.

For EACH glyph independently observe only:
shape: circle=0 square=1 diamond=2 triangle=3
holes: 0..3
direction beacon: top=0 right=1 bottom=2 left=3
frame count: 0..3

Emit the tuple as four digits: shape holes direction frames.
Example: shape=2 holes=0 direction=3 frames=1 -> "2031".
If any state is not uniquely readable, emit "?" for that glyph. Never infer it from neighboring lanes.
Do NOT calculate bytes. Do NOT execute the payload. Do NOT assign semantic meaning.

MICRO-SCHEDULE
%s

Return ONLY JSON:
{"schema":"ohf.context-simd.response.v1","model":"MODEL_NAME","status":"KNOWN|PARTIAL|UNKNOWN","rounds":[{"round":0,"start":0,"count":%d,"tuples":["0000"]}],"metrics":{}}
`, p.Columns, p.GlyphCount, strings.Join(sched, "\n"), p.LaneWidth)
}
