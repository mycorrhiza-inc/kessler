// "id": "scc7a_60_fuel",
// "name": "scc7a_60_fuel",
// "visualization_url": "https://nimbus.kessler.xyz/scc7a_60_fuel",
// "data_url": "https://nimbus.kessler.xyz/data/scc7a_60_fuel",
// "date": "10/08/2024",
export interface ModelType {
  id: string;
  name: string;
  visualization_url: string;
  data_url: string;
  date: string;
}

export interface FileType {
  id: string;
  url: string;
  name: string;
  doctype: string;
  stage: string;
  source: string;
  mdata: Object;
  display_text: string;
}

export interface TableLayout {
  columns: {
    key: string;
    label: string;
    width: string;
    enabled: boolean;
  }[];
  showExtraFeatures: boolean;
  showDisplayText: boolean;
}
