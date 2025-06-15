import { subdividedHueFromSeed } from "@/components/style/Pills/TextPills";

export interface OrgAutocompleteInfo {
  value: string;
  label: string;
  uuid: string;
}

export interface ConvoAutocompleteInfo {
  value: string;
  label: string;
  uuid: string;
}

export interface GeneralizedOption {
  value: string;
  label: string;
  color: string;
  // isFixed?: boolean;
  // isDisabled?: boolean;
}
export const orgAutocompleteToGeneralOption = (
  org: OrgAutocompleteInfo,
): GeneralizedOption => {
  return {
    value: org.value,
    label: org.label,
    color: subdividedHueFromSeed(org.uuid),
  };
};

export const convoAutocompleteToGeneralOption = (
  convo: ConvoAutocompleteInfo,
): GeneralizedOption => {
  return {
    value: convo.value,
    label: convo.label,
    color: subdividedHueFromSeed(convo.label),
  };
};

export const ConversationsAutocompleteList: ConvoAutocompleteInfo[] = [
  {
    uuid: "eb940a1c-8ef8-461d-b143-d545cb513eaf",
    value: "18-M-0084",
    label: "In the Matter of a Comprehensive Energy Efficiency Initiative.",
  },
  {
    uuid: "e913cfe5-6a84-449b-a918-0cf9bdc90df1",
    value: "17-F-0282",
    label:
      "Application of Alle-Catt Wind Energy LLC for a Certificate of Environmental Compatibility and Public Need Pursuant to Article 10 for a Proposed Wind Energy Project, Located in Allegany, Cattaraugus, and Wyoming Counties, New York, in the Towns of Arcade, Centerville, Farmersville, Freedom, and Rushford.",
  },
  {
    uuid: "d8d79012-5ae1-4d8f-8dd4-962fbd575efd",
    value: "15-E-0751",
    label: "In the Matter of the Value of Distributed Energy Resources",
  },
  {
    uuid: "2bbde660-8111-499f-a17f-6ca0ab765ad3",
    value: "16-F-0559",
    label:
      "Application of Bluestone Wind, LLC for a Certificate of Environmental Compatibility and Public Need Pursuant to Article 10 for Construction of the Bluestone Wind Farm Project Located in the Towns of Windsor and Sanford, Broome County.",
  },
  {
    uuid: "0c0571fa-fea3-44b9-b1b4-e9a3a39c8c2f",
    value: "18-T-0604",
    label:
      "Application of Deepwater Wind South Fork, LLC for a Certificate of Environmental Compatibility and Public Need for the Construction of Approximately 3.5 Miles of Submarine Export Cable from the New York State Territorial Waters Boundary to the South Shore of the Town of East Hampton in Suffolk County and Approximately 4.1 Miles of Terrestrial Export Cable from the South Shore of the Town of East Hampton to an Interconnection Facility with an Interconnection Cable Connecting to the Existing East Hampton Substation in the Town of East Hampton, Suffolk County.",
  },
  {
    uuid: "130014d9-7c42-4f24-af08-2d2e47e5bdcb",
    value: "19-E-0380",
    label:
      "Proceeding on Motion of the Commission as to the Rates, Charges, Rules and Regulations of Rochester Gas and Electric Corporation for Electric Service.",
  },
  {
    uuid: "cc575675-b3f7-4cc7-b04d-bb5e3b322356",
    value: "16-F-0205",
    label:
      "Application of Canisteo Wind Energy LLC for a Certificate of Environmental Compatibility and Public Need Pursuant to Article 10 for Construction of a Wind Energy Project in Steuben County.",
  },
  {
    uuid: "15018c5a-763a-4ca9-a50b-bdbbd8c9d570",
    value: "19-G-0309",
    label:
      "Proceeding on Motion of the Commission as to the Rates, Charges, Rules and Regulations of The Brooklyn Union Gas Company d/b/a National Grid NY for Gas Service.",
  },
  {
    uuid: "1dd52794-66b6-485a-b9ab-be18d1a8f135",
    value: "18-E-0138",
    label:
      "Proceeding on Motion of the Commission Regarding Electric Vehicle Supply Equipment and Infrastructure.",
  },
  {
    uuid: "d34f906e-78b9-4e8f-8b2d-2e0e617e22c6",
    value: "15-M-0566",
    label:
      "In the Matter of Revisions to Customer Service Performance Indicators Applicable to Gas and Electric Corporations",
  },
];
export const OrganizationsAutocompleteList: OrgAutocompleteInfo[] = [
  {
    label: "Public Service Commission",
    value: "d8dec2a3-cfac-43ae-8732-c2103aed8038",
    uuid: "d8dec2a3-cfac-43ae-8732-c2103aed8038",
  },
  {
    label: "New York State Department of Public Service",
    value: "0b544651-0226-4e0d-83af-184ef5aad4e5",
    uuid: "0b544651-0226-4e0d-83af-184ef5aad4e5",
  },
  {
    label: "New York State Electric & Gas Corporation",
    value: "13b07877-3021-44d2-a795-5576ae3c505f",
    uuid: "13b07877-3021-44d2-a795-5576ae3c505f",
  },
  {
    label: "Niagara Mohawk Power Corporation d/b/a National Grid",
    value: "f490bc1c-8150-400d-90af-7235b1e7d604",
    uuid: "f490bc1c-8150-400d-90af-7235b1e7d604",
  },
  {
    label: "Rochester Gas and Electric Corporation",
    value: "34ace092-b777-48e6-84a0-96d06eef8285",
    uuid: "34ace092-b777-48e6-84a0-96d06eef8285",
  },
  {
    label: "Consolidated Edison Company of New York, Inc.",
    value: "c7f1aca4-64bc-499d-844c-2aefdc4608d0",
    uuid: "c7f1aca4-64bc-499d-844c-2aefdc4608d0",
  },
  {
    label: "Central Hudson Gas & Electric Corporation",
    value: "bbc775e5-f6fd-4dbb-9403-c2d19c83984c",
    uuid: "bbc775e5-f6fd-4dbb-9403-c2d19c83984c",
  },
  {
    label: "Orange and Rockland Utilities, Inc.",
    value: "aeaea570-2ce5-4c2f-86b8-7800a74b22fa",
    uuid: "aeaea570-2ce5-4c2f-86b8-7800a74b22fa",
  },
  {
    label: "National Fuel Gas Distribution Corporation",
    value: "680a82de-3845-429d-b6ae-86ad82caa269",
    uuid: "680a82de-3845-429d-b6ae-86ad82caa269",
  },
  {
    label: "The Brooklyn Union Gas Company d/b/a National Grid NY",
    value: "800b9878-8193-4064-877f-0259a3db10e8",
    uuid: "800b9878-8193-4064-877f-0259a3db10e8",
  },
];
