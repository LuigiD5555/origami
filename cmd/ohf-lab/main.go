package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LuigiD5555/origami/internal/lab/artifact"
	"github.com/LuigiD5555/origami/internal/lab/baseline"
	"github.com/LuigiD5555/origami/internal/lab/canonical"
	"github.com/LuigiD5555/origami/internal/lab/dimensional"
	"github.com/LuigiD5555/origami/internal/lab/environment"
	"github.com/LuigiD5555/origami/internal/lab/evaluator"
	"github.com/LuigiD5555/origami/internal/lab/failure"
	"github.com/LuigiD5555/origami/internal/lab/groundtruth"
	"github.com/LuigiD5555/origami/internal/lab/nativebundle"
	"github.com/LuigiD5555/origami/internal/lab/observation"
	"github.com/LuigiD5555/origami/internal/lab/pipeline"
	"github.com/LuigiD5555/origami/internal/lab/prepare"
	"github.com/LuigiD5555/origami/internal/lab/runid"
	"github.com/LuigiD5555/origami/internal/lab/scorer"
	"github.com/LuigiD5555/origami/internal/lab/seed"
	"github.com/LuigiD5555/origami/internal/lab/spec"
	"github.com/LuigiD5555/origami/internal/lab/visualverify"
)

func main() {
	root, err := findRoot()
	if err != nil {
		fatal(err)
	}
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
		os.Exit(2)
	}

	switch args[0] {
	case "doctor":
		doctor(root)
	case "dimensional":
		dimensionalCmd(args[1:])
	case "baseline":
		baselineCmd(root, args[1:])
	case "experiment":
		experimentCmd(root, args[1:])
	case "seed":
		seedCmd(args[1:])
	case "run":
		runCmd(root, args[1:])
	case "evaluator":
		evaluatorCmd(args[1:])
	case "native":
		nativeCmd(root, args[1:])
	case "failure":
		failureCmd(root, args[1:])
	case "regression":
		regressionCmd(root, args[1:])
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println(`ohf-lab commands:
  doctor
  dimensional registry
  dimensional capacity [cells] [groups]
  dimensional validate-register <dvr.json>
  dimensional orthogonality <matrix.json>
  baseline verify
  baseline freeze --confirm-r310
  experiment validate <spec.json>
  experiment inspect-id <spec.json>
  experiment canonical <spec.json>
  experiment prepare <spec.json>
  experiment run <spec.json>
  run verify <run-id>
  run visual-verify <run-id>
  seed derive <master-uint64> <namespace>
  evaluator verify
  native bundle <run-id> <master-prompt.md> <query.txt>
  native verify <run-id> <bundle-id>
  native import <run-id> <observation.json>
  native score <run-id> <observation-id> <expectation.json>
  failure register <run-id> <score-id> [classification]
  regression verify <regression-id>`)
}

func doctor(root string) {
	v := baseline.Verify(root)
	env := environment.Capture(root)

	fmt.Println("OHF LAB DOCTOR")
	fmt.Println()
	fmt.Printf("Protocol             %s\n", env.ProtocolVersion)
	fmt.Printf("Baseline             %s\n", env.BaselineVersion)
	fmt.Printf("Go                   %s\n", env.GoVersion)
	fmt.Printf("OS/Arch              %s/%s\n", env.OS, env.Arch)
	fmt.Printf("Git commit           %s\n", env.GitCommit)
	fmt.Printf("Git state            %s\n", env.GitState)
	fmt.Printf("Dirty worktree       %t\n", env.DirtyWorktree)
	fmt.Printf("Baseline artifact    %s\n", v.Status)
	if len(v.Problems) > 0 {
		fmt.Printf("Baseline problems    %s\n", strings.Join(v.Problems, ", "))
	}
	fmt.Println()
	if v.Status == "PASS" {
		fmt.Println("READY = YES")
	} else {
		fmt.Println("READY = NO")
		os.Exit(1)
	}
}

