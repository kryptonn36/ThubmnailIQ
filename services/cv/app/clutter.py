import cv2
import numpy as np


def analyze_clutter(image: np.ndarray) -> dict:
    resized = cv2.resize(image, (480, 270))
    gray = cv2.cvtColor(resized, cv2.COLOR_BGR2GRAY)

    sobelx = cv2.Sobel(gray, cv2.CV_64F, 1, 0, ksize=3)
    sobely = cv2.Sobel(gray, cv2.CV_64F, 0, 1, ksize=3)
    magnitude = np.sqrt(sobelx**2 + sobely**2)

    edge_density = float(np.mean(magnitude) / 255)
    clutter_score = min(edge_density * 100, 100)

    return {
        "edge_density": round(edge_density, 4),
        "clutter_score": round(clutter_score, 2),
        # YOLOv8 object detection was skipped to avoid large model downloads
        # and heavy torch dependency; no object count/classes are produced.
        "object_count": 0,
        "objects": [],
    }
