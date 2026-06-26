# ThumbnailIQ CV Service

Lightweight FastAPI microservice that performs computer-vision analysis
(OCR, face/emotion approximation, color analysis, clutter/edge density) on
YouTube thumbnail images. Designed to be called over HTTP by the Go backend
worker.

Runs on port **8001**.

## Run locally with a venv

```bash
cd services/cv
python3 -m venv .venv
.venv/bin/pip install -r requirements.txt

# OCR requires the tesseract binary on the host (see "Known simplifications"
# below for what happens if it's missing).
sudo apt-get install -y tesseract-ocr   # Debian/Ubuntu/Kali

.venv/bin/uvicorn app.main:app --host 0.0.0.0 --port 8001
```

## Run via Docker

```bash
cd services/cv
docker build -t thumbnailiq-cv .
docker run --rm -p 8001:8001 thumbnailiq-cv
```

The Dockerfile installs `tesseract-ocr`, `libgl1`, and `libglib2.0-0` (the
last two are required by `opencv-python-headless` at import time even though
it's the "headless" build).

## API

### `GET /health`

```bash
curl -s http://localhost:8001/health
# {"status":"ok"}
```

### `POST /analyze`

Accepts either a multipart file upload (field name `image`) or a JSON body
with `image_url` (which may be an `http(s)://` URL or a local file path
readable by the service process).

```bash
# multipart file upload
curl -s -X POST http://localhost:8001/analyze -F "image=@/tmp/test_thumb.jpg"

# JSON body with a local path
curl -s -X POST http://localhost:8001/analyze \
  -H "Content-Type: application/json" \
  -d '{"image_url": "/tmp/test_thumb.jpg"}'

# JSON body with a remote URL
curl -s -X POST http://localhost:8001/analyze \
  -H "Content-Type: application/json" \
  -d '{"image_url": "https://example.com/thumb.jpg"}'
```

Response shape:

```json
{
  "ocr": {
    "text_detected": true,
    "text_strings": ["STRING ONE", "string two"],
    "text_density_pct": 12.5,
    "word_count": 4,
    "avg_text_height_pct": 9.2
  },
  "face": {
    "face_count": 1,
    "has_eye_contact": false,
    "primary_emotion": "happy",
    "faces": [
      {"bbox": [0, 0, 100, 100], "dominant_emotion": "happy", "confidence": 0.92}
    ]
  },
  "colors": {
    "dominant_colors": [
      {"hex": "#1a2b3c", "rgb": [26, 43, 60], "percentage": 34.2, "luminance": 41.1, "saturation": 28.5}
    ],
    "contrast_score": 4.5,
    "brightness_score": 120.3,
    "saturation_score": 45.2
  },
  "clutter": {
    "edge_density": 0.31,
    "clutter_score": 42.0,
    "object_count": 0,
    "objects": []
  },
  "visual_complexity": 38.7
}
```

## Known simplifications vs. original spec

The original blueprint (section 8, "AI & Computer Vision Pipeline") calls for
EasyOCR, InsightFace, DeepFace, and YOLOv8. Those require multi-GB model
downloads and heavy dependencies (PyTorch, etc.), which is impractical for
this environment. This service substitutes lightweight equivalents instead:

- **OCR**: `pytesseract` (wraps the system `tesseract` binary) instead of
  EasyOCR. If the `tesseract` binary is not installed/available on the host,
  `extract_text()` catches the resulting exception and returns
  `text_detected: false` with empty results rather than crashing the
  `/analyze` request. In the development sandbox used to build this service,
  the `tesseract-ocr` apt package could not be installed (no root/sudo
  access), so the graceful-degradation path is what was actually exercised
  during verification — install `tesseract-ocr` to get real OCR results.
- **Face detection**: OpenCV Haar Cascades (`haarcascade_frontalface_default.xml`,
  bundled with `opencv-python-headless`, no download needed) instead of
  InsightFace. Haar cascades are far less accurate than a modern face
  detector, especially on stylized/illustrated thumbnails.
- **Emotion**: there is no real emotion model (DeepFace is skipped entirely).
  `primary_emotion`/`dominant_emotion` are approximated by running a smile
  Haar cascade (`haarcascade_smile.xml`) inside each detected face's bounding
  box: a detected smile region maps to `"happy"`, otherwise `"neutral"`. This
  is a coarse heuristic, not a real emotion classifier.
- **Eye contact**: `has_eye_contact` is hardcoded to `false` always. Real
  gaze/eye-contact estimation requires facial landmark detection which is out
  of scope for this lightweight pipeline.
- **Object detection (YOLOv8)**: skipped entirely for resource reasons (no
  model download, no torch dependency). `clutter.object_count` is hardcoded
  to `0` and `clutter.objects` to `[]`. The `clutter_score` is based purely on
  Sobel edge density (no object-count component, since there's no detector).
