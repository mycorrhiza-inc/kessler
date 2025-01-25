package admin

type DocumentMetadataCheck struct {
	NamedDocketID string `json:"named_docket_id"`
	DateString    string `json:"date_string"`
	Name          string `json:"name"`
	Extension     string `json:"extension"`
	AuthorString  string `json:"author_string"`
}

// {"id": "14f92e3e-6748-4f64-8c1a-1200547bf8e9", "url": "https://documents.dps.ny.gov/public/Common/ViewDoc.aspx?DocRefId={2071B08B-0000-C730-9754-1C39327CCC7D}", "date": "11/08/2023", "lang": "en", "title": "Staff Proposal on the Transition of Utility Reported Community-Scale Energy Usage Data", "author": "New York State Department of Public Service", "source": "New York State Department of Public Service", "authors": ["New York State Department of Public Service"], "language": "en", "docket_id": "20-M-0082", "extension": "pdf", "file_class": "Plans and Proposals", "item_number": "249", "author_email": "", "author_organisation": "New York State Department of Public Service"}
