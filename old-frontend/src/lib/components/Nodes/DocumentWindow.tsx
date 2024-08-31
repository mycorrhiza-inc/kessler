import React from "react";
import { Card } from "antd";

const gridStyle: React.CSSProperties = {
  width: "25%",
  textAlign: "center",
};

const DocumentWindow: React.FC = () => (
  <div className="document-window">
    <Card title="Document Context">
      <Card.Grid style={gridStyle}>Content</Card.Grid>
      <Card.Grid hoverable={false} style={gridStyle}>
        Content
      </Card.Grid>
      <Card.Grid style={gridStyle}>Content</Card.Grid>
      <Card.Grid style={gridStyle}>Content</Card.Grid>
      <Card.Grid style={gridStyle}>Content</Card.Grid>
      <Card.Grid style={gridStyle}>Content</Card.Grid>
      <Card.Grid style={gridStyle}>Content</Card.Grid>
    </Card>
  </div>
);

export default DocumentWindow;
