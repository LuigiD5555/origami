"""Trains MicroISADecoder on a JSONL manifest produced by
microisa_dataset.py. CPU-only by design (see the plan: the available GPU
has too little free VRAM to be worth the CUDA setup for a model this
small).

Usage:
    python3 microisa_train.py --train train.jsonl --val val.jsonl --out model.pt --epochs 15
"""
import argparse
import json

import numpy as np
import torch
import torch.nn as nn
from PIL import Image
from torch.utils.data import DataLoader, Dataset

from microisa_model import ATTRIBUTES, MicroISADecoder, count_parameters


class MicroISADataset(Dataset):
    def __init__(self, manifest_path: str) -> None:
        self.records = [json.loads(line) for line in open(manifest_path)]

    def __len__(self) -> int:
        return len(self.records)

    def __getitem__(self, idx: int):
        record = self.records[idx]
        image = Image.open(record["carrier_path"]).convert("RGB")
        tensor = torch.from_numpy(np.array(image, dtype="float32") / 255.0).permute(2, 0, 1)
        labels = {attr: record[attr] for attr in ATTRIBUTES}
        return tensor, labels


def collate(batch):
    images = torch.stack([item[0] for item in batch])
    labels = {attr: torch.tensor([item[1][attr] for item in batch], dtype=torch.long) for attr in ATTRIBUTES}
    return images, labels


def evaluate(model: nn.Module, loader: DataLoader) -> dict:
    model.eval()
    correct = {attr: 0 for attr in ATTRIBUTES}
    total = 0
    with torch.no_grad():
        for images, labels in loader:
            outputs = model(images)
            total += images.size(0)
            for attr in ATTRIBUTES:
                preds = outputs[attr].argmax(dim=1)
                correct[attr] += (preds == labels[attr]).sum().item()
    return {attr: correct[attr] / total for attr in ATTRIBUTES}


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--train", required=True)
    parser.add_argument("--val", required=True)
    parser.add_argument("--out", required=True)
    parser.add_argument("--metrics-out", default=None)
    parser.add_argument("--epochs", type=int, default=15)
    parser.add_argument("--batch-size", type=int, default=32)
    parser.add_argument("--lr", type=float, default=1e-3)
    args = parser.parse_args()

    train_loader = DataLoader(MicroISADataset(args.train), batch_size=args.batch_size, shuffle=True, collate_fn=collate)
    val_loader = DataLoader(MicroISADataset(args.val), batch_size=args.batch_size, shuffle=False, collate_fn=collate)

    model = MicroISADecoder()
    print(f"model parameters: {count_parameters(model):,}")

    optimizer = torch.optim.Adam(model.parameters(), lr=args.lr)
    criterion = nn.CrossEntropyLoss()

    for epoch in range(args.epochs):
        model.train()
        total_loss = 0.0
        for images, labels in train_loader:
            optimizer.zero_grad()
            outputs = model(images)
            loss = sum(criterion(outputs[attr], labels[attr]) for attr in ATTRIBUTES)
            loss.backward()
            optimizer.step()
            total_loss += loss.item()

        val_acc = evaluate(model, val_loader)
        print(f"epoch {epoch + 1}/{args.epochs} train_loss={total_loss / len(train_loader):.4f} val_acc={val_acc}")

    torch.save(model.state_dict(), args.out)
    print(f"saved checkpoint to {args.out}")

    if args.metrics_out:
        with open(args.metrics_out, "w") as f:
            json.dump({"val_accuracy": val_acc, "parameters": count_parameters(model)}, f, indent=2)

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
