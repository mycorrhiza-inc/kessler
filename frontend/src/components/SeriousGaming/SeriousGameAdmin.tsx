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

const initialNodes = [
  { id: "1", position: { x: 0, y: 0 }, data: { label: "1" } },
  { id: "2", position: { x: 0, y: 100 }, data: { label: "2" } },
];
const initialEdges = [{ id: "e1-2", source: "1", target: "2" }];

export default function TestFlowVisuals({ user }: { user: User | null }) {
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
      <div style={{ width: "50vw", height: "50vh" }}>
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
}
