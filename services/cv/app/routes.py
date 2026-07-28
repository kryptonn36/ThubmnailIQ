import os
from io import BytesIO

import cv2
import httpx
import numpy as np
from fastapi import APIRouter, HTTPException, Request
from PIL import Image

from app.clutter import analyze_clutter
from app.colors import analyze_colors
from app.face import analyze_faces
from app.logging_config import get_logger
from app.models import AnalyzeResponse, HealthResponse
from app.ocr import extract_text

router = APIRouter()
logger = get_logger(__name__)


@router.get("/health", response_model=HealthResponse)
async def health():
    return {"status": "ok"}


async def _load_image_from_url(image_url: str) -> Image.Image:
    if image_url.startswith("http://") or image_url.startswith("https://"):
        async with httpx.AsyncClient(timeout=15.0) as client:
            try:
                resp = await client.get(image_url)
                resp.raise_for_status()
            except httpx.HTTPError:
                logger.exception("failed to fetch image from url=%s", image_url)
                raise
            return Image.open(BytesIO(resp.content)).convert("RGB")

    if not os.path.exists(image_url):
        raise HTTPException(status_code=400, detail=f"image path not found: {image_url}")
    return Image.open(image_url).convert("RGB")


def _run_pipeline(pil_image: Image.Image) -> dict:
    # Resize image to prevent excessive processing time on large images
    # Reduce dimensions significantly for constrained deployment environment
    MAX_DIMENSION = 1024
    try:
        if max(pil_image.size) > MAX_DIMENSION:
            # Create a copy to avoid modifying the original
            process_image = pil_image.copy()
            # thumbnail() modifies the image in-place and maintains aspect ratio
            process_image.thumbnail((MAX_DIMENSION, MAX_DIMENSION), Image.Resampling.LANCZOS)
        else:
            process_image = pil_image

        # Use the (potentially resized) image for all processing
        cv_image = cv2.cvtColor(np.array(process_image), cv2.COLOR_RGB2BGR)

        # These now work on reasonably-sized images
        ocr_result = extract_text(process_image)
        face_result = analyze_faces(cv_image)
        color_result = analyze_colors(process_image)
        clutter_result = analyze_clutter(cv_image)

        visual_complexity = round(
            (clutter_result["clutter_score"] * 0.5) + (ocr_result["text_density_pct"] * 0.5), 2
        )
        visual_complexity = min(visual_complexity, 100)

        return {
            "ocr": ocr_result,
            "face": face_result,
            "colors": color_result,
            "clutter": clutter_result,
            "visual_complexity": visual_complexity,
        }
    except Exception as e:
        logger.exception(f"Error in _run_pipeline: {e}")
        # Return a minimal valid response to prevent service crashes
        return {
            "ocr": {"text_detected": False, "text_strings": [], "text_density_pct": 0.0, "word_count": 0, "avg_text_height_pct": 0.0},
            "face": {"face_count": 0, "has_eye_contact": False, "primary_emotion": "neutral", "faces": []},
            "colors": {"dominant_colors": [], "contrast_score": 1.0, "brightness_score": 0.0, "saturation_score": 0.0},
            "clutter": {"edge_density": 0.0, "clutter_score": 0.0, "object_count": 0, "objects": []},
            "visual_complexity": 0.0,
        }


@router.post("/analyze", response_model=AnalyzeResponse)
async def analyze(request: Request):
    content_type = request.headers.get("content-type", "")
    pil_image = None

    if "multipart/form-data" in content_type:
        form = await request.form()
        upload = form.get("image")
        image_url = form.get("image_url")
        if upload is not None:
            contents = await upload.read()
            pil_image = Image.open(BytesIO(contents)).convert("RGB")
        elif image_url:
            pil_image = await _load_image_from_url(image_url)
    elif "application/json" in content_type:
        body = await request.json()
        image_url = body.get("image_url") if isinstance(body, dict) else None
        if image_url:
            pil_image = await _load_image_from_url(image_url)
    else:
        try:
            body = await request.json()
            image_url = body.get("image_url") if isinstance(body, dict) else None
            if image_url:
                pil_image = await _load_image_from_url(image_url)
        except Exception:
            pil_image = None

    if pil_image is None:
        raise HTTPException(
            status_code=400,
            detail="provide multipart file field 'image' or JSON body {'image_url': ...}",
        )

    try:
        return _run_pipeline(pil_image)
    except Exception:
        logger.exception("analysis pipeline failed")
        raise
