"use client";
import {
  Center,
  Circle,
  Button,
  IconButton,
  Flex,
  VStack,
  StackDivider,
  Box,
  Grid,
  GridItem,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Container,
} from "@chakra-ui/react";

import {
  Form,
  FormLayout,
  Field,
  DisplayIf,
  SubmitButton,
} from "@saas-ui/react";
import { FiArrowUpCircle } from "react-icons/fi";
import { useState } from "react";

function ChatElement(props: {}) {
  return <Box>a chat element</Box>;
}

function SourceModal() {
  // full screen modal for a given source
  //
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
      <br />
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
      tincidunt rhoncus. Vivamus id urna odio.
      <br />
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
      tincidunt rhoncus. Vivamus id urna odio.
      <br />
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
          <ModalBody>
            <SourceModal />
          </ModalBody>

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
  body: String;
}
interface SystemMessage {
  body: String;
}

interface MessageComponentProps {
  editMessage: () => {};
}

function MessageComponent(message: Message) {
  return <Box>{message.body}</Box>;
}

function ChatBox() {
  const [messages, setMessages] = useState<Message[]>([]);
  const appendMessage = (m: Message) => {
    setMessages([...messages, m]);
  };

  const sendMessage = (params: any) => {
    console.log("sending message");
    console.log(params);
    return new Promise((resolve) => {
      setTimeout(resolve, 1000);
    });
  };

  /*
		with take the current message, find the message

	*/

  return (
    <>
      <VStack
        divider={<StackDivider borderColor="gray.200" />}
        spacing={4}
        align="stretch"
        flexDir="column"
        overflow="scroll"
        h="100vh"
      >
        <SourceBox content="Source" />
        <SourceBox content="Source" />
        <SourceBox content="Source" />
        <SourceBox content="Source" />
        <SourceBox content="Source" />
        <SourceBox content="overflow" />
        <SourceBox content="overflow" />
        <SourceBox content="overflow" />
        <SourceBox content="overflow" />
        <SourceBox content="overflow" />
        <SourceBox content="overflow" />
      </VStack>
      <Form onSubmit={sendMessage}>
        <FormLayout>
          <Box
            display="flex"
            flexDir="row"
            justifySelf="center"
            justifyContent="center"
            bg="white"
            w="90%"
            borderColor="green"
            border="solid"
            borderRadius="10px"
            borderWidth="1px"
            position="absolute"
            bottom="0"
          >
            <Field
              name="messageInput"
              type="textarea"
              placeholder="chat..."
              // flex-grow="2"
              paddingLeft="20px"
              resize="none"
              border="none"
              padding="10px"
              margin="10px"
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
      </Form>
    </>
  );
}
/*
 */
export default function ChatUI({ convoID = "" }: { convoID?: string }) {
  // convoId being empty is a new chat instance
  function appendNewMessage() {}

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
        <Grid h="100%" gridTemplateColumns={"5fr 1.4fr"} gap={5}>
          <GridItem rowSpan={10} overflow="scroll clip" position="relative">
            <ChatBox />
          </GridItem>

          <GridItem rowSpan={10} overflow="scroll clip">
            <ContextSources />
          </GridItem>
        </Grid>
      </Box>
    </Center>
  );
}
