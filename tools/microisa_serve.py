"""Resident inference service for microisa-cnn-r0. Speaks the
tlaloque.CapabilityRequest/CapabilityResponse HTTP_JSON contract (see
tlaloc/behavior-lab/internal/tlaloque/http_worker.go) unmodified — the Go
client side (internal/tlaloque/microisadecoder) treats this exactly like
cmd/tlaloc-embedding-scorer: a resident micro-model service, loaded once,
called many times.

Usage:
    python3 microisa_serve.py --checkpoint ../models/microisa-cnn-r0.pt --addr :8791
"""
import argparse
import base64
import io
import json
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

import torch
from PIL import Image

from microisa_model import ATTRIBUTES, MicroISADecoder, count_parameters, predict_image

WORKER_ID = "microisa-cnn-r0"


def make_handler(model: MicroISADecoder, checkpoint_path: str):
    class Handler(BaseHTTPRequestHandler):
        def _send_json(self, status: int, payload: dict) -> None:
            body = json.dumps(payload).encode("utf-8")
            self.send_response(status)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(body)))
            self.end_headers()
            self.wfile.write(body)

        def do_GET(self) -> None:
            if self.path != "/health":
                self._send_json(404, {"error": "not found"})
                return
            self._send_json(200, {
                "status": "ok",
                "worker_id": WORKER_ID,
                "checkpoint": checkpoint_path,
                "parameters": count_parameters(model),
                "attributes": list(ATTRIBUTES),
            })

        def do_POST(self) -> None:
            if self.path != "/execute":
                self._send_json(404, {"error": "not found"})
                return

            length = int(self.headers.get("Content-Length", "0"))
            try:
                envelope = json.loads(self.rfile.read(length))
                carrier_b64 = envelope["input"]["carrier_png_base64"]
                image = Image.open(io.BytesIO(base64.b64decode(carrier_b64)))
            except Exception as exc:  # noqa: BLE001 - reported to the caller, not swallowed
                self._send_json(400, {"error": f"invalid request: {exc}"})
                return

            prediction = predict_image(model, image)
            self._send_json(200, {
                "worker_id": WORKER_ID,
                "output": prediction,
                "confidence": 1.0,
            })

        def log_message(self, format: str, *fmt_args) -> None:  # noqa: A002 - stdlib signature
            print(f"[microisa_serve] {self.address_string()} {format % fmt_args}")

    return Handler


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--checkpoint", required=True)
    parser.add_argument("--addr", default=":8791", help="host:port, or :port for all interfaces")
    args = parser.parse_args()

    model = MicroISADecoder()
    model.load_state_dict(torch.load(args.checkpoint))
    model.eval()
    print(f"loaded {args.checkpoint} ({count_parameters(model):,} parameters)")

    host, _, port = args.addr.rpartition(":")
    server = ThreadingHTTPServer((host, int(port)), make_handler(model, args.checkpoint))
    print(f"microisa_serve listening on {args.addr}")
    server.serve_forever()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
