import {
  Node,
  Edge,
  applyEdgeChanges,
  applyNodeChanges,
  NodeChange,
  EdgeChange,
  OnNodesChange,
  OnEdgesChange,
  OnConnect,
  Viewport,
  FitViewOptions,
  useReactFlow,
  FitView,
} from "reactflow";
// const { fitView } = useReactFlow();

import * as d3 from "d3";

import { createWithEqualityFn } from "zustand/traditional";
import { persist, createJSONStorage } from "zustand/middleware";

export type GraphState = {
  nodes: Node[];
  edges: Edge[];
  onNodesChange: OnNodesChange;
  onEdgesChange: OnEdgesChange;
  onConnect: OnConnect;
  viewport: Viewport;
  width: number;
  height: number;
  fv: FitView | null;
  setFitView: (fv: FitView) => void;
  nodeTypes: any;
  setNodes: (nodes: Node[]) => void;
  setEdges: (edges: Edge[]) => void;
  runSimulation: (x: number, y: number, nodes: Node[]) => void;
  setNodeData: (nodeID: string, nodeData: any) => void;
  generateRandomNodes: (n: number) => void;
  generateDebugNodeGraph: () => void;
  ClusterAroundNode: (nodeId: string) => void;
  focusedNode: Node | null;
};

const UseGraphStore = createWithEqualityFn<
  GraphState,
  [["zustand/persist", GraphState]]
>(
  persist(
    (set, get) => ({
      focusedNode: null,
      focusedSubflow: null, // a node with type : 'group'
      nodes: [],
      edges: [],
      fv: null,
      setFitView: (fv: FitView) => {
        set({ fv });
      },
      viewport: {
        x: 0,
        y: 0,
        zoom: 0.5,
      },
      width: 1000,
      height: 1000,
      // need to be able to run the sim on any set of nodes
      // NOTE: Should this be defined elsewhere and imported here, not sure about how the ideal structure of this should look? - Nic
      // Disabled for now to prevent some error that I think was caused by recursive imports.
      nodeTypes: {},
      onNodesChange: (changes: NodeChange[]) => {
        set({ nodes: applyNodeChanges(changes, get().nodes) });
      },
      onEdgesChange: (changes: EdgeChange[]) => {
        set({ edges: applyEdgeChanges(changes, get().edges) });
      },
      onConnect: () => {},
      setNodes: (nodes) => set({ nodes }),
      setEdges: (edges) => set({ edges }),
      generateRandomNodes: (n: number) => {
        // logic to generate random nodes
        let nodes: Node[] = [];
        const newnode = (name: number) => {
          return {
            id: `${name}`,
            data: {
              label: "Random Node: " + name,
            },
            position: {
              x: Math.round(Math.random() * 1000),
              y: Math.round(Math.random() * 1000),
            },
          };
        };

        for (let i = 0; i < n; i++) {
          nodes.push(newnode(i));
        }
        set({ nodes });
      },
      generateDebugNodeGraph: () => {
        // logic to generate random nodes
        let nodes: Node[] = [
          {
            id: "init-chat-node",
            type: "ChatNode",
            data: {
              chat_history: [
                {
                  role: "system",
                  content: "Behave as an autopilot for a spacecraft.",
                },
                {
                  role: "human",
                  content: "Open the pod bay doors HAL.",
                },
                {
                  role: "ai",
                  content: "I am sorry dave. I cannot do that.",
                },
              ],
              model_name: "Qwen1.5",
            },
            position: { x: 100, y: 100 },
          },
          {
            id: "init-document-node",
            type: "DocumentNode",
            data: {
              docid: {
                metadata: { title: "House of Leaves" },
                extras: { short_summary: "An extremely weird book" },
              },
              document_text:
                "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem. Ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur? Quis autem vel eum iure reprehenderit qui in ea voluptate velit esse quam nihil molestiae consequatur, vel illum qui dolorem eum fugiat quo voluptas nulla pariatur?",
            },
            position: { x: 500, y: 100 },
          },
        ];
        set({ nodes });
      },
      // TODO: This seems really inneficent potentially, this might be a misunderstanding of my understanding of react/reactflow, and how they internally store nodes. Or it could be a limitation of using a list to store nodes when something like a hashmap would be better.
      setNodeData: (nodeID: string, data_dict: any) => {
        set({
          nodes: get().nodes.map((node) => {
            if (node.id === nodeID) {
              // it's important to create a new object here, to inform React Flow about the changes
              return { ...node, data: data_dict };
            }
            return node;
          }),
        });
      },
      runSimulation: (centerX: number, centerY: number, nodes: Node[]) => {
        let edges = get().edges;
        let simNodes = nodes.map((node) => ({
          node: node,
          x: node.position.x,
          y: node.position.y,
        }));
        console.log(edges);
        const simulation = d3
          .forceSimulation(simNodes)
          .force("charge", d3.forceManyBody().strength(0))
          .force("center", d3.forceCenter(centerX, centerY))
          .force("x", d3.forceX())
          .force("y", d3.forceY())
          .force("collide", d3.forceCollide(70))
          // .force("link", d3.forceLink(edges).strength(0.05).distance(100))
          .stop();
        simulation.on("tick", () => {
          simNodes.forEach((node) => {
            node.node.position = {
              x: node.x,
              y: node.y,
            };
          });
          nodes = simNodes.reduce((nodes, node) => {
            const newNode = {
              id: node.node.id,
              data: node.node.data,
              position: {
                x: node.x,
                y: node.y,
              },
            };
            return [...nodes, newNode];
          }, [] as Node[]);
          console.log(edges);
          set({ nodes });
          get();
        });

        simulation.alpha(0.3).restart();
        setTimeout(() => {
          set({ nodes });
          get().fv?.apply({});
        }, 3000);
        // get().fitView();
      },
      ClusterAroundNode: (nodeId: string) => {
        let nodes = get().nodes;
        const node = nodes.find((n) => n.id === nodeId);
        if (node) {
          get().runSimulation(node.position.x, node.position.y, get().nodes);
        }
      },

      CenterAroundNode: (nodeId: string) => {
        // logic to center view around node
        const node = get().nodes.find((n) => n.id === nodeId);

        if (node) {
          const viewport = {
            x: node.position.x,
            y: node.position.y,
            zoom: 1,
          };

          // update viewport
          set({ viewport });
        }
      },
    }),
    {
      name: "kessler graph store",
      storage: createJSONStorage(() => localStorage),
    },
  ),
);
// export default GraphClass;
export default UseGraphStore;
