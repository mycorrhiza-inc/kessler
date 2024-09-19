import { Stack } from "@mui/joy";
import { motion, PanInfo, useDragControls } from "framer-motion";
import { useEffect, MutableRefObject, useRef, RefObject } from "react";
import { Dispatch, SetStateAction, useState } from "react";
import { CloseIcon, HamburgerIcon } from "@/components/Icons";
import "./ChatBoxAnimated.css";
import { ChatMessages, exampleChatHistory } from "./ChatHistory";
import ChatBoxInternals from "./ChatBoxInternals";

interface ChatBoxProps {
  chatVisible: boolean;
  setChatVisible: Dispatch<SetStateAction<boolean>>;
  parentRef: RefObject<HTMLDivElement>;
}

const ChatBox = ({ chatVisible, setChatVisible, parentRef }: ChatBoxProps) => {
  const [chatSidebarVisible, setChatSidebarVisible] = useState(false);
  const [chatDisplayString, setChatDisplayString] = useState("none");
  const containerRef = useRef(null);
  const [isResizing, setIsResizing] = useState(false);
  const [isDragging, setIsDragging] = useState(false);
  const [position, setPosition] = useState({ x: 0, y: 0 });

  const controls = useDragControls();
  function startDrag(event: any) {
    if (!isResizing) {
      // Only start dragging if not currently resizing
      setIsDragging(true);
      controls.start(event);
    }
  }
  useEffect(() => {
    if (!chatVisible) {
      setChatDisplayString("none");
    } else {
      setChatDisplayString("block");
    }
  }, [chatVisible]);

  const [size, setSize] = useState({ width: 300, height: 300 });
  const minSize = 50;

  // Handle resizing from edges and corners
  const resizeHandle = (direction: string) => {
    const handleMouseMove = (e: MouseEvent) => {
      let newWidth = size.width;
      let newHeight = size.height;

      if (direction.includes("right")) {
        newWidth = Math.max(e.clientX - position.x, minSize);
      } else if (direction.includes("left")) {
        newWidth = Math.max(position.x - e.clientX + size.width, minSize);
        setPosition({ ...position, x: e.clientX });
      }

      if (direction.includes("bottom")) {
        newHeight = Math.max(e.clientY - position.y, minSize);
      } else if (direction.includes("top")) {
        newHeight = Math.max(position.y - e.clientY + size.height, minSize);
        setPosition({ ...position, y: e.clientY });
      }

      setSize({ width: newWidth, height: newHeight });
    };

    const handleMouseUp = () => {
      setIsResizing(false);
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);
      setIsDragging(false); // Ensure dragging state is reset
    };

    const startResize = (e: React.MouseEvent) => {
      setIsResizing(true);
      setIsDragging(false);
      document.addEventListener("mousemove", handleMouseMove);
      document.addEventListener("mouseup", handleMouseUp);
    };

    return (
      <div className={`resize-handle ${direction}`} onMouseDown={startResize} />
    );
  };

  return (
    <motion.div
      className="overflow-y-scroll"
      animate={{
        minHeight: "40vh",
        display: chatDisplayString,
        width: size.width,
        height: size.height,
        position: "absolute",
      }}
      drag={isDragging}
      dragControls={controls}
      dragConstraints={parentRef}
      dragMomentum={false}
      data-isOpen={!chatVisible}
      style={{
        top: position.y,
        left: position.x,
        borderRadius: "10px",
        border: "2px solid grey",
        padding: "10px",
        zIndex: 10000,
        color: "black",
        display: "none",
        width: size.width,
        height: size.height,
        position: "fixed",
        pointerEvents: "none",
        overflow: "hidden",
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
      <ChatBoxInternals
      // chatSidebarVisible={chatSidebarVisible}
      // setChatSidebarVisible={setChatSidebarVisible}
      ></ChatBoxInternals>
    </motion.div>
  );
};
export default ChatBox;
