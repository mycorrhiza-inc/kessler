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

export function ChatInputForm({
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

export function ChatMessages({
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
        <MessageComponent message={m} />
      ))}
      {loading && <AwaitingMessageSkeleton />}
      <Box minHeight="300px" width="100%" color="red" />
    </VStack>
  );
}
