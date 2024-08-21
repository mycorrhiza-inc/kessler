These are the types of metadata we should gaurentee, aka, if the field is there, our scraper should try to get it, and if the value is null, it should imply that the document actually has no metadata for that field.

"date" :Date of last document change/update

"organization" : Name of the organization that made the document

Sub Organization:

Sub Organization Type:

Participants

Author

Title

Document ID for Goverment Use

ID of the Docket the document is a part of

Internal ID tracking the docket proceedings, (should often be the gov id, + extra identifying info)

Type of Document (Report, Act, Decision, Public Comment, Bill, Recording, Etc)

```json
{
  "title": {
    "description": "The title of the document.",
    "example": "Environmental Impact Report for the Uinta Basin Railway"
  },
  "author": {
    "description": "The author of the document.",
    "example": "John Doe"
  },
  "date": {
    "description": "The date of the last document change or update.",
    "example": "2023-09-15"
  },
  "organization": {
    "description": "The name of the organization that created the document.",
    "example": "Environmental Protection Agency"
  },
  "sub_organization": {
    "description": "The name of any sub-organization involved.",
    "example": "Department of Water Quality"
  },
  "sub_organization_type": {
    "description": "The type of the sub-organization involved.",
    "example": "Government Department"
  },
  "participants": {
    "description": "The participants involved in creating or contributing to the document.",
    "example": ["Alice Smith", "Bob Johnson"]
  },
  "document_id_government": {
    "description": "The document ID for government use.",
    "example": "EPA-12345-2023"
  },
  "internal_id": {
    "description": "The full identifier used in the goverment system to track the document.",
    "example": "D-2023-45678"
  },
  "document_type": {
    "description": "The type of document (e.g., Report, Act, Decision, Public Comment, Bill, Recording, etc.).",
    "example": "Report"
  }
}
```

```

```
