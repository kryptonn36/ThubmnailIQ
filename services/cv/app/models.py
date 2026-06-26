from typing import List, Optional
from pydantic import BaseModel


class AnalyzeRequest(BaseModel):
    image_url: Optional[str] = None


class OCRResult(BaseModel):
    text_detected: bool
    text_strings: List[str]
    text_density_pct: float
    word_count: int
    avg_text_height_pct: float


class FaceBox(BaseModel):
    bbox: List[int]
    dominant_emotion: str
    confidence: float


class FaceResult(BaseModel):
    face_count: int
    has_eye_contact: bool
    primary_emotion: str
    faces: List[FaceBox]


class DominantColor(BaseModel):
    hex: str
    rgb: List[int]
    percentage: float
    luminance: float
    saturation: float


class ColorResult(BaseModel):
    dominant_colors: List[DominantColor]
    contrast_score: float
    brightness_score: float
    saturation_score: float


class ClutterResult(BaseModel):
    edge_density: float
    clutter_score: float
    object_count: int
    objects: List[str]


class AnalyzeResponse(BaseModel):
    ocr: OCRResult
    face: FaceResult
    colors: ColorResult
    clutter: ClutterResult
    visual_complexity: float


class HealthResponse(BaseModel):
    status: str
