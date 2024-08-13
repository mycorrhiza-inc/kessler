"use client";
import {
  Center,
  Circle,
  Button,
  IconButton,
  Flex,
  VStack,
  HStack,
  StackDivider,
  Box,
  Grid,
  GridItem,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalCloseButton,
  Container,
  Text,
  Select,
  SkeletonText,
  DarkMode,
} from "@chakra-ui/react";
import {
  Form,
  FormLayout,
  Field,
  DisplayIf,
  SubmitButton,
} from "@saas-ui/react";
import { FiArrowUpCircle } from "react-icons/fi";
import { useState, useEffect } from "react";
import { initialState } from "node_modules/@clerk/nextjs/dist/types/app-router/server/auth";
import { useColorMode, useColorModeValue } from "@chakra-ui/react";

import MarkdownRenderer from "./MarkdownRenderer";
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

function ChatInputForm({
  sendMessage,
  setSelectedModel,
  modelOptions,
}: {
  sendMessage: (messageContent: string) => Promise<void>;
  setSelectedModel: (model: string) => void;
  modelOptions: string[];
}) {
  const [userChatbox, setUserChatbox] = useState("");

  const handleMessageSubmit = async (messageContent: string) => {
    await sendMessage(messageContent);
    setUserChatbox("");
  };

  return (
    <Form
      onSubmit={(params: { messageInput: string }) =>
        handleMessageSubmit(params.messageInput)
      }
    >
      <FormLayout>
        <Box
          display="flex"
          flexDir="row"
          justifySelf="center"
          justifyContent="center"
          w="85%"
          borderRadius="10px"
          borderWidth="1px"
          position="absolute"
          bottom="0"
        >
          <Field
            name="messageInput"
            type="textarea"
            placeholder="chat..."
            paddingLeft="20px"
            resize="none"
            border="none"
            padding="10px"
            margin="10px"
            onKeyPress={(event) => {
              if (event.key === "Enter" && !event.shiftKey) {
                event.preventDefault();
                handleMessageSubmit(userChatbox);
              }
            }}
            value={userChatbox}
            onChange={(e: any) => setUserChatbox(e.target.value)}
          />
          <Center padding="10px">
            <IconButton
              isRound={true}
              variant="solid"
              colorScheme="green"
              aria-label="Send"
              type="submit"
              icon={<FiArrowUpCircle />}
            />
          </Center>
          <Select
            placeholder="default"
            onChange={(e) => setSelectedModel(e.target.value)}
            margin="10px"
            size="sm"
            width="150px"
          >
            {modelOptions.map((model: string) => (
              <option key={model} value={model}>
                {model}
              </option>
            ))}
          </Select>
        </Box>
      </FormLayout>
    </Form>
  );
}

interface Message {
  role: string;
  content: string;
  key: symbol;
}

function MessageComponent({ message }: { message: Message }) {
  const { colorMode } = useColorMode();
  const isUser = message.role === "user";

  return (
    <HStack width="100%" justifyContent={isUser ? "flex-end" : "flex-start"}>
      <Box
        width="90%"
        background={
          isUser
            ? colorMode === "light"
              ? "teal.100"
              : "teal.700"
            : colorMode === "light"
              ? "gray.200"
              : "gray.700"
        }
        borderRadius="10px"
        overflow="auto"
        minHeight="100px"
        padding="20px"
      >
        <MarkdownRenderer>{message.content}</MarkdownRenderer>
      </Box>
    </HStack>
  );
}

function AwaitingMessageSkeleton() {
  const { colorMode } = useColorMode();

  return (
    <Box
      width="90%"
      background={colorMode === "light" ? "gray.200" : "gray.700"}
      borderRadius="10px"
      minHeight="100px"
      padding="20px"
    >
      <SkeletonText
        startColor="pink.500"
        endColor="orange.500"
        mt="4"
        noOfLines={3}
        spacing="4"
        skeletonHeight="2"
      />
    </Box>
  );
}

function ChatMessages({
  messages,
  loading,
}: {
  messages: Message[];
  loading: boolean;
}) {
  return (
    <VStack
      divider={<StackDivider borderColor="gray.200" />}
      spacing={4}
      align="stretch"
      p={4}
      borderRadius="md"
      overflowY="auto"
      flexDir="column"
      h="100vh"
    >
      {messages.length === 0 && (
        <Box p={5} textAlign="center" color="gray.500">
          <Text fontSize="lg" fontWeight="bold">
            Welcome to the Chatbot!
          </Text>
          <Text>
            Type your message in the input box below and press Enter to send.
          </Text>
        </Box>
      )}
      {messages.map((m: Message) => (
        <MessageComponent key={m.key.toString()} message={m} />
      ))}
      {loading && <AwaitingMessageSkeleton />}
      <Box minHeight="300px" width="100%" color="red" />
    </VStack>
  );
}

function ChatContainer({
  chatUrl,
  modelOptions,
}: {
  chatUrl: string;
  modelOptions: string[];
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
        return data;
      })
      .catch((e) => {
        console.log("error making request");
        console.log(JSON.stringify(e));
      });

    setMessages([...messages, result]);
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
  convoID = "",
  chatUrl,
  modelOptions,
}: {
  convoID?: string;
  chatUrl: string;
  modelOptions: string[];
}) {
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
        <Grid h="100%" gridTemplateColumns={"5fr"} gap={5}>
          <GridItem
            rowSpan={10}
            colSpan="auto"
            overflow="scroll clip"
            position="relative"
          >
            <ChatContainer chatUrl={chatUrl} modelOptions={modelOptions} />
          </GridItem>
        </Grid>
      </Box>
    </Center>
  );
}

export default ChatUI;
