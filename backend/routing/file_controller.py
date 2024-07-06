from lance_store.connection import ensure_fts_index
from rag.llamaindex import add_document_to_db_from_text
import os
from pathlib import Path
from typing import Any
from uuid import UUID
from typing import Annotated

from litestar import Controller, Request

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
)

from crawler.docingest import DocumentIngester
from docprocessing.extractmarkdown import MarkdownExtractor

from typing import List, Optional, Dict


import json

from util.niclib import rand_string


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
    metadata: Dict[str, Any]


class UrlUploadList(BaseModel):
    url: List[str]


class FileCreate(BaseModel):
    message: str


class FileUpload(BaseModel):
    message: str


class IndexFileRequest(BaseModel):
    id: UUID


# litestar only


OS_TMPDIR = Path(os.environ["TMPDIR"])
OS_GPU_COMPUTE_URL = os.environ["GPU_COMPUTE_URL"]
OS_FILEDIR = Path("/files/")


# import base64


OS_TMPDIR = Path(os.environ["TMPDIR"])
OS_GPU_COMPUTE_URL = os.environ["GPU_COMPUTE_URL"]
OS_FILEDIR = Path("/files/")
OS_HASH_FILEDIR = OS_FILEDIR / Path("raw")
OS_OVERRIDE_FILEDIR = OS_FILEDIR / Path("override")
OS_BACKUP_FILEDIR = OS_FILEDIR / Path("backup")


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

        type_adapter = TypeAdapter(FileSchema)

        return type_adapter.validate_python(obj)

    @get(path="/files/get_markdown/{file_id:uuid}")
    async def get_markdown(
        self,
        files_repo: FileRepository,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
    ) -> str:
        obj = await files_repo.get(file_id)

        type_adapter = TypeAdapter(FileSchemaWithText)

        obj_with_text=type_adapter.validate_python(obj)

        markdown_text = obj_with_text.english_text 
        if markdown_text is "":
            markdown_text="Could not find Document Markdown Text"
        return markdown_text

    @get(path="/files/all")
    async def get_all_files(
        self, files_repo: FileRepository, limit_offset: LimitOffset, request: Request
    ) -> list[FileSchema]:
        """List files."""
        results = await files_repo.list()
        logger = request.logger
        logger.info(f"{len(results)} results")
        type_adapter = TypeAdapter(list[FileSchema])
        return type_adapter.validate_python(results)

    # TODO: replace this with a jobs endpoint

    @post(path="/files/upload_file", media_type=MediaType.TEXT)
    async def handle_file_upload(
        self,
        files_repo: FileRepository,
        data: Annotated[UploadFile, Body(media_type=RequestEncodingType.MULTI_PART)],
        request: Request,
        process: bool = True,
        override_hash: bool = False,
    ) -> Any:
        supplemental_metadata = {"source": "personal"}
        logger = request.logger
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
        return await self.add_file_raw(
            final_filepath, final_metadata, process, override_hash, files_repo, logger
        )

    # TODO : (Nic) Make function that can process uploaded files
    @post(path="/files/add_url")
    async def add_url(
        self,
        files_repo: FileRepository,
        data: UrlUpload,
        request: Request,
        process: bool = True,
        override_hash: bool = False,
    ) -> Any:
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
        return await self.add_file_raw(
            tmpfile_path, metadata, process, override_hash, files_repo, logger
        )

    async def add_file_raw(
        self,
        tmp_filepath: Path,
        metadata: dict,
        process: bool,
        override_hash: bool,
        files_repo: FileRepository,
        logger: Any,
    ) -> Any:
        docingest = DocumentIngester(logger)
        document_title = metadata.get("title")
        document_doctype = metadata.get("doctype")
        document_lang = metadata.get("language")

        try:
            assert isinstance(document_title, str)
            assert isinstance(document_doctype, str)
            assert isinstance(document_lang, str)
        except Exception:
            logger.error("Illformed Metadata please fix")

        else:
            logger.info("Title, Doctype and language successfully declared")

        document_source = metadata.get("source")
        if document_source is None:
            document_source = "UNKNOWN"
            metadata["source"] = "UNKNOWN"

        logger.info("Attempting to save data to file")
        result = docingest.save_filepath_to_hash(tmp_filepath, OS_HASH_FILEDIR)
        (filehash, filepath) = result

        os.remove(tmp_filepath)

        # NOTE: this is a dangeous query
        # NOTE: Nicole- Also this doesnt allow for files with the same hash to have different metadata,
        # Scrapping it is a good idea, it was designed to solve some issues I had during testing and adding the same dataset over and over
        # FIX: fix this to not allow for users to DOS files
        query = select(FileModel).where(FileModel.hash == filehash)

        duplicate_file_objects = await files_repo.session.execute(query)
        duplicate_file_obj = duplicate_file_objects.scalar()

        if override_hash is True and duplicate_file_obj is not None:
            try:
                await files_repo.delete(duplicate_file_obj.id)
            except Exception:
                pass
            duplicate_file_obj = None

        if duplicate_file_obj is None:
            docingest.backup_metadata_to_hash(metadata, filehash)
            metadata_str = json.dumps(metadata)
            new_file = FileModel(
                url="N/A",
                name=document_title,
                doctype=document_doctype,
                lang=document_lang,
                source=document_source,
                mdata=metadata_str,
                stage="stage1",
                hash=filehash,
                summary=None,
                short_summary=None,
            )
            logger.info("new file:{file}".format(file=new_file.to_dict()))
            try:
                new_file = await files_repo.add(new_file)
            except Exception as e:
                logger.info(e)
                return e
            logger.info("added file!~")
            await files_repo.session.commit()
            logger.info("commited file to DB")

        else:
            logger.info(type(duplicate_file_obj))
            logger.info(
                f"File with identical hash already exists in DB with uuid:\
                {duplicate_file_obj.id}"
            )
            new_file = duplicate_file_obj

        if process:
            logger.info("Processing File")
            await self.process_file_raw(new_file, files_repo, logger, False)

        type_adapter = TypeAdapter(FileSchema)
        final_return = type_adapter.validate_python(new_file)

        logger.info(final_return)
        # return final_return
        return "This should return the file, but serialization is broken, I hope this is string is a fair substitute."

    @post(path="/files/add_urls")
    async def add_urls(
        self, files_repo: FileRepository, data: UrlUploadList, request: Request
    ) -> None:
        return None

    # TODO: anything but this

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
        request.logger.info(file_id)
        obj = await files_repo.get(file_id)
        # TODO : Add error for invalid document ID
        await self.process_file_raw(obj, files_repo, request.logger, regenerate)
        # TODO : Return Response code and response message
        return self.validate_and_jsonify(obj)

    async def process_file_raw(
        self, obj: FileModel, files_repo: FileRepository, logger: Any, regenerate: bool
    ):
        logger.info(type(obj))
        logger.info(obj)
        current_stage = obj.stage
        logger.info(obj.doctype)
        mdextract = MarkdownExtractor(logger, OS_TMPDIR)
        doc_metadata = json.loads(obj.mdata)

        response_code, response_message = (
            500,
            "Internal error somewhere in process.",
        )

        if regenerate:
            current_stage = "stage1"

        # text extraction
        def process_stage_one():
            # FIXME: Change to deriving the filepath from the uri.
            file_path = DocumentIngester(logger).get_default_filepath_from_hash(
                obj.hash
            )
            # This process might spit out new metadata that was embedded in the document, ignoring for now
            processed_original_text = (
                mdextract.process_raw_document_into_untranslated_text(
                    file_path, doc_metadata
                )[0]
            )
            logger.info(
                f"Successfully processed original text: {processed_original_text[0:20]}"
            )
            # FIXME: We should probably come up with a better backup protocol then doing everything with hashes
            mdextract.backup_processed_text(
                processed_original_text, obj.hash, doc_metadata, OS_BACKUP_FILEDIR
            )
            logger.info("Backed up markdown text")
            if obj.lang == "en":
                # Write directly to the english text box if
                # original text is identical to save space.
                obj.english_text = processed_original_text
                # Skip translation stage if text already english.
                return "stage3"
            else:
                obj.original_text = processed_original_text
                return "stage2"

        # text conversion
        def process_stage_two():
            if obj.lang != "en":
                try:
                    processed_english_text = mdextract.convert_text_into_eng(
                        obj.original_text, obj.lang
                    )
                    obj.english_text = processed_english_text
                except Exception as e:
                    raise Exception(
                        "failure in stage 2: \ndocument was unable to be translated to english.",
                        e,
                    )
            else:
                raise ValueError(
                    "failure in stage 2: \n Code is in an unreachable state, a document cannot be english and not english",
                )
            return "stage3"

        # text commitment
        def process_stage_three():
            logger.info("Adding Document to Vector Database")
            try:
                add_document_to_db_from_text(obj.english_text, doc_metadata)
            except Exception as e:
                raise Exception("Failure in adding document to vector database", e)
            return "completed"

        while True:
            match current_stage:
                case "stage1":
                    current_stage = process_stage_one()
                case "stage2":
                    current_stage = process_stage_two()
                case "stage3":
                    current_stage = process_stage_three()
                case "completed":
                    response_code, response_message = (
                        200,
                        "Document Fully Processed.",
                    )
                    logger.info(current_stage)
                    obj.stage = current_stage
                    logger.info(response_code)
                    logger.info(response_message)
                    _ = files_repo.update(obj)
                    await files_repo.session.commit()
                    break
                case _:
                    raise Exception(
                        "Document was incorrectly added to database, \
                        try readding it again.\
                    "
                    )
                # FIXME: The try catch exception code broke the plaintext error handling, since it still returns a 500 error, I removed it temporarially
                # try:
                #    match current_stage:
                #        case "stage1":
                #            current_stage = process_stage_one()
                #        case "stage2":
                #            current_stage = process_stage_two()
                #        case "stage3":
                #            current_stage = process_stage_three()
                #        case "completed":
                #            response_code, response_message = (
                #                200,
                #                "Document Fully Processed.",
                #            )
                #            logger.info(current_stage)
                #            obj.stage = current_stage
                #            logger.info(response_code)
                #            logger.info(response_message)
                #            _ = files_repo.update(obj)
                #            await files_repo.session.commit()
                #            break
                #        case _:
                #            raise Exception(
                #                "Document was incorrectly added to database, \
                #                try readding it again.\
                #            "
                #            )

                # except Exception as e:
                #    logger.error(e)
                #    break

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
                stage="completed",
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
                meta["uid"] = str(new_file.id)
                add_document_to_db_from_text(text=restfile, metadata=meta)
                request.app.emit("increment_processed_docs", num=1)
                request.logger.info("added a document to the db")
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
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
    ) -> None:
        fid = UUID(file_id)
        _ = await files_repo.delete(fid)
        await files_repo.session.commit()
