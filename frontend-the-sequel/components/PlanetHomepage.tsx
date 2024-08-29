"use client";
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

const BattlePage: React.FC<{ battle: Battle; childBattles: Battle[] }> = ({
  battle,
  childBattles,
}) => (
  <Paper elevation={3} style={{ padding: "20px", marginBottom: "20px" }}>
    <Typography variant="h4">{battle.title}</Typography>
    <Typography variant="body1">{battle.description}</Typography>
    {battle.actions.map((action) => (
      <ActionCard key={action.id} action={action} />
    ))}
    {childBattles && (
      <Box mt={2}>
        {childBattles.map((childBattleId) => (
          <BattleCardPreview key={childBattle.id} battle={childBattle} />
        ))}
      </Box>
    )}
  </Paper>
);

const BattleCardPreview: React.FC<{ battle: Battle }> = ({ battle }) => (
  <Paper elevation={3} style={{ padding: "20px", marginBottom: "20px" }}>
    <Typography variant="h4">{battle.title}</Typography>
    <Typography variant="body1">{battle.description}</Typography>
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
