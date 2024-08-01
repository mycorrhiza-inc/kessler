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


from litestar.params import Parameter
from litestar.di import Provide
from litestar.repository.filters import LimitOffset
from litestar.datastructures import UploadFile
from litestar.enums import RequestEncodingType
from litestar.params import Body

from pydantic import TypeAdapter
from models.utils import PydanticBaseModel as BaseModel


# from models import (
#     FileModel,
#     FileRepository,
#     FileSchema,
#     provide_files_repo,
# )
from models.files import (
    FileModel,
    FileRepository,
    FileSchema,
    FileSchemaWithText,
    provide_files_repo,
    DocumentStatus,
    docstatus_index,
    model_to_schema,
)

from logic.docingest import DocumentIngester
from logic.extractmarkdown import MarkdownExtractor

from typing import List, Optional, Dict, Annotated, Tuple, Any


from logic.filelogic import add_file_raw, process_file_raw


import json

from util.niclib import rand_string, paginate_results

from enum import Enum

from sqlalchemy import and_

from logic.databaselogic import QueryData, filter_list_mdata, querydata_to_filters


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


from constants import (
    OS_TMPDIR,
    OS_GPU_COMPUTE_URL,
    OS_FILEDIR,
    OS_HASH_FILEDIR,
    OS_OVERRIDE_FILEDIR,
    OS_BACKUP_FILEDIR,
)


# import base64


# import base64


class FileController(Controller):
    """File Controller"""

    dependencies = {"files_repo": Provide(provide_files_repo)}

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
        files_repo: FileRepository,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
        original_lang: bool = False,
    ) -> str:
        # Yes I know this is a redundant if, this looks much more readable imo.
        if original_lang is True:
            return "Feature delayed due to only supporting english documents."
        obj = await files_repo.get(file_id)

        markdown_text = obj.english_text

        if markdown_text == "":
            markdown_text = "Could not find Document Markdown Text"
        return markdown_text

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

    # TODO : (Nic) Make function that can process uploaded files
    @post(path="/files/add_url")
    async def add_url(
        self,
        files_repo: FileRepository,
        data: UrlUpload,
        request: Request,
        process: bool = True,
        override_hash: bool = False,
    ) -> str:
        logger = request.logger
        logger.info("adding files")
        logger.info(data)

        # ------------------ here be a site for refactor --------------------
        docingest = DocumentIngester(logger)

        logger.info("DocumentIngester Created")

        # tmpfile_path, metadata = (
        # LSP is giving some kind of error, I am gonna worry about it later
        tmpfile_path, metadata = docingest.url_to_filepath_and_metadata(data.url)
        new_metadata = data.metadata

        if new_metadata is not None:
            metadata.update(new_metadata)

        request.logger.info(f"Metadata Successfully Created with metadata {metadata}")

        file_obj = await add_file_raw(
            tmpfile_path, metadata, process, override_hash, files_repo, logger
        )
        # type_adapter = TypeAdapter(FileSchema)
        # final_return = model_to_schema(new_file)
        # logger.info(final_return)
        return f"Successfully added document with uuid: {file_obj.uuid}"

    # TODO: anything but this

    @post(path="/process/{file_id_str:str}")
    async def process_file(
        self,
        files_repo: FileRepository,
        request: Request,
        file_id_str: str = Parameter(
            title="File ID as hex string", description="File to retieve"
        ),
        regenerate_from: Optional[
            str
        ] = None,  # Figure out how to pass in a boolean as a query paramater
        stop_at: Optional[
            str
        ] = None,  # Figure out how to pass in a boolean as a query paramater
    ) -> None:
        """Process a File."""
        logger = request.logger
        if stop_at is None:
            stop_at = "completed"
        if regenerate_from is None:
            regenerate_from = "completed"
        # TODO: Figure out error messaging for these
        stop_at = DocumentStatus(stop_at)
        regenerate_from = DocumentStatus(regenerate_from)
        # TODO : Add error for invalid document ID
        file_id = UUID(file_id_str)
        logger.info(file_id)
        obj = await files_repo.get(file_id)
        if docstatus_index(DocumentStatus(obj.stage)) < docstatus_index(
            regenerate_from
        ):
            obj.stage = regenerate_from.value

        if docstatus_index(DocumentStatus(obj.stage)) < docstatus_index(stop_at):
            await process_file_raw(
                obj, files_repo, request.logger, stop_at, priority=True
            )
        # TODO : Return Response code and response message
        return self.validate_and_jsonify(obj)

    @post(path="/files/upload/from/md", media_type=MediaType.TEXT)
    async def upload_from_markdown(
        self,
        files_repo: FileRepository,
        request: Request,
        data: Annotated[UploadFile, Body(media_type=RequestEncodingType.MULTI_PART)],
    ) -> None:
        try:
            content = await data.read()
            filename = data.filename
            file = content.decode()
            splitfile = file.split("---")
            restfile = "".join(splitfile[2:])
            file_metadata = splitfile[1].split("\n")
            meta = {}
            for i in file_metadata:
                if i == "":
                    continue
                field = i.split(":")
                if len(field) >= 2:
                    meta[field[0]] = "".join(field[1:])

            m_text = json.dumps(meta)

            FileModel(english_text=file, metadata=m_text)
            new_file = FileModel(
                url="",
                name=filename,
                doctype="mardown",
                lang="english",
                source="markdown",
                metadata=m_text,
                stage=DocumentStatus.completed.value,
                hash="None",
                summary=None,
                short_summary=None,
                english_text=restfile,
            )
            try:
                files_repo.session.add(new_file)
                await files_repo.session.flush()
                files_repo.session.refresh(new_file)
                await files_repo.session.commit()
            except Exception as e:
                return f"issue: \n{e}"
            try:
                meta["source_id"] = str(new_file.id)
                add_document_to_db(text=restfile, metadata=meta)
            except Exception as e:
                request.logger.error(e)
                return "issue indexing file"
            return new_file.english_text

        except Exception as e:
            raise (e)

    @delete(path="/files/{file_id:uuid}")
    async def delete_file(
        self,
        files_repo: FileRepository,
        file_id: UUID = Parameter(
            title="File ID as hex string", description="File to delete"
        ),
    ) -> None:
        await files_repo.delete(file_id)
        await files_repo.session.commit()
