from uuid import UUID


from litestar import Controller, Request

from litestar.handlers.http_handlers.decorators import get, post, delete, patch
from litestar.params import Parameter
from litestar.di import Provide
from litestar.repository.filters import LimitOffset

from pydantic import TypeAdapter, validator

from db import BaseModel

from models.files import File

from modules.files.dbm.files import provide_files_repo


class FileUpload(BaseModel):
    file_metadata: dict
    # Figure out how to do a file upload datatype, maybe with werkzurg or something


class UrlUpload(BaseModel):
    url: str
    # I am going to be removing the ability for overloading metadata
    # title: str | None = None


class File(BaseModel):
    id: any  # TODO: figure out a better type for this UUID :/
    url: str
    title: str | None

    @validator("id")
    def validate_uuid(cls, value):
        if value:
            return str(value)
        return value


class FileUpdate(BaseModel):
    url: str | None = None
    title: str | None = None


class FileCreate(BaseModel):
    url: str | None = None
    title: str | None = None


# litestar only
class FileController(Controller):
    """File Controller"""

    dependencies = {"files_repo": Provide(provide_files_repo)}

    @get(path="/files/{file_id:uuid}")
    async def get_file(
        self,
        files_repo: FileRepository,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
    ) -> File:
        obj = files_repo.get(file_id)
        return File.model_validate(obj)

    @get(path="/files/all")
    async def get_all_files(
        self, files_repo: FileRepository, limit_offset: LimitOffset, request: Request
    ) -> list[File]:
        """List files."""
        results = await files_repo.list()
        type_adapter = TypeAdapter(list[File])
        return type_adapter.validate_python(results)

    @post(path="files/upload")
    async def upload_file(self) -> File:
        pass

    from crawleringest.docingest import DocumentIngester

    @post(path="/links/add")
    async def add_file(
        self, files_repo: FileRepository, data: FileCreate, request: Request
    ) -> File:
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
            file=read(raw_file_path),
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
        return File.model_validate(new_file)
    
    from docprocessing.extractmarkdown import MarkdownExtractor
    from docprocessing.genextras import GenerateExtras
    @patch(path="/docproc/{file_id:uuid}")
    async def process_File(
        self,
        files_repo: FileRepository,
        data: FileUpdate,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
    ) -> File:
        """Process a File."""
        obj = files_repo.get(file_id)
        current_stage = obj.stage
        mdextract = MarkdownExtractor()        
        genextras = GenerateExtras()
        if current_stage == "stage0":
            response_code, response_message = (422, "Failure in stage 0: Document was incorrectly added to database, try readding it again." )
        if current_stage == "stage1":
            try:
                processed_original_text = mdextract.process_raw_document_into_untranslated_text(obj.path, obj.metadata)
                obj.original_text = processed_original_text
                current_stage == "stage2"
            except:
                response_code, response_message = (422, "failure in stage 1: document was unable to be converted to markdown," )
        if current_stage == "stage2":
            try:
                processed_english_text = mdextract.convert_text_into_eng(obj.original_text, obj.lang)
                obj.english_text = processed_english_text
                current_stage == "stage3"
            except:
                response_code, response_message = (422, "failure in stage 2: document was unable to be translated to english." )
        if current_stage == "stage3":
            try:
                # TODO : Chunk and throw into chroma
                current_stage == "completed"
            except:
                response_code, response_message = (422, "failure in stage 2: document was unable to be translated to english." )
        if current_stage == "completed":
            response_code, response_message = (200, "Document Fully Processed." )
        newobj = files_repo.update(obj)
        files_repo.session.commit()
        return File.model_validate(newobj) # TODO : Return Response code and response message
    
    @patch(path="/files/{file_id:uuid}")
    async def update_File(
        self,
        files_repo: FileRepository,
        data: FileUpdate,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
    ) -> File:
        """Update a File."""
        raw_obj = data.model_dump(exclude_unset=True, exclude_none=True)
        raw_obj.update({"id": file_id})
        obj = files_repo.update(FileModel(**raw_obj))
        files_repo.session.commit()
        return File.model_validate(obj)
    
    @delete(path="/files/{file_id:uuid}")
    async def delete_file(
        self,
        files_repo: FileRepository,
        file_id: UUID = Parameter(title="File ID", description="File to retieve"),
    ) -> None:
        _ = files_repo.delete(files_repo)
        files_repo.session.commit()