func baselineCmd(root string, args []string) {
	if len(args) == 0 {
		usage()
		os.Exit(2)
	}
	switch args[0] {
	case "verify":
		v := baseline.Verify(root)
		fmt.Printf("STATUS=%s\n", v.Status)
		fmt.Printf("ARTIFACT=%s\n", v.ArtifactPath)
		if v.SHA256 != "" {
			fmt.Printf("SHA256=%s\n", v.SHA256)
			fmt.Printf("BYTES=%d\n", v.Bytes)
			fmt.Printf("ZIP_ENTRIES=%d\n", v.EntryCount)
		}
		for _, p := range v.Problems {
			fmt.Printf("PROBLEM=%s\n", p)
		}
		if err := v.Error(); err != nil {
			os.Exit(1)
		}
	case "freeze":
		if len(args) != 2 || args[1] != "--confirm-r310" {
			fatal(fmt.Errorf("baseline freeze requires explicit --confirm-r310 trust-boundary acknowledgement"))
		}
		l, err := baseline.Freeze(root)
		if err != nil {
			fatal(err)
		}
		fmt.Printf("FROZEN=R3.10\nSHA256=%s\nBYTES=%d\n", l.SHA256, l.Bytes)
		fmt.Println("NOTE=This freezes the supplied bytes; historical authenticity is operator-confirmed because the continuity record contains no original artifact hash.")
	default:
		usage()
		os.Exit(2)
	}
}

func experimentCmd(root string, args []string) {
	if len(args) < 2 {
		usage()
		os.Exit(2)
	}
	path := args[1]
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	s, err := spec.Load(path)
	if err != nil {
		fatal(err)
	}

	switch args[0] {
	case "validate":
		id, full, err := runid.ExperimentID(s)
		if err != nil {
			fatal(err)
		}
		fmt.Println("VALID=YES")
		fmt.Printf("EXPERIMENT_ID=%s\n", id)
		fmt.Printf("SPEC_SHA256=%s\n", full)
	case "inspect-id":
		id, full, err := runid.ExperimentID(s)
		if err != nil {
			fatal(err)
		}
		fmt.Printf("EXPERIMENT_ID=%s\nSPEC_SHA256=%s\n", id, full)
		for _, ns := range []string{"source", "layout", "motifs", "negatives", "degradation", "ordering", "sampling"} {
			fmt.Printf("SEED[%s]=%d\n", ns, seed.Derive(s.Seed.Value, ns))
		}
	case "prepare":
		r, err := prepare.Run(root, s)
		if err != nil {
			fatal(err)
		}
		fmt.Printf("EXPERIMENT_ID=%s\nSPEC_SHA256=%s\nRUN_ID=%s\nRUN_SHA256=%s\nSOURCE_SHA256=%s\nRUN_DIR=%s\nSTATUS=%s\n", r.ExperimentID, r.SpecSHA256, r.RunID, r.RunSHA256, r.SourceSHA256, r.RunDir, r.Status)
	case "run":
		r, err := pipeline.Run(root, s)
		if err != nil {
			fatal(err)
		}
		fmt.Printf("EXPERIMENT_ID=%s\nSPEC_SHA256=%s\nRUN_ID=%s\nRUN_SHA256=%s\nSOURCE_SHA256=%s\nPIXEL_SHA256=%s\nPNG_SHA256=%s\nPNG_BYTES=%d\nRUN_DIR=%s\nSTATUS=%s\n", r.ExperimentID, r.SpecSHA256, r.RunID, r.RunSHA256, r.SourceSHA256, r.PixelSHA256, r.PNGSHA256, r.PNGBytes, r.RunDir, r.Status)
	case "canonical":
		b, err := canonical.SpecBytes(s)
		if err != nil {
			fatal(err)
		}
		var pretty any
		if err := json.Unmarshal(b, &pretty); err != nil {
			fatal(err)
		}
		// Print exact canonical bytes, not pretty JSON.
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
	default:
		usage()
		os.Exit(2)
	}
}

