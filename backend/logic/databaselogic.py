from pydantic import BaseModel
from typing import Optional, Any, List
from advanced_alchemy.filters import SearchFilter, CollectionFilter
from models.files import DocumentStatus, docstatus_index, FileSchema


class QueryData(BaseModel):
    match_name: Optional[str] = None
    match_source: Optional[str] = None
    match_doctype: Optional[str] = None
    match_stage: Optional[str] = None
    match_metadata: Optional[dict] = None


def querydata_to_filters(query: QueryData) -> list:
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


def querydata_to_filters_strict(query: QueryData) -> list:
    filters = []
    if query.match_name is not None:
        filters.append(CollectionFilter(field_name="name", values=[query.match_name]))
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


def filter_list_mdata(
    schema_list: List[FileSchema], mdata_match: dict
) -> List[FileSchema]:

    def is_valid(item: FileSchema) -> bool:
        if item.mdata is None:
            return False
        for key in mdata_match:
            if item.mdata[key] != mdata_match["key"] or mdata_match["key"] is None:
                return False
        return True

    return list(filter(is_valid, schema_list))


# def querydata_to_filters_kwargs(query: QueryData) -> dict:
#     filters = {}
#     if query.match_name is not None:
#         filters["name"] = query.match_name
#     if query.match_source is not None:
#         filters["source"] = query.match_source
#     if query.match_doctype is not None:
#         filters["doctype"] = query.match_doctype
#     if query.match_stage is not None:
#         filters["stage"] = query.match_stage
#     return filters


def filters_docstatus_processing(
    stop_at: DocumentStatus, regenerate_from: DocumentStatus
) -> list:
    stop_index = docstatus_index(stop_at)
    regen_index = docstatus_index(regenerate_from)
    valid_values = []
    for status in DocumentStatus:
        istatus = docstatus_index(status)
        if (istatus < stop_index) or (
            stop_index > regen_index and istatus > regen_index
        ):
            valid_values.append(status.value)

    return [CollectionFilter(field_name="stage", values=valid_values)]
