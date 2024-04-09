# litestar only
from uuid import UUID

from litestar import Controller, Request


from litestar.handlers.http_handlers.decorators import get, post, delete
from litestar.params import Parameter
from litestar.di import Provide
from litestar.repository.filters import LimitOffset

from pydantic import TypeAdapter


from models.resources import Resource


class ResourceController(Controller):
    """Resource Controller"""

    dependencies = {"resource": Provide(provide_resource_repo)}

    @get(path="/resource/{Resource_id:uuid}")
    async def get_Resource(
        self,
        Resources_repo: ResourceRepository,
        Resource_id: UUID = Parameter(
            title="Resource ID", description="Resource to retieve"
        ),
    ) -> Resource:
        obj = Resources_repo.get(Resource_id)
        return Resource.model_validate(obj)

    @get(path="/resources/all")
    async def get_all_Resources(
        self, limit_offset: LimitOffset, request: Request
    ) -> list[Resource]:
        """List Resources."""
        Resource.list()
        results = await Resource.list()
        type_adapter = TypeAdapter(list[Resource])
        return type_adapter.validate_python(results)

    @post(path="resources/upload")
    async def upload_Resource(self) -> Resource:
        pass

    @delete(path="/resources/{resource_id:uuid}")
    async def delete_Resource(
        self,
        resource_id: UUID = Parameter(
            title="Resource ID", description="Resource to retieve"
        ),
    ) -> None:
        _ = Resource.delete(resource_id)
        Resource.commit()
