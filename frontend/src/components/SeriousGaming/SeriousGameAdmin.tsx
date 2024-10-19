"use client";
import React, { useCallback } from "react";
import {
  ReactFlow,
  MiniMap,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  addEdge,
} from "@xyflow/react";

import "@xyflow/react/dist/style.css";
import { User } from "@supabase/supabase-js";
import Header from "../Header";
import { useTheme } from "next-themes";
import { themeDataDictionary } from "../ThemeSelector";

interface SeriousGame {
  teams: SeriousGameTeam[];
  turns: SeriousGameTurn[];
  description: string;
  completed: boolean;
  post_game: PostGameInfo | null; // Give this its own type at some point.
}
interface Organization {}
interface PostGameInfo {}

interface SeriousGameTeam {
  members: Organization[];
  description: string;
  objective: string;
}
interface SeriousGameTurn {
  completed: boolean;
  context: string;
  actions: SeriousGameAction[];
  synthesis: string | null;
}
interface SeriousGameAction {
  team_uuid: string;
  context: string;
  suggested_action: string;
  action_probabilities: ActionProbability[];
  chosen_probability: number;
  description: string;
}
interface ActionProbability {
  probability: number; // float between 0 and 1
  definitive_action_description: string;
}

const exampleTurn: SeriousGameTurn = {
  completed: true,
  context: "The local fantasy town council meeting is underway.",
  actions: [
    {
      team_uuid: "wizard-team",
      context: "Wizard tries to sway town council with a vision of prosperity.",
      suggested_action: "Cast an illusion spell of a bountiful future.",
      action_probabilities: [
        {
          probability: 0.6,
          definitive_action_description:
            "The council is intrigued but skeptical.",
        },
      ],
      chosen_probability: 0.4,
      description: "The wizard's illusion was partially effective.",
    },
    {
      team_uuid: "sorcerer-team",
      context: "Sorcerer aims to manipulate emotions to incite support.",
      suggested_action: "Charm council members with emotional manipulation.",
      action_probabilities: [
        {
          probability: 0.7,
          definitive_action_description:
            "Council members are swayed by the emotional appeal.",
        },
      ],
      chosen_probability: 0.5,
      description: "The sorcerer's charm spell had a significant impact.",
    },
  ],
  synthesis:
    "The council favors neither party completely but is more open to magic in political affairs.",
};

type NodeType = {
  id: string;
  position: { x: number; y: number };
  data: { label: string };
};

type EdgeType = {
  id: string;
  source: string;
  target: string;
};

// 2. A way to render the output of a SeriousGameTurn as a network of nodes and edges, the context should be on the leftmost side with the actions for each user being stacked vertically on top (the nodes should start with the context, continue to the action, then the described outcomes, then the probability roll and the final outcome selection.) they should all combine on the rightmost side for the final description of what happened in said turn
const createNodesAndEdgesFromTurn = (
  turn: SeriousGameTurn,
): { nodes: NodeType[]; edges: EdgeType[] } => {
  let nodes: NodeType[] = [];
  let edges: EdgeType[] = [];

  // Root context node
  nodes.push({
    id: "context",
    position: { x: 0, y: 0 },
    data: { label: turn.context },
  });

  let currentYPosition = 100;

  turn.actions.forEach((action, index) => {
    const actionNodeId = `action-${index}`;
    // Add Action Node
    nodes.push({
      id: actionNodeId,
      position: { x: 100, y: currentYPosition },
      data: { label: action.suggested_action },
    });

    edges.push({
      id: `context-to-action-${index}`,
      source: "context",
      target: actionNodeId,
    });

    const descriptionNodeId = `description-${index}`;
    // Add Description Node
    nodes.push({
      id: descriptionNodeId,
      position: { x: 200, y: currentYPosition },
      data: { label: action.description },
    });

    edges.push({
      id: `action-to-description-${index}`,
      source: actionNodeId,
      target: descriptionNodeId,
    });

    action.action_probabilities.forEach((probability, probIndex) => {
      const probNodeId = `probability-${index}-${probIndex}`;
      // Add Probability Node
      nodes.push({
        id: probNodeId,
        position: { x: 300, y: currentYPosition },
        data: {
          label: `Prob: ${probability.probability}, Chosen: ${action.chosen_probability}`,
        },
      });

      edges.push({
        id: `description-to-probability-${index}-${probIndex}`,
        source: descriptionNodeId,
        target: probNodeId,
      });

      const outcomeNodeId = `outcome-${index}-${probIndex}`;
      // Add Outcome Node
      nodes.push({
        id: outcomeNodeId,
        position: { x: 400, y: currentYPosition },
        data: { label: probability.definitive_action_description },
      });

      edges.push({
        id: `probability-to-outcome-${index}-${probIndex}`,
        source: probNodeId,
        target: outcomeNodeId,
      });

      currentYPosition += 100;
    });
  });

  const synthesisNodeId = "synthesis";
  // Add Synthesis Node
  if (turn.synthesis) {
    nodes.push({
      id: synthesisNodeId,
      position: { x: 500, y: currentYPosition / 2 },
      data: { label: turn.synthesis },
    });

    turn.actions.forEach((_, index) => {
      edges.push({
        id: `outcome-to-synthesis-${index}`,
        source: `outcome-${index}-0`, // Assume there's only one outcome now, adjust if multiple
        target: synthesisNodeId,
      });
    });
  }

  return { nodes, edges };
};
const initialNodes = [
  { id: "1", position: { x: 0, y: 0 }, data: { label: "1" } },
  { id: "2", position: { x: 0, y: 100 }, data: { label: "2" } },
];
const initialEdges = [{ id: "e1-2", source: "1", target: "2" }];

const TestFlowVisuals = ({ user }: { user: User | null }) => {
  const { theme } = useTheme();
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);
  // Fix Later
  // @ts-ignore
  const themeLightDark = themeDataDictionary[theme].lightdark || "dark";
  const onConnect = useCallback(
    (params: any) => setEdges((eds) => addEdge(params, eds)),
    [setEdges],
  );

  return (
    <div>
      <Header user={user} />
      <div style={{ width: "90vw", height: "80vh" }}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          colorMode={themeLightDark}
        >
          <Controls />
          <MiniMap />
          {/* Did react forget how enums work? */}
          {/* @ts-ignore */}
          <Background variant="dots" gap={12} size={1} />
        </ReactFlow>
      </div>
    </div>
  );
};
export default TestFlowVisuals;
