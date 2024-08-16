"use client";
import { Center, Box, Grid, GridItem } from "@chakra-ui/react";
import { useState, useEffect } from "react";

import MarkdownRenderer from "./MarkdownRenderer";
import { ChatMessages, ChatInputForm } from "./ChatUIUtils";
import FileTable, { defaultLayout } from "./FileTable";

interface ChatAgent {
  role: boolean;
}

interface Message {
  role: string;
  content: string;
  key: symbol;
}

interface MessageComponentProps {
  message: Message;
  editMessage: () => {};
}

function ChatContainer({
  chatUrl,
  modelOptions,
  setCitations,
}: {
  chatUrl: string;
  modelOptions: string[];
  setCitations: (citations: any[]) => void;
}) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [needsResponse, setResponse] = useState(false);
  const [loadingResponse, setLoadingResponse] = useState(false);
  const [selectedModel, setSelectedModel] = useState("default");

  const getResponse = async () => {
    let chat_hist = messages.map((m) => {
      let { key, ...rest } = m;
      return rest;
    });

    const modelToSend = selectedModel === "default" ? undefined : selectedModel;
    setLoadingResponse(true);

    let result = await fetch(chatUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        accept: "application/json",
      },
      body: JSON.stringify({
        model: modelToSend,
        chat_history: chat_hist,
      }),
    })
      .then((resp) => {
        setLoadingResponse(false);
        if (resp.status < 200 || resp.status > 299) {
          return "failed request";
        }
        return resp.json();
      })
      .then((data) => {
        setCitations(data.citations || []);
        return data;
      })
      .catch((e) => {
        console.log("error making request");
        console.log(JSON.stringify(e));
      });

    setMessages([...messages, result.message]);
  };

  const sendMessage = async (message_content: string) => {
    let m: Message = {
      role: "user",
      content: message_content,
      key: Symbol(),
    };

    setMessages([...messages, m]);
    setResponse(true);
  };

  useEffect(() => {
    if (needsResponse) {
      getResponse();
      setResponse(false);
    }
  }, [needsResponse]);

  return (
    <>
      <ChatMessages messages={messages} loading={loadingResponse} />
      <ChatInputForm
        sendMessage={sendMessage}
        setSelectedModel={setSelectedModel}
        modelOptions={modelOptions}
      />
    </>
  );
}

function ChatUI({
  chatUrl,
  modelOptions,
}: {
  convoID?: string;
  chatUrl: string;
  modelOptions: string[];
}) {
  const [citations, setCitations] = useState<any[]>([]);

  return (
    <Center width="100%" height="100%">
      <Box
        border="solid"
        borderColor="oklch(92.83% 0.01 286.37)"
        borderRadius="10px"
        borderWidth="1px"
        width="95%"
        height="95vh"
        margin="20px"
        padding="20px"
        justifySelf="center"
        overflow="clip"
        position="relative"
      >
        <Grid h="100%" gridTemplateColumns="4fr 1fr" gap={5}>
          <GridItem
            rowSpan={10}
            colSpan="auto"
            overflow="scroll clip"
            position="relative"
          >
            <ChatContainer
              chatUrl={chatUrl}
              modelOptions={modelOptions}
              setCitations={setCitations}
            />
          </GridItem>
          <GridItem overflow="scroll clip" position="relative">
            <FileTable files={citations} layout={defaultLayout} />
          </GridItem>
        </Grid>
      </Box>
    </Center>
  );
}

export default ChatUI;
