from pydantic import BaseModel
from typing import Optional


class QueryData(BaseModel):
    match_name: Optional[str] = None
    match_source: Optional[str] = None
    match_doctype: Optional[str] = None
    match_stage: Optional[str] = None


def querydata_to_filters(query: QueryData) -> dict:
    filters = {}
    if query.match_name is not None:
        filters["name"] = query.match_name
    if query.match_source is not None:
        filters["source"] = query.match_source
    if query.match_doctype is not None:
        filters["doctype"] = query.match_doctype
    if query.match_stage is not None:
        filters["stage"] = query.match_stage
    return filters
