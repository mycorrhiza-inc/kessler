import { useEffect } from "react";
import useGraphStore, { GraphState } from "../utils/GraphUtilities";
import { Node, Edge } from "reactflow";
import axios from "axios";
import { useState } from "react";
import { shallow } from "zustand/shallow";
import * as React from "react";

import { redirect, usePathname, useRouter } from "next/navigation";

// mui elements
import Box from "@mui/joy/Box";
import ListItemDecorator from "@mui/joy/ListItemDecorator";
import Tabs from "@mui/joy/Tabs";
import TabList from "@mui/joy/TabList";
import Tab, { tabClasses } from "@mui/joy/Tab";

// icons
import HomeRoundedIcon from "@mui/icons-material/HomeRounded";
import CollectionsBookmarkIcon from "@mui/icons-material/CollectionsBookmark";
import { SatelliteAlt } from "@mui/icons-material";
import Search from "@mui/icons-material/Search";
import SettingsIcon from "@mui/icons-material/Settings";
import { Tooltip } from "@mui/joy";

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
  const [index, setIndex] = React.useState(-1);
  const color = "success";
  const pathval: { [key: string]: number } = {
    realms: 1,
    search: 2,
    browse: 3,
    settings: 4,
  };
  const handleNav = (event: any, value: any) => {
    setIndex(value as number);
    if (value == 0) {
      redirect("/");
    }
    let path = Object.keys(pathval).find((key) => pathval[key] === value);
    redirect(`/${path}`);
  };
  useEffect(() => {
    let path = pathname.substring(pathname.lastIndexOf("/") + 1);

    console.log("PATH:", path);

    if (path == "") {
      setIndex(0);
      return;
    }
    if (pathval[path]) setIndex(pathval[path]);
  }, []);
  return (
    <Box
      sx={{
        flexGrow: 1,
        bottom: 1,
        position: "fixed",
        m: -3,
        p: 4,
        borderTopLeftRadius: "12px",
        borderTopRightRadius: "12px",
      }}
      className="toolbar"
    >
      <Tabs
        size="lg"
        aria-label="Bottom Navigation"
        value={index}
        onChange={handleNav}
        sx={(theme) => ({
          p: 1,
          borderRadius: 16,
          maxWidth: 400,
          mx: "auto",
          boxShadow: theme.shadow.sm,
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
          <Tooltip title="Home">
            <Tab
              disableIndicator
              orientation="vertical"
              {...(index === 0 && { color: color })}
            >
              <ListItemDecorator>
                <HomeRoundedIcon />
              </ListItemDecorator>
              
            </Tab>
          </Tooltip>
          <Tooltip title="Realms">
            <Tab
              disableIndicator
              orientation="vertical"
              {...(index === 1 && { color: color })}
            >
              <ListItemDecorator>
                <SatelliteAlt />
              </ListItemDecorator>
              
            </Tab>
          </Tooltip>
          <Tooltip title="Search">
            <Tab
              disableIndicator
              orientation="vertical"
              {...(index === 2 && { color: color })}
            >
              <ListItemDecorator>
                <Search />
              </ListItemDecorator>
              
            </Tab>
          </Tooltip>
          <Tooltip title="Browse">
            <Tab
              disableIndicator
              orientation="vertical"
              {...(index === 3 && { color: color })}
            >
              <ListItemDecorator>
                <CollectionsBookmarkIcon />
              </ListItemDecorator>
              
            </Tab>
          </Tooltip>
          <Tooltip title="Settings">
            <Tab
              disableIndicator
              orientation="vertical"
              {...(index === 4 && { color: color })}
            >
              <ListItemDecorator>
                <SettingsIcon />
              </ListItemDecorator>
            </Tab>
          </Tooltip>
        </TabList>
      </Tabs>
    </Box>
  );
};
export default ToolBar;
