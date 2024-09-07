class OrganizationSchema(BaseModel):
    id: UUID
    name: str
    description: Optional[str]


class IndividualSchema(BaseModel):
    id: UUID
    name: str
    current_org: Optional[UUID]


class Faction(BaseModel):
    name: str
    description: str
    orgs: List[OrganizationSchema]


class EncounterSchema(BaseModel):
    id: UUID
    name: str
    created_at: datetime
    document_set: List[FileSchema]
    description: str
    factions: List[Faction]
