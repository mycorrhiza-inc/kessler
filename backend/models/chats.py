from typing import List

from litestar.contrib.sqlalchemy.base import UUIDAuditBase

from sqlalchemy.orm import Mapped


from .utils import PydanticBaseModel


class MessageModel(UUIDAuditBase):
    __tablename__ = "message"
    sourceModel: Mapped[str]  # Human or ModelName
    Text: Mapped[str]


class ChatModel(UUIDAuditBase):
    pass


class MessageSchema(PydanticBaseModel):
    id: any
    body: str


class ChatSchema(PydanticBaseModel):
    id: any
    messages: List[MessageSchema]
