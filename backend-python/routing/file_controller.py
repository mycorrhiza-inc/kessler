from util.file_io import S3FileManager
from vecstore.docprocess import add_document_to_db
import os
from pathlib import Path
from uuid import UUID

from litestar import Controller, Request, Response

from litestar.handlers.http_handlers.decorators import (
    get,
    post,
    delete,
    MediaType,
)


from sqlalchemy import select

from sqlalchemy.ext.asyncio import AsyncSession

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
    provide_files_repo,
    model_to_schema,
    get_texts_from_file_uuid,
)

from models.utils import provide_async_session

from common.file_schemas import (
    FileSchema,
    DocumentStatus,
    FileTextSchema,
    docstatus_index,
)


from typing import List, Optional, Dict, Annotated, Tuple, Any


import json

from common.niclib import rand_string, paginate_results

from enum import Enum

from sqlalchemy import and_

from logic.databaselogic import QueryData, filter_list_mdata, querydata_to_filters

from constants import (
    OS_TMPDIR,
)


class UUIDEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, UUID):
            # if the obj is uuid, we simply return the value of uuid
            return obj.hex
        return json.JSONEncoder.default(self, obj)


# TODO : Create test that adds a file once we know what the file DB schema is going to look like


class FileUpdate(BaseModel):
    message: str
    metadata: Dict[str, Any]


class UrlUpload(BaseModel):
    url: str
    metadata: Dict[str, Any] = {}


class UrlUploadList(BaseModel):
    url: List[str]


class FileCreate(BaseModel):
    message: str


class FileUpload(BaseModel):
    message: str


class IndexFileRequest(BaseModel):
    id: UUID


class FileTextUpload(BaseModel):
    text: str
    doctype: str
    metadata: Dict[str, Any]


def ensure_metadata(metadata: Dict[str, Any]) -> None:
    assert metadata.get("title") is not None


