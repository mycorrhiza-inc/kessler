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
  id: number;
  name: string;
  description: string;
}

export interface Action {
  id: number;
  organization: Organization;
  description: string;
  date: Date;
}

export interface Battle {
  id: number;
  title: string;
  description: string;
  parentBattleId?: number;
  childBattleIds?: number[];
  actions: Action[];
}
// Example data

export const exampleOrganizations: Organization[] = [
  {
    id: 1,
    name: "Unnamed Affinity Group from Denver",
    description: "A group fighting for climate justice.",
  },
  {
    id: 2,
    name: "Public Trust Doctrine Advocates",
    description: "An organization using legal means to fight climate change.",
  },
  {
    id: 3,
    name: "Student BDS Movement",
    description:
      "Students promoting Boycott, Divestment, and Sanctions against apartheid in Israel.",
  },
];

export const exampleActions: Action[] = [
  {
    id: 1,
    organization: exampleOrganizations[0],
    description: "Performed Banner Drop in Glenwood Springs Opposing UBR",
    date: new Date("2023-01-15"),
  },
  {
    id: 2,
    organization: exampleOrganizations[1],
    description: "Filed Juliana vs United States lawsuit",
    date: new Date("2023-02-20"),
  },
  {
    id: 3,
    organization: exampleOrganizations[2],
    description: "Organized Auaria Protest",
    date: new Date("2023-03-05"),
  },
];

export const exampleBattles: Battle[] = [
  {
    id: 1,
    title: "Stopping New Fossil Fuel Development in the USA",
    description: "Efforts to halt fossil fuel projects in the USA.",
    childBattleIds: [2],
    actions: [],
  },
  {
    id: 2,
    title: "Stopping the Uinta Basin Railway",
    description: "Efforts to stop the Uinta Basin Railway project.",
    parentBattleId: 1,
    childBattleIds: [3],
    actions: [],
  },
  {
    id: 3,
    title:
      "Getting the Mayor of Greenwood Springs to publicly condemn the project",
    description: "Local level effort to gain political support against UBR.",
    parentBattleId: 2,
    actions: [exampleActions[0]],
  },
  {
    id: 4,
    title: "Fight Climate Change with Judicial System",
    description: "Using legal means to combat climate change.",
    childBattleIds: [5],
    actions: [],
  },
  {
    id: 5,
    title:
      "Use the Public Trust Doctrine/ to Establish a Constitutional Right to a Livable Climate",
    description: "Legal strategy to establish environmental rights.",
    parentBattleId: 4,
    childBattleIds: [6],
    actions: [],
  },
  {
    id: 6,
    title: "Juliana vs United States",
    description:
      "Landmark lawsuit to secure a constitutional right to a livable climate.",
    parentBattleId: 5,
    actions: [exampleActions[1]],
  },
  {
    id: 7,
    title: "Use BDS as a tool of stopping apartheid in Israel",
    description: "Promoting BDS to fight apartheid policies.",
    childBattleIds: [8],
    actions: [],
  },
  {
    id: 8,
    title: "Use Student Demonstrations to Encourage Divestment from Colleges",
    description: "Encouraging colleges to divest through student activism.",
    parentBattleId: 7,
    childBattleIds: [9],
    actions: [],
  },
  {
    id: 9,
    title: "Auaria Protest",
    description: "A specific protest organized by students.",
    parentBattleId: 8,
    actions: [exampleActions[2]],
  },
];
