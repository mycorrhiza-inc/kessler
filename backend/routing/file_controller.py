from uuid import UUID
from typing import Annotated, assert_type
import logging

from litestar import Controller, Request

from litestar.handlers.http_handlers.decorators import (
    get,
    post,
    delete,
    patch,
    MediaType,
)


from sqlalchemy import select
from sqlalchemy.exc import IntegrityError, NoResultFound
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column




from litestar.params import Parameter
from litestar.di import Provide
from litestar.repository.filters import LimitOffset
from litestar.datastructures import UploadFile
from litestar.enums import RequestEncodingType
from litestar.params import Body
from litestar.logging import LoggingConfig

from pydantic import TypeAdapter, BaseModel


from models import FileModel, FileRepository, FileSchema, provide_files_repo


from crawler.docingest import DocumentIngester
from docprocessing.extractmarkdown import MarkdownExtractor
from docprocessing.genextras import GenerateExtras

from typing import List, Optional, Union, Any, Dict


from util.niclib import get_blake2

import json


class UUIDEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, UUID):
            # if the obj is uuid, we simply return the value of uuid
            return obj.hex
        return json.JSONEncoder.default(self, obj)


# for testing purposese
emptyFile = FileModel(
    url="",
    name="",
    doctype="",
    lang="en",
    path="",
    # file=raw_tmpfile,
    doc_metadata={},
    stage="stage0",
    hash="",
    summary=None,
    short_summary=None,
)



class FileUpdate(BaseModel):
    message: str
    metadata : Optional[Dict[str,Any]]


class UrlUpload(BaseModel):
    url: str
    metadata : Optional[Dict[str,Any]]


class UrlUploadList(BaseModel):
    url: List[str]


class FileCreate(BaseModel):
    message: str


class FileUpload(BaseModel):
    message: str


# litestar only

from pathlib import Path
import os

OS_TMPDIR = Path(os.environ["TMPDIR"])
OS_GPU_COMPUTE_URL = os.environ["GPU_COMPUTE_URL"]
OS_FILEDIR = Path("/files/")


