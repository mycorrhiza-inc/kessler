from uuid import UUID
from typing import Annotated


from litestar import Controller, Request

from litestar.handlers.http_handlers.decorators import \
    get, post, delete, patch, MediaType

from litestar.params import Parameter
from litestar.di import Provide
from litestar.repository.filters import LimitOffset
from litestar.datastructures import UploadFile
from litestar.enums import RequestEncodingType
from litestar.params import Body

from pydantic import TypeAdapter, BaseModel


from models.files import FileRepository, File, FileModel

# for testing purposese
emptyFile = FileModel(
    path="",
    doctype="",
    lang="",
    name="",
    stage="unprocessed",
    summary=None,
    short_summary=None)


class FileUpdate(BaseModel):
    message: str


class FileCreate(BaseModel):
    message: str


class FileUpload(BaseModel):
    message: str

# litestar only


class FileController(Controller):
    """File Controller"""

    dependencies = {"files_repo": Provide(FileModel.provide_repo)}

    @get(path="/files/{file_id:uuid}")
    async def get_file(
        self,
        files_repo: FileRepository,
        file_id: UUID = Parameter(
            title="File ID", description="File to retieve"),
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

    @post(path="/files/upload", media_type=MediaType.TEXT)
    async def handle_file_upload(
        self,
        files_repo: FileRepository,
        data: Annotated[UploadFile, Body(media_type=RequestEncodingType.MULTI_PART)],
    ) -> FileUpload:
        content = await data.read()
        newFileObj = emptyFile()
        newFileObj.name = data.filename

        obj = files_repo.add(newFileObj)

        # TODO: emit event for celery to process this file

        # return f"{newFileObj.name}, {content.decode()}"
        return obj

    @patch(path="/files/{file_id:uuid}")
    async def update_file(
        self,
        files_repo: FileRepository,
        data: FileUpdate,
        file_id: UUID = Parameter(
            title="File ID", description="File to retieve"),
    ) -> File:
        """Update a File."""
        raw_obj = data.model_dump(exclude_unset=True, exclude_none=True)
        raw_obj.update({"id": file_id})
        obj = files_repo.update(File(**raw_obj))
        files_repo.session.commit()
        return File.model_validate(obj)

    @delete(path="/files/{file_id:uuid}")
    async def delete_file(
        self,
        files_repo: FileRepository,
        file_id: UUID = Parameter(
            title="File ID", description="File to retieve"),
    ) -> None:
        _ = files_repo.delete(files_repo)
        files_repo.session.commit()
