from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column
from litestar.contrib.sqlalchemy.base import UUIDAuditBase, AuditColumns


class TextResourceModel(AuditColumns):
    __tablename__ = "TextResource"
    # used to get Links from resource IDs
    resource_id = mapped_column(ForeignKey("resource.id"))
    # these are different so we can update link objects without regard to their resource id
    text_id = mapped_column(ForeignKey("text_object.id"), primary_key=True)


class TextObject(UUIDAuditBase):
    __tablename__ = "text_object"
    original_text: Mapped[str]
    en_text: Mapped[str | None]
    resource: Mapped[ForeignKey] = mapped_column()
