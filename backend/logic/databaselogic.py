from pydantic import BaseModel
from typing import Optional, Any
from advanced_alchemy.filters import SearchFilter, CollectionFilter


class QueryData(BaseModel):
    match_name: Optional[str] = None
    match_source: Optional[str] = None
    match_doctype: Optional[str] = None
    match_stage: Optional[str] = None


def querydata_to_filters(query: QueryData) -> Any:
    filters = []
    if query.match_name is not None:
        filters.append(SearchFilter(field_name="name", value=query.match_name))
    if query.match_source is not None:
        filters.append(
            CollectionFilter(field_name="source", values=[query.match_source])
        )
    if query.match_doctype is not None:
        filters.append(
            CollectionFilter(field_name="doctype", values=[query.match_doctype])
        )
    if query.match_stage is not None:
        filters.append(CollectionFilter(field_name="stage", values=[query.match_stage]))
    return filters


def querydata_to_filters_kwargs(query: QueryData) -> dict:
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
