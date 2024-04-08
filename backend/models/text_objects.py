from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column
from litestar.contrib.sqlalchemy.base import UUIDAuditBase, AuditColumns



class TextObject(UUIDAuditBase):
    __tablename__ = "text_object"
    original_text: Mapped[str]
    en_text: Mapped[str | None]
    resource: Mapped[ForeignKey] = mapped_column()
