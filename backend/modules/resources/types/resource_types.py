from db import BaseModel


class ResourceGet(BaseModel):
    id: any


class ResourceCreate(BaseModel):
    id: any


class Resource(BaseModel):
    id: any  # TODO: figure out a better type for this UUID :/
    type: str
