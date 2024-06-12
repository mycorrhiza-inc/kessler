import numpy as np
from typing import List, Union, Optional
from datetime import datetime

from litestar.config.app import AppConfig
from litestar.plugins import InitPluginProtocol
from litestar.di import Provide

from pydantic import BaseModel

import lancedb
from lancedb.db import AsyncConnection
from lancedb.pydantic import LanceModel


class LanceDBAsyncConfig:
    uri: str = ""

    def __init__(cls, uri: str = ""):
        cls.uri = uri

    async def get_conn(cls) -> AsyncConnection:
        # TODO: experiment with read_consistency_interval
        return await lancedb.connect_async(uri=cls.uri)

    async def create_all(cls, tables: List[LanceModel]):
        pass


class LanceDBInitPlugin(InitPluginProtocol):
    config: LanceDBAsyncConfig

    def __init__(cls, config=LanceDBAsyncConfig):
        cls.config = config

    async def on_app_init(cls, app_config: AppConfig) -> AppConfig:
        conn = await cls.config.get_conn()
        app_config.dependencies["lancedb"] = Provide(conn)


lancedb_config = LanceDBAsyncConfig(uri="/tmp/lancedb/", workers=2)


async def provide_async_lance_repo():
    pass


lancedb_plugin = LanceDBInitPlugin
