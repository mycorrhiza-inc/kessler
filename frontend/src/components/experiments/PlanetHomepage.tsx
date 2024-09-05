"use client";
import React from "react";
import {
  Box,
  TextField,
  Typography,
  Container,
  Grid,
  Paper,
} from "@mui/material";
import {
  article,
  exampleBattles,
  Organization,
  Battle,
  Faction,
  Action,
} from "@/utils/interfaces";

import TestJuristictionViewer from "./JuristictionViewer";

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

import Calendar from "react-calendar";
// Interfaces for new thing:
//
//

// export interface Organization {
//   id: string;
//   name: string;
//   description: string;
// }
//
// export interface Action {
//   id: string;
//   organization: Organization;
//   description: string;
//   date: Date;
// }
//
// export interface Faction {
//   id: string;
//   title: string; // Something like "Opponents", "Proponents", "Third Party Intervenors". For more complicated stuff like a climate bill with nuclear energy, it might be something like "Proponents of bill", "Want nuclear gone from bill", "Want Solar Investments Gone from Bill", and "Completely Opposed".
//   description: string;
// }

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
    <Typography variant="body2">{action.date.toDateString()}</Typography>
    <Typography variant="caption">{action.description}</Typography>
  </Paper>
);

const BattleCardPreview: React.FC<{ battle: Battle }> = ({ battle }) => (
  <Paper elevation={3} style={{ padding: "20px", marginBottom: "20px" }}>
    <Typography variant="h4">{battle.title}</Typography>
    <Typography variant="body1">{battle.description}</Typography>
  </Paper>
);
//

const FactionBox: React.FC<{ faction: Faction; color: string }> = ({
  faction,
  color,
}) => (
  <Paper
    elevation={3}
    style={{ padding: "20px", marginBottom: "20px", backgroundColor: color }}
  >
    <Typography variant="h5">{faction.title}</Typography>
    <Typography>{faction.description}</Typography>
    <Box mt={2}>
      {faction?.organizations?.length > 0 &&
        faction.organizations.map((org) => (
          <OrganizationCard key={org.id} organization={org} />
        ))}
    </Box>
  </Paper>
);

const BattlePage: React.FC<{ battle: Battle; childBattles: Battle[] }> = ({
  battle,
  childBattles,
}) => {
  const factionColors = [
    "oklch(87% 0.1 0)",
    "oklch(87% 0.1 200)",
    "oklch(87% 0.1 140)",
    "oklch(87% 0.1 80)",
  ]; // Example colors
  console.log(battle);
  console.log(childBattles);

  return (
    <Paper elevation={3} style={{ padding: "20px", borderRadius: "15px" }}>
      <Typography variant="h4">{battle.title}</Typography>
      <Typography variant="body1">{battle.description}</Typography>

      {battle?.actions?.length > 0 && (
        <Grid container spacing={2} mt={2}>
          {battle.actions
            .sort((a, b) => a.date.getTime() - b.date.getTime())
            .map((action) => (
              <Grid item xs={12} sm={6} md={4} key={action.id}>
                <ActionCard action={action} />
              </Grid>
            ))}
        </Grid>
      )}
      {childBattles?.length > 0 && (
        <Grid container spacing={2} mt={2}>
          {childBattles.map((action) => (
            <Grid item xs={12} sm={6} md={4} key={action.id}>
              <BattleCardPreview battle={action} />
            </Grid>
          ))}
        </Grid>
      )}
      {battle?.factions?.length > 0 && (
        <Grid container spacing={2} mt={2}>
          {battle.factions.map((faction, index) => (
            <Grid item xs={12} sm={6} md={4} key={faction.id}>
              <FactionBox
                faction={faction}
                color={factionColors[index % factionColors.length]}
              />
            </Grid>
          ))}
        </Grid>
      )}

      <Grid container spacing={2} mt={2}>
        <Grid item xs={12}>
          <Box
            display="flex"
            justifyContent="space-between"
            alignItems="stretch"
          >
            <Box width="50%">
              <Typography variant="h6" gutterBottom>
                Calendar
              </Typography>
              <Box height="100%">
                <Calendar />
              </Box>
            </Box>
            <Box width="50%">
              <Typography variant="h6" gutterBottom>
                Test Juristiction Viewer
              </Typography>
              <Box height="100%">
                <TestJuristictionViewer />
              </Box>
            </Box>
          </Box>
        </Grid>
      </Grid>
    </Paper>
  );
};
const PlanetStartPage: React.FC<PlanetStartPageProps> = ({ articles }) => (
  <Container>
    <BattlePage battle={exampleBattles[0]} childBattles={[]}></BattlePage>
  </Container>
);

// 2.71828182845904523536028747135266249775724709369995957496696762772407663035354759457138217852516642742746
export default PlanetStartPage;