from hashlib import blake2b

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
        return self.validate_and_jsonify(obj)

    @get(path="/test")
    async def test_api(self) -> None:
        return None

    @get(path="/files/all")
    async def get_all_files(
        self, files_repo: FileRepository, limit_offset: LimitOffset, request: Request
    ) -> list[FileSchema]:
        """List files."""
        results = await files_repo.list()
        type_adapter = TypeAdapter(list[FileSchema])
        return type_adapter.validate_python(results)

    @post(path="/files/upload", media_type=MediaType.TEXT)
    async def handle_file_upload(
        self,
        files_repo: FileRepository,
        data: Annotated[UploadFile, Body(media_type=RequestEncodingType.MULTI_PART)],
    ) -> Optional[FileUpload]:
        content = await data.read()
        newFileObj = emptyFile()
        newFileObj.name = data.filename

    @post(path="/files/add_url")
    async def add_url(
        self, files_repo: FileRepository, data: UrlUpload, request: Request, process : bool = False
    ) -> Any:
        request.logger.info("adding files")
        request.logger.info(data)
        # New stuff here, is this where this code belongs? <new stuff>
        docingest = DocumentIngester(request.logger)
        request.logger.info("DocumentIngester Created")
        tmpfile_path, metadata = docingest.url_to_filepath_and_metadata(data.url)
        override_metadata = data.metadata
        if override_metadata is not None:
            metadata = metadata.update(override_metadata)
        request.logger.info(f"Metadata Successfully Created with metadata {metadata}")
        document_title = metadata.get("title")
        document_doctype = metadata.get("doctype")
        document_lang = metadata.get("language")
        try:
            assert isinstance(document_title, str)
            assert isinstance(document_doctype, str)
            assert isinstance(document_lang, str)
        except:
            request.logger.error("Illformed Metadata please fix")
        else:
            request.logger.info(f"Title, Doctype and language successfully declared")
        request.logger.info("Attempting to save data to file")
        result = docingest.save_filepath_to_hash(tmpfile_path)
        (filehash, filepath) = result
        query = select(FileModel).where(FileModel.hash == filehash)
        duplicate_file_objects = await files_repo.session.execute(query)
        duplicate_file_obj = duplicate_file_objects.scalar()
        if duplicate_file_obj is None:
            docingest.backup_metadata_to_hash(metadata, filehash)
            new_file = FileModel(
                url=data.url,
                name=document_title,
                doctype=document_doctype,
                lang=document_lang,
                path=str(filepath),
                # file=raw_tmpfile,
                doc_metadata=metadata,
                stage="stage1",
                hash=filehash,
                summary=None,
                short_summary=None,
            )
            # </new stuff>
            request.logger.info("new file:{file}".format(file=new_file.to_dict()))
            try:
                new_file = await files_repo.add(new_file)
            except Exception as e:
                request.logger.info(e)
                return e
            request.logger.info("added file!~")
            await files_repo.session.commit()
            request.logger.info("commited file to DB")
        else:
            request.logger.info(type(duplicate_file_obj))
            request.logger.info(f"File with identical hash already exists in DB with uuid: {duplicate_file_obj.id}")
            new_file=duplicate_file_obj
        if process:
            request.logger.info("Processing File")
            await self.process_file_raw(new_file,files_repo,request.logger,False)
        return self.validate_and_jsonify(new_file)

    @post(path="/files/add_urls")
    async def add_urls(
        self, files_repo: FileRepository, data: UrlUploadList, request: Request
    ) -> None:
        return None

    @post(path="/process/{file_id_str:str}")
    async def process_file(
        self,
        files_repo: FileRepository,
        request: Request,
        file_id_str: str = Parameter(
            title="File ID as hex string", description="File to retieve"
        ),
        regenerate: bool = True,  # Figure out how to pass in a boolean as a query paramater
    ) -> FileSchema:
        """Process a File."""
        file_id = UUID(file_id_str)
        logger.info(file_id)
        obj = await files_repo.get(file_id)
        # TODO : Add error for invalid document ID
        await self.process_file_raw(obj,files_repo,request.logger,regenerate)
        return self.validate_and_jsonify(
            newobj
        )  # TODO : Return Response code and response message

    async def process_file_raw(self,obj : FileModel, files_repo : FileRepository , logger : Any,regenerate: bool):
        logger.info(type(obj))
        logger.info(obj)
        current_stage = obj.stage
        doctype = obj.doctype
        logger.info(obj.doctype)
        mdextract = MarkdownExtractor(logger, OS_GPU_COMPUTE_URL, OS_TMPDIR)
        genextras = GenerateExtras()

        response_code, response_message = (
            500,
            "Internal error somewhere in process.",
        )
        if current_stage == "stage0":
            response_code, response_message = (
                422,
                "Failure in stage 0: Document was incorrectly added to database, try readding it again.",
            )
            current_stage = "stage1"
        if regenerate and current_stage != "stage0":
            current_stage = "stage1"
        if current_stage == "stage1":
            try:
                processed_original_text = (
                    mdextract.process_raw_document_into_untranslated_text(
                        Path(obj.path),obj.doctype, obj.lang
                    )
                )
            except:
                response_code, response_message = (
                    422,
                    "failure in stage 1: document was unable to be converted to markdown,",
                )
            else:
                obj.original_text = processed_original_text
                current_stage = "stage2"
        if current_stage == "stage2":
            obj.english_text = obj.original_text
            current_stage = "stage3"
            # TODO : Fix issue of existence of languages aside from english
            # try:
            #     processed_english_text = mdextract.convert_text_into_eng(
            #         obj.original_text, obj.lang
            #     )
            # except:
            #     response_code, response_message = (
            #         422,
            #         "failure in stage 2: document was unable to be translated to english.",
            #     )
            # else:
        if current_stage == "stage3":
            links = genextras.extract_markdown_links(obj.original_text)
            long_sum = genextras.summarize_document_text(obj.original_text)
            short_sum = genextras.gen_short_sum_from_long_sum(long_sum)
            try:
                x=3
            except:
                response_code, response_message = (
                    422,
                    "failure in stage 3: Unable to generate summaries and links for document.",
                )
            else:
                obj.links = links
                obj.long_summary = long_sum
                obj.short_summary = short_sum
                current_stage = "stage4"
        if current_stage == "stage4":
            try:
                # TODO : Chunk document and generate embeddings.
                print("Create Embeddings.")
            except:
                response_code, response_message = (
                    422,
                    "failure in stage 2: document was unable to be translated to english.",
                )
            else:
                current_stage = "stage5"

        if current_stage == "completed":
            response_code, response_message = (200, "Document Fully Processed.")
        logger.info(current_stage)
        obj.stage=current_stage
        logger.info(response_code)
        logger.info(response_message)
        newobj = files_repo.update(obj)

        await files_repo.session.commit()
    # @patch(path="/files/{file_id:uuid}")
    # async def update_file(
    #     self,
    #     files_repo: FileRepository,
    #     data: FileUpdate,
    #     file_id: UUID = Parameter(
    #         title="File ID", description="File to retieve"),
    # ) -> FileSchema:
    #     """Update a File."""
    #     raw_obj = data.model_dump(exclude_unset=True, exclude_none=True)
    #     raw_obj.update({"id": file_id})
    #     obj = files_repo.update(FileModel(**raw_obj))
    #     await files_repo.session.commit()
    #     return self.validate_and_jsonify(obj)

    @delete(path="/files/{file_id:uuid}")
    async def delete_file(
        self,
        files_repo: FileRepository,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
    ) -> None:
        _ = await files_repo.delete(file_id)
        await files_repo.session.commit()

    def validate_and_jsonify(self,file_object : FileModel) -> dict:
        validated =  FileSchema.model_validate(file_object)
        validated_dict = dict(validated)
        validated_dict["id"]= str(validated_dict["id"])
        return validated_dict
