from PIL import Image
import pytesseract


def extract_text(image: Image.Image) -> dict:
    width, height = image.size
    image_area = width * height

    # pytesseract shells out to the `tesseract` binary; if it isn't installed
    # on the host/container, this raises TesseractNotFoundError. Degrade
    # gracefully instead of failing the whole /analyze request.
    try:
        data = pytesseract.image_to_data(image, output_type=pytesseract.Output.DICT)
    except Exception:
        return {
            "text_detected": False,
            "text_strings": [],
            "text_density_pct": 0.0,
            "word_count": 0,
            "avg_text_height_pct": 0.0,
        }

    text_strings = []
    total_box_area = 0
    heights_pct = []

    n = len(data.get("text", []))
    for i in range(n):
        raw = data["text"][i].strip()
        if not raw:
            continue
        try:
            conf = float(data["conf"][i])
        except (ValueError, TypeError):
            conf = -1.0
        if conf < 0:
            continue

        w = data["width"][i]
        h = data["height"][i]
        text_strings.append(raw)
        total_box_area += w * h
        if height > 0:
            heights_pct.append((h / height) * 100)

    word_count = len(text_strings)
    text_density_pct = (total_box_area / image_area) * 100 if image_area > 0 else 0.0
    avg_text_height_pct = sum(heights_pct) / len(heights_pct) if heights_pct else 0.0

    return {
        "text_detected": word_count > 0,
        "text_strings": text_strings,
        "text_density_pct": round(text_density_pct, 2),
        "word_count": word_count,
        "avg_text_height_pct": round(avg_text_height_pct, 2),
    }
