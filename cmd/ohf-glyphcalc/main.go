package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/LuigiD5555/origami/internal/lab/glyphcalc"
	"github.com/LuigiD5555/origami/internal/lab/microisa"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "demo":
		demo(os.Args[2:])
	case "verify":
		verify(os.Args[2:])
	case "glyph":
		glyph(os.Args[2:])
	case "native-check":
		nativeCheck(os.Args[2:])
	case "microbench-prepare":
		microbenchPrepare(os.Args[2:])
	case "microbench-score":
		microbenchScore(os.Args[2:])
	case "microbench-profile":
		microbenchProfile(os.Args[2:])
	case "microbench-next":
		microbenchNext(os.Args[2:])
	case "simd-prepare":
		simdPrepare(os.Args[2:])
	case "simd-check":
		simdCheck(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println(`ohf-glyphcalc commands:
  demo --pattern mixed --length 4096 --seed 817392 --out runs/glyphcalc-mixed
  verify --run runs/glyphcalc-mixed
  glyph 0xAA
  native-check --run runs/glyphcalc-mixed --response response.json
  microbench-prepare --out runs/microbench --widths 1,4,8,16 --cases 2 --seed 817392
  microbench-score --case runs/microbench/observe4/w04/c00 --response response.json
  microbench-profile --root runs/microbench --min-exact 0.98
  microbench-next --root runs/microbench --min-exact 0.98
  simd-prepare --run runs/glyphcalc-native-small --width 8 --out runs/glyphcalc-native-small/context-simd-w08
  simd-check --bundle runs/glyphcalc-native-small/context-simd-w08 --response response.json`)
}

func demo(args []string) {
	fs := flag.NewFlagSet("demo", flag.ExitOnError)
	pattern := fs.String("pattern", "mixed", "runs|sequence|motif|mixed|random")
	length := fs.Int("length", 4096, "source bytes <= 65535")
	seed := fs.Uint64("seed", 817392, "deterministic seed")
	out := fs.String("out", "runs/glyphcalc-demo", "output directory")
	_ = fs.Parse(args)
	dir := abs(*out)
	m, err := glyphcalc.Run(dir, *pattern, *length, *seed)
	fatalIf(err)
	fmt.Printf("STATUS=PASS\nRUN_DIR=%s\nMODE=%s\nSOURCE_BYTES=%d\nGLYPHS=%d\nSOURCE_BYTES_PER_GLYPH=%.3f\nPERCEPTION_REDUCTION_PCT=%.2f\nEXACT_RECONSTRUCTION=%t\nCARRIER_PNG_BYTES=%d\n", dir, m.SelectedMode, m.SourceBytes, m.GlyphCount, m.SourceBytesPerGlyph, m.PerceptionReductionPct, m.ExactReconstruction, m.CarrierPNGBytes)
}

func verify(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	run := fs.String("run", "", "run directory")
	_ = fs.Parse(args)
	if *run == "" {
		fatalIf(fmt.Errorf("--run required"))
	}
	dir := abs(*run)
	fatalIf(glyphcalc.Verify(dir))
	fmt.Printf("STATUS=PASS\nRUN_DIR=%s\nCHECKSUMS=PASS\nRECONSTRUCTION=PASS\n", dir)
}

func glyph(args []string) {
	if len(args) != 1 {
		usage()
		os.Exit(2)
	}
	s := strings.TrimPrefix(strings.ToLower(args[0]), "0x")
	v, err := strconv.ParseUint(s, 16, 8)
	fatalIf(err)
	g := glyphcalc.EncodeByte(0, byte(v))
	b, _ := json.MarshalIndent(g, "", "  ")
	fmt.Println(string(b))
	fmt.Printf("FORMULA=%d + 4*%d + 16*%d + 64*%d = %d\n", g.Shape, g.Holes, g.Direction, g.Frames, g.Byte)
}

func nativeCheck(args []string) {
	fs := flag.NewFlagSet("native-check", flag.ExitOnError)
	run := fs.String("run", "", "run directory")
	resp := fs.String("response", "", "VLM JSON response")
	_ = fs.Parse(args)
	if *run == "" || *resp == "" {
		fatalIf(fmt.Errorf("--run and --response required"))
	}
	c, err := glyphcalc.CheckNative(abs(*run), abs(*resp))
	fatalIf(err)
	b, _ := json.MarshalIndent(c, "", "  ")
	fmt.Println(string(b))
	if !c.PerceptionExact || !c.ExactReconstruction {
		os.Exit(1)
	}
}

