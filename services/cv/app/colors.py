import colorsys

import numpy as np
from PIL import Image
from sklearn.cluster import KMeans


def _relative_luminance(rgb):
    r, g, b = [x / 255 for x in rgb]
    r = r / 12.92 if r <= 0.03928 else ((r + 0.055) / 1.055) ** 2.4
    g = g / 12.92 if g <= 0.03928 else ((g + 0.055) / 1.055) ** 2.4
    b = b / 12.92 if b <= 0.03928 else ((b + 0.055) / 1.055) ** 2.4
    return 0.2126 * r + 0.7152 * g + 0.0722 * b


def calculate_wcag_contrast(rgb1, rgb2) -> float:
    l1 = _relative_luminance(rgb1) + 0.05
    l2 = _relative_luminance(rgb2) + 0.05
    return max(l1, l2) / min(l1, l2)


def analyze_colors(image: Image.Image, n_colors: int = 5) -> dict:
    img = image.convert("RGB")
    img.thumbnail((150, 150))
    pixels = np.array(img).reshape(-1, 3)

    k = min(n_colors, len(np.unique(pixels, axis=0)))
    k = max(k, 1)
    kmeans = KMeans(n_clusters=k, n_init=4, random_state=42)
    labels = kmeans.fit_predict(pixels)
    centers = kmeans.cluster_centers_

    counts = np.bincount(labels, minlength=k)
    total = counts.sum()

    dominant_colors = []
    order = np.argsort(-counts)
    for idx in order:
        rgb = [int(round(c)) for c in centers[idx]]
        rgb = [max(0, min(255, c)) for c in rgb]
        pct = float(counts[idx] / total * 100)
        luminance = _relative_luminance(rgb) * 100
        r, g, b = [c / 255 for c in rgb]
        _, _, s = colorsys.rgb_to_hls(r, g, b)
        dominant_colors.append(
            {
                "hex": "#{:02x}{:02x}{:02x}".format(*rgb),
                "rgb": rgb,
                "percentage": round(pct, 2),
                "luminance": round(luminance, 2),
                "saturation": round(s * 100, 2),
            }
        )

    if len(dominant_colors) >= 2:
        contrast_score = calculate_wcag_contrast(
            dominant_colors[0]["rgb"], dominant_colors[1]["rgb"]
        )
    else:
        contrast_score = 1.0

    gray_img = np.array(image.convert("L"))
    brightness_score = float(np.mean(gray_img))

    hsv_img = np.array(image.convert("HSV"))
    saturation_score = float(np.mean(hsv_img[:, :, 1]) / 255 * 100)

    return {
        "dominant_colors": dominant_colors,
        "contrast_score": round(contrast_score, 2),
        "brightness_score": round(brightness_score, 2),
        "saturation_score": round(saturation_score, 2),
    }