func runCmd(root string, args []string) {
	if len(args) != 2 {
		usage()
		os.Exit(2)
	}
	runID := args[1]
	if strings.Contains(runID, "/") || strings.Contains(runID, "\\") || runID == "" {
		fatal(fmt.Errorf("invalid run id %q", runID))
	}
	dir := filepath.Join(root, "runs", runID)

	switch args[0] {
	case "verify":
		if err := artifact.Verify(dir); err != nil {
			fatal(err)
		}
		truthBytes, err := os.ReadFile(filepath.Join(dir, "ground_truth.json"))
		if err != nil {
			fatal(err)
		}
		var truth groundtruth.GroundTruth
		if err := json.Unmarshal(truthBytes, &truth); err != nil {
			fatal(err)
		}
		sourceBytes, err := os.ReadFile(filepath.Join(dir, "source", "source.bin"))
		if err != nil {
			fatal(err)
		}
		if err := groundtruth.Verify(truth, sourceBytes); err != nil {
			fatal(err)
		}
		fmt.Printf("RUN_ID=%s\nARTIFACT_MANIFEST=PASS\nGROUND_TRUTH=PASS\nSOURCE_SHA256=%s\n", runID, truth.Source.SHA256)
		if _, err := os.Stat(filepath.Join(dir, "render.json")); err == nil {
			v, err := visualverify.Verify(dir)
			if err != nil {
				fatal(err)
			}
			fmt.Printf("VISUAL_REOPEN=PASS\nPIXEL_SHA256=%s\nPNG_SHA256=%s\nPNG_BYTES=%d\n", v.PixelSHA256, v.PNGSHA256, v.PNGBytes)
		}
	case "visual-verify":
		v, err := visualverify.Verify(dir)
		if err != nil {
			fatal(err)
		}
		fmt.Printf("RUN_ID=%s\nVISUAL_REOPEN=PASS\nWIDTH=%d\nHEIGHT=%d\nPIXEL_SHA256=%s\nPNG_SHA256=%s\nPNG_BYTES=%d\n", runID, v.Width, v.Height, v.PixelSHA256, v.PNGSHA256, v.PNGBytes)
	default:
		usage()
		os.Exit(2)
	}
}

func evaluatorCmd(args []string) {
	if len(args) != 1 || args[0] != "verify" {
		usage()
		os.Exit(2)
	}
	if err := evaluator.ValidateIsolation(); err != nil {
		fatal(err)
	}
	c, _ := evaluator.For(evaluator.Computational)
	n, _ := evaluator.For(evaluator.Native)
	fmt.Printf("EVALUATOR_ISOLATION=PASS\nCOMPUTATIONAL_GROUND_TRUTH=%t\nNATIVE_GROUND_TRUTH=%t\nNATIVE_INPUTS=MASTER_PROMPT+ONE_IMAGE+QUERY\n", c.GroundTruth, n.GroundTruth)
}

func nativeCmd(root string, args []string) {
	if len(args) == 0 {
		usage()
		os.Exit(2)
	}
	switch args[0] {
	case "bundle":
		if len(args) != 4 {
			usage()
			os.Exit(2)
		}
		prompt := args[2]
		query := args[3]
		if !filepath.IsAbs(prompt) {
			prompt = filepath.Join(root, prompt)
		}
		if !filepath.IsAbs(query) {
			query = filepath.Join(root, query)
		}
		b, err := nativebundle.Build(root, args[1], prompt, query)
		if err != nil {
			fatal(err)
		}
		fmt.Printf("BUNDLE_ID=%s\nRUN_ID=%s\nBLIND_PATH=%s\nBLIND_FILES=3\nISOLATION=PASS\n", b.BundleID, b.RunID, b.BlindRelativePath)
	case "verify":
		if len(args) != 3 {
			usage()
			os.Exit(2)
		}
		b, err := nativebundle.Verify(root, args[1], args[2])
		if err != nil {
			fatal(err)
		}
		fmt.Printf("BUNDLE_ID=%s\nBLIND_FILES=%d\nISOLATION=PASS\n", b.BundleID, len(b.Files))
	case "import":
		if len(args) != 3 {
			usage()
			os.Exit(2)
		}
		path := args[2]
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, path)
		}
		id, out, err := observation.Import(root, args[1], path)
		if err != nil {
			fatal(err)
		}
		fmt.Printf("OBSERVATION_ID=%s\nPATH=%s\n", id, out)
	case "score":
		if len(args) != 4 {
			usage()
			os.Exit(2)
		}
		o, _, err := observation.Find(root, args[1], args[2])
		if err != nil {
			fatal(err)
		}
		expPath := args[3]
		if !filepath.IsAbs(expPath) {
			expPath = filepath.Join(root, expPath)
		}
		e, err := scorer.LoadExpectation(expPath)
		if err != nil {
			fatal(err)
		}
		sc := scorer.Evaluate(args[1], args[2], o, e)
		p, err := scorer.Write(root, args[1], sc)
		if err != nil {
			fatal(err)
		}
		fmt.Printf("SCORE_ID=%s\nVERDICT=%s\nFAILURE_KIND=%s\nMATCHED=%d\nTOTAL=%d\nPATH=%s\n", sc.ScoreID, sc.Verdict, sc.FailureKind, sc.Matched, sc.Total, p)
		if sc.Verdict != "PASS" {
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(2)
	}
}