func abs(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	wd, _ := os.Getwd()
	return filepath.Join(wd, p)
}

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
}

func microbenchPrepare(args []string) {
	fs := flag.NewFlagSet("microbench-prepare", flag.ExitOnError)
	out := fs.String("out", "runs/microbench", "campaign output")
	widthText := fs.String("widths", "1,4,8,16", "comma separated lane widths")
	cases := fs.Int("cases", 2, "cases per width/op")
	seed := fs.Uint64("seed", 817392, "master seed")
	_ = fs.Parse(args)
	var widths []int
	for _, x := range strings.Split(*widthText, ",") {
		v, err := strconv.Atoi(strings.TrimSpace(x))
		fatalIf(err)
		widths = append(widths, v)
	}
	fatalIf(microisa.PrepareCampaign(abs(*out), widths, *cases, *seed))
	fmt.Printf("STATUS=PASS\nROOT=%s\nWIDTHS=%v\nCASES=%d\nEXTERNAL_REQUESTS_MADE=0\n", abs(*out), widths, *cases)
}

func microbenchScore(args []string) {
	fs := flag.NewFlagSet("microbench-score", flag.ExitOnError)
	caseDir := fs.String("case", "", "probe case directory")
	resp := fs.String("response", "", "model response JSON")
	_ = fs.Parse(args)
	if *caseDir == "" || *resp == "" {
		fatalIf(fmt.Errorf("--case and --response required"))
	}
	s, err := microisa.ScoreResponse(abs(*caseDir), abs(*resp))
	fatalIf(err)
	b, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(b))
	if !s.SafeFalseExact {
		os.Exit(1)
	}
}

func microbenchProfile(args []string) {
	fs := flag.NewFlagSet("microbench-profile", flag.ExitOnError)
	root := fs.String("root", "runs/microbench", "campaign root")
	minExact := fs.Float64("min-exact", 0.98, "minimum exact-known lane rate")
	_ = fs.Parse(args)
	p, err := microisa.BuildProfile(abs(*root), *minExact)
	fatalIf(err)
	b, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(b))
}

func simdPrepare(args []string) {
	fs := flag.NewFlagSet("simd-prepare", flag.ExitOnError)
	run := fs.String("run", "", "glyphcalc run dir")
	width := fs.Int("width", 8, "logical lane width")
	out := fs.String("out", "", "output bundle dir")
	_ = fs.Parse(args)
	if *run == "" || *out == "" {
		fatalIf(fmt.Errorf("--run and --out required"))
	}
	fatalIf(microisa.PrepareSIMD(abs(*run), abs(*out), *width))
	fmt.Printf("STATUS=PASS\nBUNDLE=%s\nLANE_WIDTH=%d\nEXTERNAL_REQUESTS_MADE=0\n", abs(*out), *width)
}

func simdCheck(args []string) {
	fs := flag.NewFlagSet("simd-check", flag.ExitOnError)
	bundle := fs.String("bundle", "", "context SIMD bundle")
	resp := fs.String("response", "", "model response JSON")
	_ = fs.Parse(args)
	if *bundle == "" || *resp == "" {
		fatalIf(fmt.Errorf("--bundle and --response required"))
	}
	c, err := microisa.CheckSIMD(abs(*bundle), abs(*resp))
	fatalIf(err)
	b, _ := json.MarshalIndent(c, "", "  ")
	fmt.Println(string(b))
	if !c.FalseExactZero {
		os.Exit(1)
	}
}

func microbenchNext(args []string) {
	fs := flag.NewFlagSet("microbench-next", flag.ExitOnError)
	root := fs.String("root", "runs/microbench", "campaign root")
	minExact := fs.Float64("min-exact", 0.98, "minimum exact-known rate")
	_ = fs.Parse(args)
	n, err := microisa.RecommendNext(abs(*root), *minExact)
	fatalIf(err)
	b, _ := json.MarshalIndent(n, "", "  ")
	fmt.Println(string(b))
}
