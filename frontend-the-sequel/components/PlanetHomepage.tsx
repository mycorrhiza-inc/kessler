import React from "react";

// This app will have a search bar in the middle of the page, then have a bunch of article previews below the box, could you make seperate components for each element?
// Please add any imports you need, but try to keep it to the joy material ui library
// ## Datatypes
//
// Battles - Fundamental Unit Type, Every Battle is has child battles, that are each battles occurring on a smaller more local level that are directly dependant for winning the parent battle. For one example, a chain might look like this.
//
// Stopping New Fossil Fuel Development in the USA > Stopping the Uinta Basin Railway > Getting the Mayor of Greenwood Springs to publicly condom the project.
// or
// Fight Climate Change with Judicial System > Use the Public Trust Doctrine/ to Establish a Constitutional Right to a Livable Climate > Juliana vs United States
// or
// Use BDS as a tool of stopping apartheid in Israel > Use Student Demonstrations to Encourage Divestement from Colleges > Auaria Protest
// You could potentially have a parent battle such as "Solve Climate Change", but I think that even in the highest levels it is essential to pair every societal problem with an explicit theory of change. Stuff like "End Capitalism" is so broad that its impossible for anyone to know what to do next. But explicitly framing every battle as.
//
// 1. In context of some large problem.
//
// 2. A method that people are using to try to solve that problem.
//
// Is a good framework for analysis for NGOs/Companies, since it actively forces you to consider theories of change, in connection with every campaign that they are working on. And its also great for beginners, since every problem is immediately bundled with actions that one can take to solve it.
//
// Battles with no children/"leaf battles" should ideally be similar to "SMART Goals", namely they should be narrow and well defined, with measurable outcomes for success and where the next steps for any party are reasonable to find.
//
// Likewise the main point of battles with children / "branch battles" is not to provide guidance as to what political action to take, but to try and provide context for each of the landscape that each of the leaf battles are taking place in. As well as providing a location for people to find actions they want to participate in, and network and build community.
//
// Each Battle should have the following attributes:
//
// - A list of "actions" taken by each organisation to try to w in the battle. (Maybe "skirmishes" would be a good name for these) For instance taking the above example, a skirmish would be "Unnamed Affinity Group from Denver, preforms Banner Drop in Glenwood Springs Opposing UBR". (Notifications for these sounds amazing, since if we could ever figure out a way to measure the impact of skirmishes, you could use that as a filter for notifications.)
//
// As a result you should also have an interface for organisations and associated data, and a type for actions/skirmishes and their data.
//
//
import {
  Box,
  TextField,
  Typography,
  Container,
  Grid,
  Paper,
} from "@mui/material";
import { article } from "@/interfaces";

import TestViewer from "./JuristictionViewer";

interface PlanetStartPageProps {
  articles: article[];
}
interface Organization {
  id: number;
  name: string;
  description: string;
}

interface Action {
  id: number;
  organization: Organization;
  description: string;
  date: Date;
}

interface Battle {
  id: number;
  title: string;
  description: string;
  parentBattleId?: number;
  childBattleIds?: number[];
  actions: Action[];
}
// Example data

const organizations: Organization[] = [
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

const actions: Action[] = [
  {
    id: 1,
    organization: organizations[0],
    description: "Performed Banner Drop in Glenwood Springs Opposing UBR",
    date: new Date("2023-01-15"),
  },
  {
    id: 2,
    organization: organizations[1],
    description: "Filed Juliana vs United States lawsuit",
    date: new Date("2023-02-20"),
  },
  {
    id: 3,
    organization: organizations[2],
    description: "Organized Auaria Protest",
    date: new Date("2023-03-05"),
  },
];

const battles: Battle[] = [
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
    actions: [actions[0]],
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
    actions: [actions[1]],
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
    actions: [actions[2]],
  },
];

const SearchBar: React.FC = () => (
  <Box display="flex" justifyContent="center" mt={5}>
    <TextField label="Search..." variant="outlined" fullWidth />
  </Box>
);

const ArticlePreview: React.FC<{ title: string; summary: string }> = ({
  title,
  summary,
}) => (
  <Grid item xs={12} md={6}>
    <Paper elevation={3} style={{ padding: "20px", marginBottom: "20px" }}>
      <Typography variant="h6">{title}</Typography>
      <Typography variant="body2">{summary}</Typography>
    </Paper>
  </Grid>
);

const ArticlesList: React.FC<{
  articles: { title: string; summary: string }[];
}> = ({ articles }) => (
  <Grid container spacing={3} mt={3}>
    {articles.map((article, index) => (
      <ArticlePreview
        key={index}
        title={article.title}
        summary={article.summary}
      />
    ))}
  </Grid>
);

const PlanetStartPage: React.FC<PlanetStartPageProps> = ({ articles }) => (
  <Container>
    <SearchBar />
    <ArticlesList articles={articles} />
    <TestViewer></TestViewer>
  </Container>
);

// 2.71828182845904523536028747135266249775724709369995957496696762772407663035354759457138217852516642742746
export default PlanetStartPage;

const OrganizationCard: React.FC<{ organization: Organization }> = ({
  organization,
}) => (
  <Paper elevation={3} style={{ padding: "20px", marginBottom: "20px" }}>
    <Typography variant="h6">{organization.name}</Typography>
    <Typography variant="body2">{organization.description}</Typography>
  </Paper>
);

const ActionCard: React.FC<{ action: Action }> = ({ action }) => (
  <Paper elevation={3} style={{ padding: "20px", marginBottom: "20px" }}>
    <Typography variant="h6">{action.organization.name}</Typography>
    <Typography variant="body2">{action.description}</Typography>
    <Typography variant="caption">{action.date.toDateString()}</Typography>
  </Paper>
);

const BattleCard: React.FC<{ battle: Battle }> = ({ battle }) => (
  <Paper elevation={3} style={{ padding: "20px", marginBottom: "20px" }}>
    <Typography variant="h4">{battle.title}</Typography>
    <Typography variant="body1">{battle.description}</Typography>
    {battle.actions.map((action) => (
      <ActionCard key={action.id} action={action} />
    ))}
    {battle.childBattles && (
      <Box mt={2}>
        {battle.childBattles.map((childBattleId) => (
          <BattleCard key={childBattle.id} battle={childBattle} />
        ))}
      </Box>
    )}
  </Paper>
);

const BattlesList: React.FC<{ battles: Battle[] }> = ({ battles }) => (
  <Grid container spacing={3} mt={3}>
    {battles.map((battle) => (
      <Grid item xs={12} key={battle.id}>
        <BattleCard battle={battle} />
      </Grid>
    ))}
  </Grid>
);
