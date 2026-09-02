"""A tiny CNN, trained from scratch (no HuggingFace download), that reads a
108x108 SAFE_MICRO_ISA carrier glyph and predicts its four discrete
attributes (shape/holes/direction/frames, each a 4-way classification).
This is deliberately small: the task is local geometric perception on a
tiny synthetic image, not general visual understanding, so a few
convolutional layers are the right-sized tool — no attention, no
pretrained backbone.
"""
import numpy as np
import torch
import torch.nn as nn
from PIL import Image


class MicroISADecoder(nn.Module):
    def __init__(self) -> None:
        super().__init__()
        self.features = nn.Sequential(
            nn.Conv2d(3, 8, kernel_size=3, padding=1), nn.ReLU(), nn.MaxPool2d(2),   # 108 -> 54
            nn.Conv2d(8, 16, kernel_size=3, padding=1), nn.ReLU(), nn.MaxPool2d(2),  # 54 -> 27
            nn.Conv2d(16, 32, kernel_size=3, padding=1), nn.ReLU(), nn.AdaptiveAvgPool2d(6),  # -> 6x6
        )
        trunk_dim = 32 * 6 * 6
        self.shape_head = nn.Linear(trunk_dim, 4)
        self.holes_head = nn.Linear(trunk_dim, 4)
        self.direction_head = nn.Linear(trunk_dim, 4)
        self.frames_head = nn.Linear(trunk_dim, 4)

    def forward(self, x: torch.Tensor):
        trunk = self.features(x).flatten(1)
        return {
            "shape": self.shape_head(trunk),
            "holes": self.holes_head(trunk),
            "direction": self.direction_head(trunk),
            "frames": self.frames_head(trunk),
        }


ATTRIBUTES = ("shape", "holes", "direction", "frames")


def count_parameters(model: nn.Module) -> int:
    return sum(p.numel() for p in model.parameters() if p.requires_grad)


def image_to_tensor(image: Image.Image) -> torch.Tensor:
    """Converts a PIL image into the model's input tensor, in memory —
    shared by every caller (training, file-based inference, the HTTP
    service) so the preprocessing never drifts between them."""
    rgb = image.convert("RGB")
    tensor = torch.from_numpy(np.array(rgb, dtype="float32") / 255.0).permute(2, 0, 1)
    return tensor.unsqueeze(0)


def predict_image(model: "MicroISADecoder", image: Image.Image) -> dict:
    model.eval()
    with torch.no_grad():
        outputs = model(image_to_tensor(image))
    return {attr: int(outputs[attr].argmax(dim=1).item()) for attr in ATTRIBUTES}
