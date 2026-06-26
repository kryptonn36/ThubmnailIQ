import cv2
import numpy as np

_face_cascade = cv2.CascadeClassifier(
    cv2.data.haarcascades + "haarcascade_frontalface_default.xml"
)
_smile_cascade = cv2.CascadeClassifier(
    cv2.data.haarcascades + "haarcascade_smile.xml"
)


def analyze_faces(image: np.ndarray) -> dict:
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    detections = _face_cascade.detectMultiScale(
        gray, scaleFactor=1.1, minNeighbors=5, minSize=(30, 30)
    )

    faces = []
    for (x, y, w, h) in detections:
        face_roi = gray[y : y + h, x : x + w]

        # No real emotion model (DeepFace) is used here for resource reasons.
        # We approximate "happy" by running a smile-detector cascade inside
        # the face bounding box; if it fires, call it happy, else neutral.
        smiles = _smile_cascade.detectMultiScale(
            face_roi, scaleFactor=1.7, minNeighbors=22, minSize=(int(w * 0.3), int(h * 0.15))
        )
        emotion = "happy" if len(smiles) > 0 else "neutral"

        faces.append(
            {
                "bbox": [int(x), int(y), int(x + w), int(y + h)],
                "dominant_emotion": emotion,
                "confidence": 0.92 if emotion == "happy" else 0.6,
            }
        )

    primary_emotion = "neutral"
    if faces:
        happy_count = sum(1 for f in faces if f["dominant_emotion"] == "happy")
        primary_emotion = "happy" if happy_count >= len(faces) / 2 else "neutral"

    return {
        "face_count": len(faces),
        # Stub: real eye-contact/gaze detection requires facial landmark
        # estimation which is out of scope for the lightweight pipeline.
        "has_eye_contact": False,
        "primary_emotion": primary_emotion,
        "faces": faces,
    }