class FileController(Controller):
    """File Controller"""

    dependencies = {
        "files_repo": Provide(provide_files_repo),
        "db_session": Provide(provide_async_session),
    }

    # def jsonify_validate_return(self,):
    #     return None

    @get(path="/files/{file_id:uuid}")
    async def get_file(
        self,
        files_repo: FileRepository,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
    ) -> FileSchema:
        obj = await files_repo.get(file_id)

        return model_to_schema(obj)

    @get(path="/files/markdown/{file_id:uuid}")
    async def get_markdown(
        self,
        db_connection: AsyncSession,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
        original_lang: bool = False,
        match_lang: Optional[str] = None,
    ) -> str:
        # TODO: Return 404 if doc not found.
        texts = await get_texts_from_file_uuid(db_connection, file_id)

        def filter_original_lang(doc: FileTextSchema) -> bool:
            return doc.is_original_lang

        if original_lang:
            texts = list(filter(filter_original_lang, texts))
        # Every doc should have original lang, so if it has a translated lang but not an original, something is very wrong
        if len(texts) == 0:
            raise Exception(f"No Texts Found for document {file_id}")
        search_first_lang = match_lang
        if search_first_lang is None:
            search_first_lang = "en"

        def filter_lang(doc: FileTextSchema) -> bool:
            return doc.language == search_first_lang

        match_lang_texts = list(filter(filter_lang, texts))
        if len(texts) > 0:
            markdown_text = match_lang_texts[0].text
        if match_lang is not None:
            raise Exception(
                "Could not find Text matching language, also this should return a 404 error you can bother the devs to fix it"
            )
        return texts[0].text

    @get(path="/files/raw/{file_id:uuid}")
    async def get_raw(
        self,
        files_repo: FileRepository,
        request: Request,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
    ) -> Response:
        logger = request.logger
        obj = await files_repo.get(file_id)
        if obj is None:
            return Response(content="ID does not exist", status_code=404)

        filehash = obj.hash

        file_path = S3FileManager(logger=logger).generate_local_filepath_from_hash(
            filehash, ensure_network=False
        )
        if file_path is None or not file_path.is_file():
            return Response(content="File not found", status_code=404)

        # Read the file content
        with open(file_path, "rb") as file:
            file_content = file.read()
        # currently doesnt work unfortunately
        # file_name = obj.name
        # headers = {
        #     "Content-Disposition": f'attachment; filename="{file_name}"'
        # }

        return Response(content=file_content, media_type="application/octet-stream")

        # Return as a result of the get request, the file at file_path. Also make sure to include the correct return type.

    @get(path="/files/metadata/{file_id:uuid}")
    async def get_metadata(
        self,
        files_repo: FileRepository,
        request: Request,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
    ) -> dict:
        logger = request.logger
        obj = await files_repo.get(file_id)
        if obj is None:
            return Response(content="ID does not exist", status_code=404)

        obj = model_to_schema(obj)
        return obj.mdata

    async def get_all_files_raw(
        self, files_repo: FileRepository, logger: Any
    ) -> list[FileSchema]:
        results = await files_repo.list()
        logger.info(f"{len(results)} results")
        valid_results = list(map(model_to_schema, results))
        return valid_results

    @post(path="/files/all")
    async def get_all_files(
        self,
        files_repo: FileRepository,
        request: Request,
        ensure_all_on_s3: bool = False,
    ) -> list[FileSchema]:
        """List files."""
        valid_results = await self.get_all_files_raw(files_repo, request.logger)
        s3 = S3FileManager()
        if ensure_all_on_s3:
            for result in valid_results:
                hash = result.hash
                if hash is not None:
                    s3.generate_local_filepath_from_hash(
                        hash, ensure_network=True, download_local=False
                    )
        return valid_results

    @post(path="/files/all/paginate")
    async def get_all_files_paginated(
        self,
        files_repo: FileRepository,
        request: Request,
        num_results: Optional[int],
        page: Optional[int],
    ) -> Tuple[list[FileSchema], int]:
        """List files."""
        valid_results = await self.get_all_files_raw(files_repo, request.logger)
        return paginate_results(valid_results, num_results, page)

    async def query_all_files_raw(
        self, files_repo: FileRepository, query: QueryData, logger: Any
    ) -> List[FileSchema]:
        filters = querydata_to_filters(query)
        results = await files_repo.list(*filters)
        logger.info(f"{len(results)} results")
        valid_results = list(map(model_to_schema, results))
        if query.match_metadata is None or query.match_metadata == {}:
            return valid_results
        filtered_valid_results = filter_list_mdata(valid_results, query.match_metadata)
        return filtered_valid_results

    @post(path="/files/query")
    async def query_all_files(
        self,
        files_repo: FileRepository,
        data: QueryData,
        request: Request,
    ) -> list[FileSchema]:
        """List files."""
        return await self.query_all_files_raw(files_repo, data, request.logger)

    @post(path="/files/query/paginate")
    async def query_all_files_paginated(
        self,
        files_repo: FileRepository,
        data: QueryData,
        request: Request,
        num_results: Optional[int],
        page: Optional[int],
    ) -> Tuple[list[FileSchema], int]:
        valid_results = await self.query_all_files_raw(files_repo, data, request.logger)
        return paginate_results(valid_results, num_results, page)

    # TODO: replace this with a jobs endpoint

    @post(path="/files/upload_file", media_type=MediaType.TEXT)
    async def handle_file_upload(
        self,
        files_repo: FileRepository,
        data: Annotated[UploadFile, Body(media_type=RequestEncodingType.MULTI_PART)],
        request: Request,
        process: bool = True,
        override_hash: bool = False,
    ) -> str:
        raise Exception(
            "This should probably be forwarded to thaumaturgy, since ingest is covered with that asynchronously. Priority Queue should make it pretty fast."
        )
        logger = request.logger
        logger.info("Process initiated.")
        supplemental_metadata = {"source": "personal"}

        docingest = DocumentIngester(logger)
        input_directory = OS_TMPDIR / Path("formdata_uploads") / Path(rand_string())
        # Ensure the directories exist
        os.makedirs(input_directory, exist_ok=True)
        # Save the PDF to the output directory
        filename = data.filename
        final_filepath = input_directory / Path(filename)
        with open(final_filepath, "wb") as f:
            f.write(data.file.read())
        additional_metadata = docingest.infer_metadata_from_path(final_filepath)
        additional_metadata.update(supplemental_metadata)
        final_metadata = additional_metadata
        if final_metadata.get("lang") is None:
            final_metadata["lang"] = "en"

        file_obj = await add_file_raw(
            final_filepath, final_metadata, process, override_hash, files_repo, logger
        )
        return f"Successfully added document with uuid: {file_obj.id}"

    # TODO: Implement security for this method so people cant just delete everyting from db + remove from vector DB
    # @delete(path="/files/{file_id:uuid}")
    # async def delete_file(
    #     self,
    #     files_repo: FileRepository,
    #     file_id: UUID = Parameter(
    #         title="File ID as hex string", description="File to delete"
    #     ),
    # ) -> None:
    #     await files_repo.delete(file_id)
    #     await files_repo.session.commit()
