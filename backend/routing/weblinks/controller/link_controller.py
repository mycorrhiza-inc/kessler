from uuid import UUID


from litestar import Controller, Request

from litestar.handlers.http_handlers.decorators import get, post, delete, patch
from litestar.params import Parameter
from litestar.di import Provide
from litestar.repository.filters import LimitOffset

from pydantic import TypeAdapter, validator

from db import BaseModel

from routing.links.dbm import LinkRepository, provide_Links_repo, LinkModel


class Link(BaseModel):
    id: any  # TODO: figure out a better type for this UUID :/
    url: str
    title: str | None

    @validator("id")
    def validate_uuid(cls, value):
        if value:
            return str(value)
        return value


class LinkUpdate(BaseModel):
    url: str | None = None
    title: str | None = None


class LinkCreate(BaseModel):
    url: str | None = None
    title: str | None = None


# litestar only
class LinkController(Controller):
    """Link Controller"""

    dependencies = {"Links_repo": Provide(provide_Links_repo)}

    @get(path="/links/{Link_id:uuid}")
    async def get_Link(
        self,
        Links_repo: LinkRepository,
        Link_id: UUID = Parameter(
            title="Link ID", description="Link to retieve"),
    ) -> Link:
        obj = Links_repo.get(Link_id)
        return Link.model_validate(obj)

    @get(path="/links/all")
    async def get_all_Links(
        self, Links_repo: LinkRepository, limit_offset: LimitOffset, request: Request
    ) -> list[Link]:
        """List Links."""
        results = await Links_repo.list()
        type_adapter = TypeAdapter(list[Link])
        return type_adapter.validate_python(results)

    @post(path="Links/upload")
    async def upload_Link(self) -> Link:
        pass

    @post(path="/links/add")
    async def add_Link(
        self, Links_repo: LinkRepository, data: LinkCreate, request: Request
    ) -> Link:
        request.logger.info("adding Links")
        request.logger.info(data)
        new_Link = LinkModel(url=data.url, title="")
        request.logger.info("new Link:{Link}".format(Link=new_Link.to_dict()))
        try:
            new_Link = await Links_repo.add(new_Link)
        except Exception as e:
            request.logger.info(e)
            return e
        request.logger.info("added Link!~")
        await Links_repo.session.commit()
        return Link.model_validate(new_Link)

    @patch(path="/links/{Link_id:uuid}")
    async def update_Link(
        self,
        Links_repo: LinkRepository,
        data: LinkUpdate,
        Link_id: UUID = Parameter(
            title="Link ID", description="Link to retieve"),
    ) -> Link:
        """Update a Link."""
        raw_obj = data.model_dump(exclude_unset=True, exclude_none=True)
        raw_obj.update({"id": Link_id})
        obj = Links_repo.update(LinkModel(**raw_obj))
        Links_repo.session.commit()
        return Link.model_validate(obj)

    @delete(path="/links/{Link_id:uuid}")
    async def delete_Link(
        self,
        Links_repo: LinkRepository,
        Link_id: UUID = Parameter(
            title="Link ID", description="Link to retieve"),
    ) -> None:
        _ = Links_repo.delete(Links_repo)
        Links_repo.session.commit()
