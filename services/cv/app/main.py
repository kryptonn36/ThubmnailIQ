from fastapi import FastAPI

from app.logging_config import configure_logging
from app.middleware import RequestLoggingMiddleware
from app.routes import router

configure_logging()

app = FastAPI(title="ThumbnailIQ CV Service", version="0.1.0")
app.add_middleware(RequestLoggingMiddleware)
app.include_router(router)
