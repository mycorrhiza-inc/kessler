import os
from pathlib import Path
from uuid import UUID

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


from typing import List, Optional, Dict


import logging
from models import utils
from logic.filelogic import process_fileid_raw
import asyncio
from litestar.contrib.sqlalchemy.base import UUIDBase

from sqlalchemy.ext.asyncio import AsyncSession

from logic.databaselogic import (
    QueryData,
    querydata_to_filters_strict,
    filters_docstatus_processing,
)

from util.gpu_compute_calls import get_total_connections
import random


class MiscController(Controller):
    """Miscellaneous Controller"""

    dependencies = {"files_repo": Provide(provide_files_repo)}

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