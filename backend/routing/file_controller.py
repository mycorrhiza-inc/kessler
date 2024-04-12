from uuid import UUID
from typing import Annotated


from litestar import Controller, Request

from litestar.handlers.http_handlers.decorators import (
    get,
    post,
    delete,
    patch,
    MediaType,
)

from litestar.params import Parameter
from litestar.di import Provide
from litestar.repository.filters import LimitOffset
from litestar.datastructures import UploadFile
from litestar.enums import RequestEncodingType
from litestar.params import Body

from pydantic import TypeAdapter, BaseModel


from models import FileModel, FileRepository, FileSchema, provide_files_repo


from crawler.docingest import DocumentIngester
from docprocessing.extractmarkdown import MarkdownExtractor
from docprocessing.genextras import GenerateExtras


# for testing purposese
emptyFile = FileModel(
    path="",
    doctype="",
    lang="",
    name="",
    stage="unprocessed",
    summary=None,
    short_summary=None,
)


class FileUpdate(BaseModel):
    message: str


class FileCreate(BaseModel):
    message: str


class FileUpload(BaseModel):
    message: str


# litestar only


class FileController(Controller):
    """File Controller"""

    dependencies = {"files_repo": Provide(provide_files_repo)}

    @get(path="/files/{file_id:uuid}")
    async def get_file(
        self,
        files_repo: FileRepository,
        file_id: UUID = Parameter(
            title="File ID", description="File to retieve"),
    ) -> FileSchema:
        obj = files_repo.get(file_id)
        return FileSchema.model_validate(obj)

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
    ) -> FileUpload:
        content = await data.read()
        newFileObj = emptyFile()
        newFileObj.name = data.filename

    @post(path="/files/add")
    async def add_file(
        self, files_repo: FileRepository, data: FileCreate, request: Request
    ) -> FileSchema:
        request.logger.info("adding files")
        request.logger.info(data)
        # New stuff here, is this where this code belongs? <new stuff>
        docingest = DocumentIngester()
        metadata, raw_file_path = docingest.url_to_file_and_metadata(data.url)
        new_file = FileModel(
            url=data.url,
            title=metadata["title"],
            doctype=metadata["doctype"],
            lang=metadata["lang"],
            file=open(raw_file_path),
            metadata=metadata,
            stage="stage0",
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
        return FileSchema.model_validate(new_file)

    @patch(path="/files/{file_id:uuid}")
    async def process_File(
        self,
        files_repo: FileRepository,
        data: FileUpdate,
        file_id: UUID = Parameter(
            title="File ID", description="File to retieve"),
        regenerate: bool = False,  # Figure out how to pass in a boolean as a query paramater
    ) -> FileSchema:
        """Process a File."""
        obj = files_repo.get(file_id)
        current_stage = obj.stage
        mdextract = MarkdownExtractor()
        genextras = GenerateExtras()

        if current_stage == "stage0":
            response_code, response_message = (
                422,
                "Failure in stage 0: Document was incorrectly added to database, try readding it again.",
            )
        if regenerate and current_stage != "stage0":
            current_stage = "stage1"
        if current_stage == "stage1":
            try:
                processed_original_text = (
                    mdextract.process_raw_document_into_untranslated_text(
                        obj.path, obj.metadata
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
            try:
                processed_english_text = mdextract.convert_text_into_eng(
                    obj.original_text, obj.lang
                )
            except:
                response_code, response_message = (
                    422,
                    "failure in stage 2: document was unable to be translated to english.",
                )
            else:
                obj.english_text = processed_english_text
                current_stage = "stage3"
        if current_stage == "stage3":
            try:
                links = genextras.extract_markdown_links(
                    obj.original_text, obj.lang)
                long_sum = genextras.summarize_document_text(obj.original_text)
                short_sum = genextras.gen_short_sum_from_long_sum(long_sum)
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

                current_stage = "completed"
        if current_stage == "completed":
            response_code, r3esponse_message = (
                200, "Document Fully Processed.")
        newobj = files_repo.update(obj)
        await files_repo.session.commit()
        return FileSchema.model_validate(
            newobj
        )  # TODO : Return Response code and response message

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
    #     return FileSchema.model_validate(obj)

    @delete(path="/files/{file_id:uuid}")
    async def delete_file(
        self,
        files_repo: FileRepository,
        file_id: UUID = Parameter(
            title="File ID", description="File to retieve"),
    ) -> None:
        _ = await files_repo.delete(file_id)
        await files_repo.session.commit()
