import { Stack } from "@mui/joy";
import { useEffect, MutableRefObject, useRef, RefObject } from "react";
import { Dispatch, SetStateAction, useState } from "react";
import { CloseIcon, HamburgerIcon } from "@/components/Icons";
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
  const [position, setPosition] = useState({ x: 0, y: 0 });

  useEffect(() => {
    if (!chatVisible) {
      setChatDisplayString("none");
    } else {
      setChatDisplayString("block");
    }
  }, [chatVisible]);

  const [size, setSize] = useState({ width: 300, height: 300 });
  const minSize = 50;

  return (
    <div
      style={{
        minHeight: "40vh",
        display: chatDisplayString,
        width: size.width,
        height: size.height,
        position: "absolute",
        top: position.y,
        left: position.x,
        backgroundColor: "white",
        borderRadius: "10px",
        border: "2px solid grey",
        padding: "10px",
        zIndex: 10000,
        color: "black",
        pointerEvents: "none",
        overflow: "hidden",
      }}
      ref={containerRef}
    >
      <ChatBoxInternals
        chatSidebarVisible={chatSidebarVisible}
        setChatSidebarVisible={setChatSidebarVisible}
      ></ChatBoxInternals>
    </div>
  );
};

export default ChatBox;
