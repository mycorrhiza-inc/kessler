use std::collections::HashMap;

use chrono::{DateTime, Utc};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use serde_with::skip_serializing_none;

use crate::common::{file_extension::FileExtension, hash::Blake2bHash};

#[derive(Serialize, Deserialize, Debug, JsonSchema, Clone, Hash, PartialEq, Eq)]
pub struct JurisdictionInfo {
    pub country: String,
    pub state: String,
    pub jurisdiction: String,
}
impl Default for JurisdictionInfo {
    fn default() -> Self {
        let unknown_static = "unknown";
        JurisdictionInfo {
            country: unknown_static.to_string(),
            state: unknown_static.to_string(),
            jurisdiction: unknown_static.to_string(),
        }
    }
}

impl JurisdictionInfo {
    pub fn new_usa(jurisdiction: &str, state: &str) -> Self {
        JurisdictionInfo {
            country: "usa".to_string(),
            state: state.to_string(),
            jurisdiction: jurisdiction.to_string(),
        }
    }
}

#[derive(Serialize, Deserialize, Debug, JsonSchema, Clone)]
pub struct CaseWithJurisdiction {
    pub case: GenericCase,
    pub jurisdiction: JurisdictionInfo,
}

#[skip_serializing_none]
#[derive(Serialize, Deserialize, Debug, JsonSchema, Clone)]
pub struct GenericAttachment {
    pub name: String,
    pub url: String,
    pub document_extension: Option<String>,
    pub extra_metadata: HashMap<String, serde_json::Value>,
    pub hash: Option<Blake2bHash>,
}

#[skip_serializing_none]
#[derive(Serialize, Deserialize, Debug, JsonSchema, Clone, Default)]
pub struct GenericFilingLegacy {
    pub name: String,
    pub filed_date: DateTime<Utc>,
    pub party_name: String,
    pub filing_type: String,
    pub description: String,
    pub attachments: Vec<GenericAttachment>,
    pub extra_metadata: HashMap<String, serde_json::Value>,
}

#[skip_serializing_none]
#[derive(Serialize, Deserialize, Debug, JsonSchema, Clone, Default)]
pub struct GenericFiling {
    pub name: String,
    pub filed_date: DateTime<Utc>,
    pub organization_authors: Vec<String>,
    pub individual_authors: Vec<String>,
    pub filing_type: String,
    pub description: String,
    pub attachments: Vec<GenericAttachment>,
    pub extra_metadata: HashMap<String, serde_json::Value>,
}

impl From<GenericFilingLegacy> for GenericFiling {
    fn from(value: GenericFilingLegacy) -> Self {
        GenericFiling {
            name: value.name,
            filed_date: value.filed_date,
            organization_authors: vec![value.party_name],
            individual_authors: vec![],
            filing_type: value.filing_type,
            description: value.description,
            attachments: value.attachments,
            extra_metadata: value.extra_metadata,
        }
    }
}

#[skip_serializing_none]
#[derive(Serialize, Deserialize, Debug, JsonSchema, Clone, Default)]
pub struct GenericCase {
    pub case_number: String,
    pub case_name: String,
    pub case_url: String,
    pub case_type: Option<String>,
    pub description: Option<String>,
    pub industry: Option<String>,
    pub petitioner: Option<String>,
    pub hearing_officer: Option<String>,
    pub opened_date: Option<DateTime<Utc>>,
    pub closed_date: Option<DateTime<Utc>>,
    pub filings: Vec<GenericFiling>,
    pub extra_metadata: HashMap<String, serde_json::Value>,
    pub indexed_at: DateTime<Utc>,
}
#[skip_serializing_none]
#[derive(Serialize, Deserialize, Debug, JsonSchema, Clone, Default)]
pub struct GenericCaseLegacy {
    pub case_number: String,
    pub case_name: String,
    pub case_url: String,
    pub case_type: Option<String>,
    pub description: Option<String>,
    pub industry: Option<String>,
    pub petitioner: Option<String>,
    pub hearing_officer: Option<String>,
    pub opened_date: Option<DateTime<Utc>>,
    pub closed_date: Option<DateTime<Utc>>,
    pub filings: Vec<GenericFilingLegacy>,
    pub extra_metadata: HashMap<String, serde_json::Value>,
    pub indexed_at: DateTime<Utc>,
}
impl From<GenericCaseLegacy> for GenericCase {
    fn from(value: GenericCaseLegacy) -> Self {
        GenericCase {
            case_number: value.case_number,
            case_name: value.case_name,
            case_url: value.case_url,
            case_type: value.case_type,
            description: value.description,
            industry: value.industry,
            petitioner: value.petitioner,
            hearing_officer: value.hearing_officer,
            opened_date: value.opened_date,
            closed_date: value.closed_date,
            filings: value.filings.into_iter().map(|f| f.into()).collect(),
            extra_metadata: value.extra_metadata,
            indexed_at: value.indexed_at,
        }
    }
}

#[derive(Serialize, Deserialize, Debug, Clone, Copy, JsonSchema)]
pub enum AttachmentTextQuality {
    #[serde(rename = "low")]
    Low,
    #[serde(rename = "high")]
    High,
}

#[derive(Serialize, Deserialize, Debug, Clone, JsonSchema)]
pub struct RawAttachmentText {
    pub quality: AttachmentTextQuality,
    pub language: String,
    pub text: String,
    pub timestamp: DateTime<Utc>,
}

#[skip_serializing_none]
#[derive(Serialize, Deserialize, Debug, JsonSchema, Clone)]
pub struct RawAttachment {
    pub hash: Blake2bHash,
    pub jurisdiction_info: JurisdictionInfo,
    pub name: String,
    pub extension: FileExtension,
    pub text_objects: Vec<RawAttachmentText>,
    pub date_added: chrono::DateTime<Utc>,
    pub date_updated: chrono::DateTime<Utc>,
    pub extra_metadata: Option<HashMap<String, String>>,
}
