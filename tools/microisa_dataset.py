"""Generates labeled SAFE_MICRO_ISA training data by shelling out to the
real Go generator (cmd/ohf-glyphcalc microbench-prepare) — it is not
reimplemented here. Only the DECODE8 op's cases are indexed: the four
underlying glyph attributes (shape/holes/direction/frames) are what a model
needs to learn regardless of which micro-op eventually reports them, so one
op's images are enough training data. Splits (train/val/test) use disjoint
--seed values, one generator invocation each, so a single call produces
`--cases N` deterministically distinct examples in one shot (case_index
0..N-1 are already derived independently inside the generator — no need to
invoke it once per example).

Usage:
    python3 microisa_dataset.py --seed 1001 --cases 4000 --out /tmp/microisa-train --manifest /tmp/microisa-train.jsonl
"""
import argparse
import json
import subprocess
from pathlib import Path

ORIGAMI_ROOT = Path(__file__).resolve().parent.parent


def generate(out_dir: Path, seed: int, cases: int) -> None:
    """Runs the real Go generator once, producing `cases` examples."""
    subprocess.run(
        [
            "go", "run", "./cmd/ohf-glyphcalc", "microbench-prepare",
            "--out", str(out_dir),
            "--widths", "1",
            "--cases", str(cases),
            "--seed", str(seed),
        ],
        cwd=ORIGAMI_ROOT,
        check=True,
    )


def index(out_dir: Path, cases: int, manifest_path: Path) -> int:
    written = 0
    with manifest_path.open("w") as manifest:
        for c in range(cases):
            case_dir = out_dir / "decode8" / "w01" / f"c{c:02d}"
            probe_path = case_dir / "probe.json"
            carrier_path = case_dir / "blind" / "carrier.png"
            if not probe_path.exists() or not carrier_path.exists():
                continue
            truth = json.loads(probe_path.read_text())["truth"][0]
            manifest.write(json.dumps({
                "carrier_path": str(carrier_path),
                "shape": truth["shape"],
                "holes": truth["holes"],
                "direction": truth["direction"],
                "frames": truth["frames"],
            }) + "\n")
            written += 1
    return written


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--seed", type=int, required=True)
    parser.add_argument("--cases", type=int, required=True)
    parser.add_argument("--out", required=True, help="directory to hold the generated campaign")
    parser.add_argument("--manifest", required=True, help="output JSONL manifest path")
    args = parser.parse_args()

    out_dir = Path(args.out)
    manifest_path = Path(args.manifest)
    generate(out_dir, args.seed, args.cases)
    written = index(out_dir, args.cases, manifest_path)
    print(f"wrote {written}/{args.cases} examples to {manifest_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
