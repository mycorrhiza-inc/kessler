
I am trying to migrate some legacy code right now, typically this endpoing used to handle data from multiple different sources, but now all of those different data sources are handled and brought into a common format elsewhere. Namely these are accessible using api endpoints like:

```py
class GenericCase(BaseModel, extra=Extra.allow):
    """Model representing case data.

    Attributes:
        case_number (str): The unique case number.
        case_type (Optional[str]): The type of the case (e.g., civil, regulatory).
        description (Optional[str]): A detailed description of the case.
        industry (Optional[str]): The industry related to the case.
        petitioner (Optional[str]): The name of the petitioner in the case.
        hearing_officer (Optional[str]): The hearing officer for the case.
        opened_date (Optional[date]): The date the case was opened.
        closed_date (Optional[date]): The date the case was closed.
        filings (Optional[list[Filing]]): A list of filings associated with the case.
    """

    case_number: str
    case_name: str = ""
    case_url: str = ""
    case_type: Optional[str] = None
    description: Optional[str] = None
    industry: Optional[str] = None
    petitioner: Optional[str] = None
    hearing_officer: Optional[str] = None
    opened_date: Optional[RFC3339Time] = None
    closed_date: Optional[RFC3339Time] = None
    filings: Optional[list[GenericFiling]] = None
    extra_metadata: Dict[str, Any] = {}
    indexed_at: RFC3339Time = rfc_time_now()

    def model_post_init(self, __context: Any) -> None:
        if self.model_extra:
            self.extra_metadata.update(self.model_extra)

class GenericFiling(BaseModel, extra=Extra.allow):
    """Model representing filing data within a case.

    Attributes:
        filed_date (date): The date the filing was made.
        party_name (str): The name of the party submitting the filing.
        filing_type (str): The type of filing (e.g., brief, testimony).
        description (str): A description of the filing.
        attachments (Optional[list[Attachment]]): A list of associateda ttachments.
    """

    name: str = ""
    filed_date: RFC3339Time
    party_name: str
    filing_type: str
    description: str
    attachments: List[GenericAttachment] = []
    extra_metadata: Dict[str, Any] = {}

		 def model_post_init(self, __context: Any) -> None:
		     if self.model_extra:
		         self.extra_metadata.update(self.model_extra)

```
Could you refactor this project to use these 2 data types as the primary interface. Currently its using a type of ScraperInfoPayload in /ingest/tasks/schemas.go
Some of the openscraper types already exist in a pimative form in /ingest/openscrapers.go. So that might be helpful.
