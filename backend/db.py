from sqlalchemy.orm import DeclarativeBase
from pydantic import BaseModel as _BaseModel


class BaseModel(_BaseModel):
    """Extends Pydantics BaseModel to enable ORM mode"""

    model_config = {"from_attributes": True, "arbitrary_types_allowed": True}
