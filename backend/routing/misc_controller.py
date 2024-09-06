from litestar import Controller, Request, Response

from litestar.handlers.http_handlers.decorators import (
    get,
    post,
)
from litestar.events import listener


from litestar.params import Parameter
from litestar.di import Provide
from pydantic import TypeAdapter
from models.utils import PydanticBaseModel as BaseModel


from models.files import (
    FileSchema,
    FileModel,
    provide_files_repo,
    DocumentStatus,
    docstatus_index,
    FileRepository,
)


class MiscController(Controller):
    """Miscellaneous Controller"""

    dependencies = {"files_repo": Provide(provide_files_repo)}

    @get(path="/misc/test")
    async def basic_test(
        self,
    ) -> str:
        return "Test Successful"

    @get(path="/misc/example_csv")
    async def large_example_csv(
        self,
    ) -> str:
        rows = ["id,name,age"]
        for i in range(1, 1001):
            rows.append(f"{i},Name {i},{20 + (i % 30)}")
        csv_content = "\n".join(rows)
        return csv_content

    @get(path="/misc/allowable_fields")
    async def get_metadata(
        self,
        files_repo: FileRepository,
        request: Request,
    ) -> dict:
        # Find out way to refresh and generate from postgres then cache for duration of application.
        source_list = [
            "colorado-puc-electricity",
            "ny-puc-energyefficiency-filedocs",
            "personal",
        ]
        return {"stage": [s.value for s in DocumentStatus], "source": source_list}