func failureCmd(root string, args []string) {
	if len(args) < 3 || len(args) > 4 || args[0] != "register" {
		usage()
		os.Exit(2)
	}
	class := ""
	if len(args) == 4 {
		class = args[3]
	}
	r, err := failure.Register(root, args[1], args[2], class)
	if err != nil {
		fatal(err)
	}
	fmt.Printf("REGRESSION_ID=%s\nCLASSIFICATION=%s\nEVIDENCE_FILES=%d\nSTATUS=REGISTERED\n", r.RegressionID, r.Classification, len(r.SourceFiles))
}

func regressionCmd(root string, args []string) {
	if len(args) != 2 || args[0] != "verify" {
		usage()
		os.Exit(2)
	}
	if err := failure.Verify(root, args[1]); err != nil {
		fatal(err)
	}
	fmt.Printf("REGRESSION_ID=%s\nREPLAY_MANIFEST=PASS\n", args[1])
}

func dimensionalCmd(args []string) {
	if len(args) < 1 {
		usage()
		os.Exit(2)
	}
	r := dimensional.DefaultRegistry()
	switch args[0] {
	case "registry":
		if err := dimensional.ValidateRegistry(r); err != nil {
			fatal(err)
		}
		b, err := dimensional.MarshalCanonical(r)
		if err != nil {
			fatal(err)
		}
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
	case "capacity":
		cells, groups := 1, 1
		if len(args) >= 2 {
			if _, err := fmt.Sscan(args[1], &cells); err != nil {
				fatal(fmt.Errorf("invalid cells: %w", err))
			}
		}
		if len(args) == 3 {
			if _, err := fmt.Sscan(args[2], &groups); err != nil {
				fatal(fmt.Errorf("invalid groups: %w", err))
			}
		} else if len(args) > 3 {
			usage()
			os.Exit(2)
		}
		report, err := dimensional.Capacity(r, dimensional.OrthogonalityMatrix{Schema: "ohf.orthogonality.v1"}, cells, groups)
		if err != nil {
			fatal(err)
		}
		b, err := dimensional.MarshalCanonical(report)
		if err != nil {
			fatal(err)
		}
		os.Stdout.Write(b)
		os.Stdout.Write([]byte("\n"))
	case "validate-register":
		if len(args) != 2 {
			usage()
			os.Exit(2)
		}
		b, err := os.ReadFile(args[1])
		if err != nil {
			fatal(err)
		}
		var reg dimensional.Register
		if err := json.Unmarshal(b, &reg); err != nil {
			fatal(err)
		}
		if err := dimensional.ValidateRegister(r, reg); err != nil {
			fatal(err)
		}
		fmt.Println("DVR_VALID=PASS")
	case "orthogonality":
		if len(args) != 2 {
			usage()
			os.Exit(2)
		}
		b, err := os.ReadFile(args[1])
		if err != nil {
			fatal(err)
		}
		var matrix dimensional.OrthogonalityMatrix
		if err := json.Unmarshal(b, &matrix); err != nil {
			fatal(err)
		}
		report, err := dimensional.Capacity(r, matrix, 1, 1)
		if err != nil {
			fatal(err)
		}
		out, err := dimensional.MarshalCanonical(report)
		if err != nil {
			fatal(err)
		}
		os.Stdout.Write(out)
		os.Stdout.Write([]byte("\n"))
	default:
		usage()
		os.Exit(2)
	}
}

func seedCmd(args []string) {
	if len(args) != 3 || args[0] != "derive" {
		usage()
		os.Exit(2)
	}
	var master uint64
	if _, err := fmt.Sscan(args[1], &master); err != nil {
		fatal(fmt.Errorf("invalid master seed: %w", err))
	}
	fmt.Println(seed.Explain(master, args[2]))
}

func findRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("could not locate repository root (go.mod)")
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "ERROR:", err)
	os.Exit(1)
}
