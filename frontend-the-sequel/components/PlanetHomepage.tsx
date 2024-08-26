import React from "react";

// This app will have a search bar in the middle of the page, then have a bunch of article previews below the box, could you make seperate components for each element?
// Please add any imports you need, but try to keep it to the joy material ui library
import {
  Box,
  TextField,
  Typography,
  Container,
  Grid,
  Paper,
} from "@mui/material";
import { article } from "@/interfaces";

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
