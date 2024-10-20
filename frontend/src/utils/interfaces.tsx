export interface article {
  title: string;
  summary: string;
  text: string;
  sources: documentInfo[];
}

export interface documentInfo {
  title: string;
  display_text: string;
  summary: string;
  metadata: Object;
}
export const exampleArticle: article = {
  title: "The Evolution of AI",
  summary: "An overview of the advancements in artificial intelligence.",
  text: "Artificial intelligence has rapidly evolved over the past few decades, transforming various industries and creating new opportunities...",
  sources: [
    {
      title: "The Rise of Machine Learning",
      display_text: "Machine Learning advancements",
      summary:
        "An in-depth look at how machine learning has become a cornerstone of AI development.",
      metadata: {
        author: "Jane Doe",
        publicationDate: "2021-05-15",
        journal: "Tech Journal",
      },
    },
    {
      title: "Neural Networks Explained",
      display_text: "Understanding Neural Networks",
      summary:
        "A comprehensive guide to how neural networks function and their applications.",
      metadata: {
        author: "John Smith",
        publicationDate: "2020-11-01",
        journal: "AI Review",
      },
    },
  ],
};

export const exampleDocumentInfo: documentInfo = {
  title: "Deep Learning Innovations",
  display_text: "Innovations in Deep Learning",
  summary:
    "Exploring the latest innovations in deep learning techniques and their impact.",
  metadata: {
    author: "Alice Johnson",
    publicationDate: "2022-01-10",
    journal: "Deep Learning Today",
  },
};

// Interfaces for new thing:
//
//

export interface Organization {
  id: string;
  name: string;
  description: string;
}

export interface Action {
  id: string;
  organization: Organization;
  description: string;
  date: Date;
}

export interface Faction {
  id: string;
  title: string; // Something like "Opponents", "Proponents", "Third Party Intervenors". For more complicated stuff like a climate bill with nuclear energy, it might be something like "Proponents of bill", "Want nuclear gone from bill", "Want Solar Investments Gone from Bill", and "Completely Opposed".
  description: string;
  organizations: Organization[];
}

export interface Battle {
  id: string;
  title: string;
  description: string;
  parentBattleIds?: number[];
  childBattleIds?: number[];
  actions: Action[];
  factions: Faction[];
}
// Example data

export const exampleOrganizations: Organization[] = [
  {
    id: "us-govt",
    name: "US Government",
    description: "A group fighting for climate justice.",
  },
  {
    id: "oct",
    name: "Our Childrens Trust",
    description: "An organization using legal means to fight climate change.",
  },
  {
    id: "naacp",
    name: "National Association for the Advancement of Colored People",
    description: "Test Description",
  },
  {
    id: "shell",
    name: "Shell Oil Company",
    description: "Description for Shell Oil Company",
  },
  {
    id: "bp",
    name: "BP",
    description: "Description for BP",
  },
  {
    id: "exxon",
    name: "ExxonMobil",
    description: "Description for ExxonMobil",
  },
];

export const exampleActions: Action[] = [
  {
    id: "action-oct-1",
    organization: exampleOrganizations[1],
    description: "Filed Juliana vs United States lawsuit",
    date: new Date("2015-02-20"),
  },
  {
    id: "action-oct-2",
    organization: exampleOrganizations[1],
    description: "Juliana Recived Bad Ruling from 9th Circuit",
    date: new Date("2019-02-20"),
  },
  {
    id: "action-oct-2",
    organization: exampleOrganizations[1],
    description: "Juliana Recived Bad Ruling from 9th Circuit Part 2",
    date: new Date("2023-02-20"),
  },
];

export const exampleFactions: Faction[] = [
  {
    id: "juliana-proponents",
    title: "Proponents of the lawsuit",
    description: "Groups who support the Juliana vs United States lawsuit.",
    organizations: [exampleOrganizations[1], exampleOrganizations[2]],
  },
  {
    id: "juliana-opponents",
    title: "Opponents of the lawsuit",
    description: "Groups who are against the Juliana vs United States lawsuit.",
    organizations: [exampleOrganizations[0]],
  },
  {
    id: "fossil-fuel-lobby",
    title: "Private Oil Companies",
    description:
      "Companies who intervened against the Juliana vs United States lawsuit, but dropped out to avoid discovery requests",
    organizations: [
      exampleOrganizations[3],
      exampleOrganizations[4],
      exampleOrganizations[5],
    ],
  },
];

export const exampleBattles: Battle[] = [
  {
    id: "battle-juliana-vs-us",
    title: "Juliana vs United States",
    description:
      "Landmark lawsuit to secure a constitutional right to a livable climate.",
    parentBattleIds: [5], // Assuming battle with id 5 exists
    childBattleIds: [],
    actions: [exampleActions[1]],
    factions: exampleFactions,
  },
];

export interface extraProperties {
  match_name: string;
  match_source: string;
  match_doctype: string;
  match_docket_id: string;
  match_document_class: string;
  match_author: string;
}
export const emptyExtraProperties: extraProperties = {
  match_name: "",
  match_source: "",
  match_doctype: "",
  match_docket_id: "",
  match_document_class: "",
  match_author: "",
};

export interface Filters {
  author: string;
  date: string;
  docket_id: string;
  doctype: string;
  lang: string;
  language: string;
  source: string;
  title: string;
}

export interface SearchRequest {
  query: string;
  filters: Filters;
}
