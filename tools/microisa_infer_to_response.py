"""Runs the trained MicroISADecoder against real held-out microbench cases
and writes response.json in the exact schema each op's MASTER_PROMPT.md
demands, so the existing Go harness (microbench-score) can grade it with
the same yardstick used for claude-opus-5 — no separate scoring logic here.

Usage:
    python3 microisa_infer_to_response.py --checkpoint model.pt --root ../runs/microbench --model-name microisa-cnn-r0
"""
import argparse
import json
from pathlib import Path

import torch
from PIL import Image

from microisa_model import MicroISADecoder, predict_image

OPS = ("observe4", "tuple4", "decode8")


def predict(model: MicroISADecoder, carrier_path: Path) -> dict:
    return predict_image(model, Image.open(carrier_path))


def format_lane(op: str, i: int, prediction: dict) -> dict:
    if op == "observe4":
        return {"i": i, "status": "KNOWN", **prediction}
    if op == "tuple4":
        tuple4 = "".join(str(prediction[a]) for a in ("shape", "holes", "direction", "frames"))
        return {"i": i, "status": "KNOWN", "tuple4": tuple4}
    if op == "decode8":
        byte_value = prediction["shape"] + 4 * prediction["holes"] + 16 * prediction["direction"] + 64 * prediction["frames"]
        return {"i": i, "status": "KNOWN", "byte_hex": f"{byte_value:02x}"}
    raise ValueError(f"unknown op {op}")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--checkpoint", required=True)
    parser.add_argument("--root", required=True, help="microbench root, e.g. ../runs/microbench")
    parser.add_argument("--model-name", default="microisa-cnn-r0")
    args = parser.parse_args()

    model = MicroISADecoder()
    model.load_state_dict(torch.load(args.checkpoint))
    model.eval()

    root = Path(args.root)
    written = []
    for op in OPS:
        for case_dir in sorted((root / op / "w01").glob("c*")):
            carrier_path = case_dir / "blind" / "carrier.png"
            if not carrier_path.exists():
                continue
            prediction = predict(model, carrier_path)
            response = {
                "schema": "ohf.microisa.response.v1",
                "model": args.model_name,
                "status": "KNOWN",
                "lanes": [format_lane(op, 0, prediction)],
                "metrics": {},
            }
            response_path = case_dir / "response.json"
            response_path.write_text(json.dumps(response, indent=2))
            written.append(str(response_path))

    print(f"wrote {len(written)} response.json files")
    for path in written:
        print(f"  {path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
