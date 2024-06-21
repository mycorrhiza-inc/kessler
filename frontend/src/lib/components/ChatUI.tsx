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
import Markdown from 'react-markdown'
import { message } from "antd";
import { start } from "repl";
import { initialState } from "node_modules/@clerk/nextjs/dist/types/app-router/server/auth";


import MarkdownRenderer from "./MarkdownRenderer"
interface ChatAgent {
  role: boolean;
}

function SourceModal() {
  // full screen modal for a given source

  // TODO: cache the sourceModal text in the zustand state manager
  return (
    <Box>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur
      ultricies vehicula velit, at condimentum diam tristique a. Sed commodo,
      metus quis scelerisque porta, turpis libero placerat erat, quis semper
      sapien nisi eu neque. Donec id maximus lacus. Proin dolor erat, tempus ac
      scelerisque fringilla, imperdiet in orci. Curabitur egestas magna ut
      mollis sollicitudin. Sed nec pulvinar eros. Donec porta tempor convallis.
      Aliquam erat volutpat. Nulla facilisi. Aenean faucibus ipsum sit amet
      dictum lobortis. Cras congue magna sapien, et facilisis arcu efficitur
      vitae. Integer nec tellus nec lectus molestie tristique. Vestibulum a sem
      aliquam, cursus erat ac, luctus lectus. Pellentesque et augue facilisis,
      ullamcorper ante at, facilisis nibh. Mauris placerat ut sapien ac pretium.
      Fusce maximus dignissim diam quis blandit. Aliquam eget ultrices velit.
      Cras in rutrum arcu. Praesent volutpat, sapien vitae pulvinar commodo,
      urna neque lobortis arcu, a condimentum libero arcu sagittis leo. Nulla
      ante arcu, pharetra id turpis quis, consectetur semper ex. Pellentesque
      habitant morbi tristique senectus et netus et malesuada fames ac turpis
      egestas. Praesent pulvinar id turpis sed semper. Praesent ac dignissim
      odio. Nullam accumsan tincidunt augue. Nullam id interdum metus. Vivamus
      sodales laoreet dolor a condimentum. Maecenas tincidunt semper lorem, sed
      lacinia tellus lacinia sit amet. Ut augue odio, fermentum ut luctus at,
      cursus quis risus. Vivamus lacinia augue lectus. Morbi facilisis nibh
      massa, in vestibulum mauris fermentum at. Lorem ipsum dolor sit amet,
      consectetur adipiscing elit. Ut sollicitudin fermentum justo, vitae
      malesuada lorem ullamcorper id. Mauris vitae dui ornare, luctus libero
      nec, posuere massa. Suspendisse at odio tristique, maximus nunc eu,
      ullamcorper mauris. Mauris ultricies sit amet tortor non gravida. Cras
      pellentesque justo dui, posuere ultricies ante lobortis quis. Sed bibendum
      mattis fringilla. Quisque id pharetra ligula. In consectetur sagittis enim
      sed pellentesque. Vestibulum quam libero, posuere et nisi ut, egestas
      laoreet augue. Aliquam diam ex, ultricies quis pellentesque ut, venenatis
      et erat. Nam at elementum augue, sit amet fermentum lorem. Cras aliquet
      elit eu dui mollis gravida vitae non sapien. Donec sit amet est eu est
      tincidunt rhoncus. Vivamus id urna odio. Lorem ipsum dolor sit amet,
      consectetur adipiscing elit. Curabitur ultricies vehicula velit, at
      condimentum diam tristique a. Sed commodo, metus quis scelerisque porta,
      turpis libero placerat erat, quis semper sapien nisi eu neque. Donec id
      maximus lacus. Proin dolor erat, tempus ac scelerisque fringilla,
      imperdiet in orci. Curabitur egestas magna ut mollis sollicitudin. Sed nec
      pulvinar eros. Donec porta tempor convallis. Aliquam erat volutpat. Nulla
      facilisi. Aenean faucibus ipsum sit amet dictum lobortis. Cras congue
      magna sapien, et facilisis arcu efficitur vitae. Integer nec tellus nec
      lectus molestie tristique. Vestibulum a sem aliquam, cursus erat ac,
      luctus lectus. Pellentesque et augue facilisis, ullamcorper ante at,
      facilisis nibh. Mauris placerat ut sapien ac pretium. Fusce maximus
      dignissim diam quis blandit. Aliquam eget ultrices velit. Cras in rutrum
      arcu. Praesent volutpat, sapien vitae pulvinar commodo, urna neque
      lobortis arcu, a condimentum libero arcu sagittis leo. Nulla ante arcu,
      pharetra id turpis quis, consectetur semper ex. Pellentesque habitant
      morbi tristique senectus et netus et malesuada fames ac turpis egestas.
      Praesent pulvinar id turpis sed semper. Praesent ac dignissim odio. Nullam
      accumsan tincidunt augue. Nullam id interdum metus. Vivamus sodales
      laoreet dolor a condimentum. Maecenas tincidunt semper lorem, sed lacinia
      tellus lacinia sit amet. Ut augue odio, fermentum ut luctus at, cursus
      quis risus. Vivamus lacinia augue lectus. Morbi facilisis nibh massa, in
      vestibulum mauris fermentum at. Lorem ipsum dolor sit amet, consectetur
      adipiscing elit. Ut sollicitudin fermentum justo, vitae malesuada lorem
      ullamcorper id. Mauris vitae dui ornare, luctus libero nec, posuere massa.
      Suspendisse at odio tristique, maximus nunc eu, ullamcorper mauris. Mauris
      ultricies sit amet tortor non gravida. Cras pellentesque justo dui,
      posuere ultricies ante lobortis quis. Sed bibendum mattis fringilla.
      Quisque id pharetra ligula. In consectetur sagittis enim sed pellentesque.
      Vestibulum quam libero, posuere et nisi ut, egestas laoreet augue. Aliquam
      diam ex, ultricies quis pellentesque ut, venenatis et erat. Nam at
      elementum augue, sit amet fermentum lorem. Cras aliquet elit eu dui mollis
      gravida vitae non sapien. Donec sit amet est eu est tincidunt rhoncus.
      Vivamus id urna odio.
    </Box>
  );
}

