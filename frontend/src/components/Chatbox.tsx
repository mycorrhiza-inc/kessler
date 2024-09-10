import { Stack } from "@mui/joy";
import { motion, PanInfo, useDragControls } from "framer-motion";
import { useEffect, MutableRefObject, useRef, RefObject } from "react";

import { Dispatch, SetStateAction, use, useState } from "react";
import { CloseIcon, HamburgerIcon } from "@/components/Icons";
import "./Chatbox.css";
import { set } from "lodash-es";

import { ChatMessages, exampleChatHistory } from "./ChatHistory";

interface ChatBoxProps {
  chatVisible: boolean;
  setChatVisible: Dispatch<SetStateAction<boolean>>;
  parentRef: RefObject<HTMLDivElement>;
}

const ChatBox = ({ chatVisible, setChatVisible, parentRef }: ChatBoxProps) => {
  const [chatSidebarVisible, setChatSidebarVisible] = useState(false);
  const [chatDisplayString, setChatDisplayString] = useState("none");
  const containerRef = useRef(null);
  const [isResizing, setIsResizing] = useState(false); // Track resizing state to trigger reset
  const [isDragging, setIsDragging] = useState(true); // Track dragging state to trigger reset
  const [position, setPosition] = useState({ x: 0, y: 0 });

  const controls = useDragControls();
  function startDrag(event: any) {
    controls.start(event);
  }
  useEffect(() => {
    if (!chatVisible) {
      setTimeout(() => setChatDisplayString("none"), 0.5);
    } else {
      setChatDisplayString("block");
    }
  });

  const [size, setSize] = useState({ width: 300, height: 300 });
  // const [position, setPosition] = useState({ x: 0, y: 0 });

  const minSize = 50;

  // Handle resizing from edges and corners

  // Render a resize handle
  const resizeHandle = (direction: string) => {
    const handleMouseMove = (e: MouseEvent) => {
      console.log("mouse moving");
      console.log("changing size");
      let newWidth = size.width;
      let newHeight = size.height;

      if (direction.includes("right")) {
        newWidth = Math.max(e.clientX - position.x, minSize);
      } else if (direction.includes("left")) {
        newWidth = Math.max(position.x - e.clientX + size.width, minSize);
        setPosition((prevPosition) => ({ ...prevPosition, x: e.clientX }));
      }

      if (direction.includes("bottom")) {
        newHeight = Math.max(e.clientY - position.y, minSize);
      } else if (direction.includes("top")) {
        newHeight = Math.max(position.y - e.clientY + size.height, minSize);
        setPosition((prevPosition) => ({ ...prevPosition, y: e.clientY }));
      }

      setSize({ width: newWidth, height: newHeight });
    };

    const handleMouseUp = () => {
      setIsResizing(false);
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);
    };

    // Start resizing process
    const startResize = (e: any) => {
      console.log("resizing");

      setIsResizing(true);
      console.log(isResizing);
      document.addEventListener("mousemove", handleMouseMove);
      document.addEventListener("mouseup", handleMouseUp);
    };
    return (
      <div className={`resize-handle ${direction}`} onMouseDown={startResize} />
    );
  };
  return (
    <motion.div
      animate={{
        minHeight: "40vh",
        // y: chatVisible ? "20vh" : "150vh", // if mobile, animate to the middle of the screen
        // x: "0", // if mobile, animate to the middle of the screen
        display: chatDisplayString,
        width: size.width,
        height: size.height,
        position: "absolute",
      }}
      drag={isDragging}
      dragConstraints={parentRef}
      dragControls={controls}
      dragMomentum={false}
      // dragElastic={0.01}
      // transition={{ type: "spring", stiffness: 100, damping: 20 }} // Customize the animation
      data-isOpen={!chatVisible}
      style={{
        top: position.y,
        left: position.x,
        backgroundColor: "white",
        borderRadius: "10px",
        border: "2px solid grey",
        padding: "10px",
        zIndex: 10000,
        color: "black",
        display: "none",
        width: size.width,
        height: size.height,
        position: "fixed",
        pointerEvents: "none", // Prevent the div from blocking pointer events of children
        overflow: "hidden", // Prevent children from overflowing the div
      }}
      ref={containerRef}
    >
      {resizeHandle("bottom-right")}
      {resizeHandle("bottom-left")}
      {resizeHandle("top-right")}
      {resizeHandle("top-left")}
      {resizeHandle("top")}
      {resizeHandle("bottom")}
      {resizeHandle("left")}
      {resizeHandle("right")}

      <div
        className="chatbox-banner"
        style={{
          position: "sticky",
          top: "0",
          backgroundColor: "#f1f1f1",
          padding: "20px",
          textAlign: "center",
          zIndex: "1000" /* Ensures it's on top of other content */,
          borderBottom: "1px solid #ccc",
          height: "auto",
          pointerEvents: "auto",
        }}
      >
        <Stack
          direction="row"
          justifyContent="space-between"
          onPointerDown={startDrag}
        >
          <button
            onClick={() => setChatSidebarVisible((prevState) => !prevState)}
          ></button>
        </Stack>
      </div>
      <div className="chatbox-banner sticky top-0 bg-[#f5f5f5] dark:bg-gray-700 p-5 text-center z-50 border-b border-gray-300 h-auto">
        <div className="flex flex-row justify-between">
          <button>
            <HamburgerIcon />
          </button>
          <button
            onClick={() => {
              setChatVisible((prev) => !prev);
            }}
          >
            <CloseIcon />
          </button>
        </div>
      </div>
      <div>
        <ChatMessages
          messages={exampleChatHistory}
          loading={false}
        ></ChatMessages>
      </div>

      <div
        className="chatContainer"
        style={{
          display: "flex",
          flexDirection: "row",
          height: "90%",
          width: "100%",
          padding: "2px",
        }}
      >
        <motion.div
          className="chatSidebar"
          animate={{
            width: chatSidebarVisible ? "20%" : "0%",
          }}
          style={{ backgroundColor: "red", overflow: "scroll" }}
        >
          sidebar contents
        </motion.div>
        <div> chat contents</div>
      </div>
    </motion.div>
  );
};

export default ChatBox;
