from pydantic import BaseModel
from typing import Optional


class QueryData(BaseModel):
    match_name: Optional[str] = None
    match_source: Optional[str] = None
    match_doctype: Optional[str] = None
    match_stage: Optional[str] = None
    match_metadata: Optional[dict] = None