function SourceBox({ content }: { content: string }) {
  // on click show modal of the document
  const [openModal, changeModal] = useState(false);
  const toggleModal = () => {
    changeModal(!openModal);
  };
  return (
    <>
      <Modal isOpen={openModal} onClose={toggleModal}>
        <ModalOverlay />
        <ModalContent maxH="3000px" maxW="1500px" overflow="scroll">
          <ModalHeader>Modal Title</ModalHeader>
          <ModalCloseButton />
          <ModalContent>
            <SourceModal />
          </ModalContent>

          <ModalFooter>
            <Button colorScheme="blue" mr={3} onClick={toggleModal}>
              Close
            </Button>
            <Button variant="ghost">Secondary Action</Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
      <Box
        width="95%"
        margin="5px"
        border="solid"
        borderColor="oklch(92.83% 0.01 286.37)"
        borderRadius="10px"
        borderWidth="1px"
        height="100px"
        padding="10px"
        onClick={toggleModal}
      >
        {content}
      </Box>
    </>
  );
}

function ContextSources() {
  // TODO: generate sources from some list of soruce objects
  return (
    <VStack
      divider={<StackDivider borderColor="gray.200" />}
      spacing={4}
      align="stretch"
      overflowY="scroll"
      h="100vh"
    >
      <SourceBox content="Source" />
      <SourceBox content="Source" />
      <SourceBox content="Source" />
      <SourceBox content="Source" />
      <SourceBox content="Source" />
      <SourceBox content="Source" />
      <SourceBox content="Overflow" />
      <SourceBox content="Overflow" />
      <SourceBox content="Overflow" />
      <SourceBox content="Overflow" />
      <SourceBox content="Overflow" />
      <SourceBox content="Overflow" />
    </VStack>
  );
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

function MessageComponent({
  message = {
    role: "user",
    content: "",
    key: Symbol(),
  },
}: {
  message: Message;
}) {
  return (
    <Box
      width="90%"
      background={message.role == "user" ? "aquamarine" : "antiquewhite"}
      borderRadius="10px"
      // maxWidth="800px"
      height="auto"
      overflow="auto"
      minHeight="100px"
      justifyContent={message.role ? "right" : "left"}
      padding="20px"
    >
      {/* enclosing the message */}
      <VStack
        // divider={<StackDivider borderColor="gray.200" />}
        spacing={4}
        // align="stretch"
        overflowY="scroll"
      // justifyContent={message.role == true ? "right" : "left"}
      // h="100vh"
      >
        {/* role message */}
        <div><MarkdownRenderer>{message.content}</MarkdownRenderer></div>
        {/* <Box width="100%" height="50px">
          {!message.role && <div>Regenerate</div>}{" "}
          {message.role && <div>Edit</div>}
        </Box> */}
      </VStack>
    </Box>
  );
}
function ChatBox({ chatUrl }: { chatUrl: string }) {
  // const [messages, setMessages] = useState<Message[]>(startingMessages);
  const [messages, setMessages] = useState<Message[]>([]);
  const [needsResponse, setResponse] = useState(false);
  const [userChatbox, setUserChatbox] = useState("");
  // let messages: Message[] = [];

  const getResponse = async () => {
    let chat_hist = messages.map((m) => {
      let { key, ...rest } = m;
      return rest;
    });
    console.log(`chat history`);
    console.log(chat_hist);
    let result = await fetch(
      // FIXME : Add the base url instead of localhost to make it more amenable to this stuff.
      // "http://localhost/api/rag/rag_chat",
      chatUrl,
      {
        method: "POST",
        mode: "cors",
        headers: {
          "Content-Type": "application/json",
          accept: "application/json",
          // "Access-Control-Allow-Origin": "*",
          "Referrer-Policy": "no-referrer",
        },
        body: JSON.stringify({
          model: null,
          chat_history: chat_hist,
        }),
      },
    )
      .then((resp) => {
        console.log("completed request");
        console.log(resp);
        if (resp.status < 200 || resp.status > 299) {
          console.log(`error with request:\n${resp}`);
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
    console.log(result);
    setMessages([...messages, result]);
    // let c = result.body;

    // setMessages(result);
  };
  useEffect(() => {
    if (needsResponse) getResponse(); // This is be executed when `loading` state changes
  }, [needsResponse]);

  interface msgSent {
    messageInput: string;
  }
  // TODO : Fix Horrible Buggy passing the any function
  const sendMessage = async (params: msgSent) => {
    console.log("sending message");
    console.log(`msg: ${params.messageInput}`);
    let m: Message = {
      role: "user",
      content: params.messageInput,
      key: Symbol(),
    };
    console.log(`appending message "${m.content}"`);
    if (messages.length == 0) {
      // messages = [m];
      setMessages([m]);
      console.log(messages);
    } else {
      // messages = [...messages, m];
      setMessages([...messages, m]);
      console.log(messages);
    }
    setResponse(true);
    setUserChatbox("");
  };

  /*
    with take the current message, find the message
 
  */

  // get a response every time a message is sent
  useEffect(() => {
    console.log("getting response");
    if (needsResponse) {
      getResponse();
      setResponse(false);
    }
  }, [needsResponse]);
  return (
    <>
      <VStack
        divider={<StackDivider borderColor="gray.200" />}
        spacing={4}
        flexDir="column"
        overflow="scroll"
        h="100vh"
      >
        {messages.length === 0 && (
          <Box p={5} textAlign="center" color="gray.500">
            <Text fontSize="lg" fontWeight="bold">Welcome to the Chatbot!</Text>
            <Text>Type your message in the input box below and press Enter to send.</Text>
          </Box>
        )}
        {messages.map((m: Message) => {
          return <MessageComponent
            // key={m.key.toString()} 
            message={m}
          />;
        })}
        <Box minHeight="300px" width="100%" color="red" />
      </VStack>
      <Form onSubmit={sendMessage}>
        <FormLayout>
          <Box
            display="flex"
            flexDir="row"
            justifySelf="center"
            justifyContent="center"
            bg="white"
            w="85%"
            borderColor="green"
            borderRadius="10px"
            borderWidth="1px"
            position="absolute"
            bottom="0"
          >
            <Field
              name="messageInput"
              type="text"
              placeholder="chat..."
              // flex-grow="2"
              paddingLeft="20px"
              resize="none"
              border="none"
              padding="10px"
              margin="10px"
            // onKeyPress={(event) => {
            //   if (event.key === 'Enter') {
            //     event.preventDefault();
            //     this.myFormRef.requestSubmit();
            //   }
            // }}
            // value={userChatbox}
            // // FIXME : Figure out the proper type for this
            // onChange={(e: any) => setUserChatbox(e.targetvalue)}
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
          </Box>
        </FormLayout>
      </Form >
    </>
  );
}
/*
 */
function ChatUI({ convoID = "", chatUrl }: { convoID?: string, chatUrl: string }) {
  // convoId being empty is a new chat instance

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
            <ChatBox chatUrl={chatUrl} />
          </GridItem >

          {/* <GridItem rowSpan={10} overflow="scroll clip">
            <ContextSources />
          </GridItem> */}
        </Grid>
      </Box>
    </Center>
  );
}

export default ChatUI;
