import useGraphStore, { GraphState } from "../utils/GraphUtilities";
import { Node, Edge } from "reactflow";
import axios from "axios";
import { useState } from "react";
import { shallow } from "zustand/shallow";
import * as React from "react";

import { usePathname, useRouter } from "next/navigation";

// mui elements
import Box from "@mui/joy/Box";
import ListItemDecorator from "@mui/joy/ListItemDecorator";
import Tabs from "@mui/joy/Tabs";
import TabList from "@mui/joy/TabList";
import Tab, { tabClasses } from "@mui/joy/Tab";

// icons
import HomeRoundedIcon from "@mui/icons-material/HomeRounded";
import FavoriteBorder from "@mui/icons-material/FavoriteBorder";
import Search from "@mui/icons-material/Search";
import Person from "@mui/icons-material/Person";

const selector = (state: GraphState) => ({
  setNodes: state.setNodes,
  setEdges: state.setEdges,
  generateRandomNodes: state.generateRandomNodes,
  generateDebugNodeGraph: state.generateDebugNodeGraph,
  // ClusterNodesAround: state.ClusterNodesAround,
  ClusterAroundNode: state.ClusterAroundNode,
});
const ToolBar2 = () => {
  const {
    setNodes,
    setEdges,
    generateRandomNodes,
    generateDebugNodeGraph,
    // ClusterNodesAround,
    ClusterAroundNode,
  } = useGraphStore(selector, shallow);
  const [numToGenerate, setGenerate] = useState(0);
  const [toggle, setToggle] = useState(false);

  const RequestGraph = async (_id: string) => {
    console.log("Fetching graph from backend");
    try {
      const { data, status } = await axios.get(
        "http://localhost:5000/api/graph/",
        {
          headers: {
            "Content-Type": "application/json",
            "Access-Control-Allow-Origin": "http://localhost",
          },
        }
      );
      if (status === 200) {
        console.log("Request succeeded");
        console.log(data);
        setNodes(data.graph.nodes);
        setEdges(data.graph.edges);
      } else {
        console.log(`Request failed with status ${status}:\n data: ${data}`);
      }
    } catch (e: any) {
      console.log(`Request failed with error:\n data: ${e}`);
    }
  };
  return (
    <div className="chat-bar">
      <div id="toolbar" className="flex-item">
        {/* this will contain all of the files used in the current notebook */}
        <div id="toolbar-toggle">
          <button onClick={() => setToggle(!toggle)}>Toggle Toolbar</button>
        </div>
        {toggle && (
          <>
            <div id="fetch-graph" className="toolbar-item">
              <br />
              <button onClick={() => RequestGraph("")}>Fetch Graph</button>
            </div>
            <div id="add-url" className="toolbar-item">
              <button>Add Url</button>
            </div>
            <div id="add-file" className="toolbar-item">
              <button>Add File</button>
            </div>
            <div id="add-random-nodes" className="toolbar-item">
              <button
                onClick={() => {
                  generateRandomNodes(numToGenerate);
                }}
              >
                Generate Random Nodes
              </button>
              <input
                type="number"
                onChange={(e) => {
                  let i = parseInt(e.target.value);
                  i == undefined ? setGenerate(0) : setGenerate(i);
                }}
              />
            </div>
            <div id="add-starting-nodes" className="toolbar-item">
              <button
                onClick={() => {
                  generateDebugNodeGraph();
                }}
              >
                Generate Debug Nodes
              </button>
            </div>
            <div className="toolbar-item">
              <button
                onClick={() => {
                  ClusterAroundNode(`${0}`);
                }}
              >
                Cluster Nodes
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
};

const ToolBar = () => {
  const pathname = usePathname();
  const [index, setIndex] = React.useState(0);
  const colors = ["primary", "danger", "success", "warning"] as const;
  return (
    <Box
      sx={{
        flexGrow: 1,
        m: -3,
        p: 4,
        borderTopLeftRadius: "12px",
        borderTopRightRadius: "12px",
        bgcolor: `${colors[index]}.500`,
      }}
      className="toolbar"
    >
      <Tabs
        size="lg"
        aria-label="Bottom Navigation"
        value={index}
        onChange={(event, value) => setIndex(value as number)}
        sx={(theme) => ({
          p: 1,
          borderRadius: 16,
          maxWidth: 400,
          mx: "auto",
          boxShadow: theme.shadow.sm,
          "--joy-shadowChannel": theme.vars.palette[colors[index]].darkChannel,
          [`& .${tabClasses.root}`]: {
            py: 1,
            flex: 1,
            transition: "0.3s",
            fontWeight: "md",
            fontSize: "md",
            [`&:not(.${tabClasses.selected}):not(:hover)`]: {
              opacity: 0.7,
            },
          },
        })}
      >
        <TabList
          variant="plain"
          size="sm"
          disableUnderline
          sx={{ borderRadius: "lg", p: 0 }}
        >
          <Tab
            disableIndicator
            orientation="vertical"
            {...(index === 0 && { color: colors[0] })}
          >
            <ListItemDecorator>
              <HomeRoundedIcon />
            </ListItemDecorator>
            Home
          </Tab>
          <Tab
            disableIndicator
            orientation="vertical"
            {...(index === 1 && { color: colors[1] })}
          >
            <ListItemDecorator>
              <FavoriteBorder />
            </ListItemDecorator>
            Likes
          </Tab>
          <Tab
            disableIndicator
            orientation="vertical"
            {...(index === 2 && { color: colors[2] })}
          >
            <ListItemDecorator>
              <Search />
            </ListItemDecorator>
            Search
          </Tab>
          <Tab
            disableIndicator
            orientation="vertical"
            {...(index === 3 && { color: colors[3] })}
          >
            <ListItemDecorator>
              <Person />
            </ListItemDecorator>
            Profile
          </Tab>
        </TabList>
      </Tabs>
    </Box>
  );
};
export default ToolBar;
