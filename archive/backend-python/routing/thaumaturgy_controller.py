from util.file_io import S3FileManager
from vecstore.docprocess import add_document_to_db
from uuid import UUID

from litestar import Controller, Request, Response

from litestar.handlers.http_handlers.decorators import (
    get,
    post,
    delete,
    MediaType,
)


from sqlalchemy import select


from litestar.params import Parameter
from litestar.di import Provide
from litestar.repository.filters import LimitOffset
from litestar.datastructures import UploadFile
from litestar.enums import RequestEncodingType
from litestar.params import Body

from pydantic import TypeAdapter
from models.utils import PydanticBaseModel as BaseModel


from models.files import (
    FileModel,
    FileRepository,
    get_full_file_from_uuid,
    provide_files_repo,
    file_model_to_schema,
)
from common.file_schemas import FileSchema, DocumentStatus, docstatus_index


from typing import List, Dict, Any


import json

from common.niclib import rand_string, paginate_results


from sqlalchemy import and_

from logic.databaselogic import QueryData, filter_list_mdata, querydata_to_filters

from constants import (
    OS_TMPDIR,
)

from common.file_schemas import DocumentStatus, FileSchemaFull
from sqlalchemy.ext.asyncio import AsyncSession

from models.files import upsert_file_from_full_schema
import logging


default_logger = logging.getLogger(__name__)


class UUIDEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, UUID):
            # if the obj is uuid, we simply return the value of uuid
            return obj.hex
        return json.JSONEncoder.default(self, obj)


# TODO : Create test that adds a file once we know what the file DB schema is going to look like


class FileEmbeddings(BaseModel):
    file_id: UUID
    text: str
    metadata: Dict[str, Any]
    strings: List[str] | None
    embeddings: List[List[float]] | None


# import base64


# import base64


class ThaumaturgyController(Controller):
    """File Controller"""

    # def jsonify_validate_return(self,):
    #     return None
    # TODO: ADD some kind of authentication to this entire controller
    @get(path="/thaumaturgy/full_file/{file_id:uuid}")
    async def get_file_by_id(
        self,
        db_session: AsyncSession,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
    ) -> FileSchemaFull:
        result = await get_full_file_from_uuid(db_session, file_id)
        return result

    @post(path="/thaumaturgy/upsert_file", media_type=MediaType.TEXT)
    async def upsert_file_dangerous(
        self,
        db_session: AsyncSession,
        data: FileSchemaFull,
        process: bool = True,
        override_hash: bool = False,
    ) -> str:
        file = data
        await upsert_file_from_full_schema(db_session, file)
        return f"Successfully added document with uuid: {file.id}"

    @post(path="/thaumaturgy/insert_file_embeddings", media_type=MediaType.TEXT)
    async def create_file_embeddings(
        self,
        db_session: AsyncSession,
        data: FileEmbeddings,
        request: Request,
        remove_old: bool = False,
    ) -> None:
        await self.create_file_embeddings_raw(data=data, remove_old=remove_old)

    async def create_file_embeddings_raw(
        self,
        data: FileEmbeddings,
        remove_old: bool,
    ) -> None:
        if remove_old:
            pass
        logger = default_logger
        source_id = data.file_id
        doc_metadata = data.metadata
        embedding_text = data.text

        logger.info("Adding Document to Vector Database")

        def generate_searchable_metadata(initial_metadata: dict) -> dict:
            return_metadata = {
                "title": initial_metadata.get("title"),
                "author": initial_metadata.get("author"),
                "source": initial_metadata.get("source"),
                "date": initial_metadata.get("date"),
                "source_id": source_id,
            }

            def guarentee_field(field: str, default_value: Any = "unknown"):
                if return_metadata.get(field) is None:
                    return_metadata[field] = default_value

            guarentee_field("title")
            guarentee_field("author")
            guarentee_field("source")
            guarentee_field("date")
            return return_metadata

        searchable_metadata = generate_searchable_metadata(doc_metadata)
        try:
            add_document_to_db(embedding_text, metadata=searchable_metadata)
        except Exception as e:
            raise Exception("Failure in adding document to vector database", e)
