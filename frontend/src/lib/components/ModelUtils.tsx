import React, { useState, useEffect } from "react";
import {
  Button,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Grid,
  GridItem,
  useDisclosure,
  useToast,
} from "@chakra-ui/react";

import { Table, Thead, Tbody, Tr, Th, Td } from "@chakra-ui/react";

const ViewIframeModal: React.FC<{ iframeURL: string; buttonName: string }> = ({
  iframeURL,
  buttonName,
}) => {
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);

  const handleOpenModal = () => {
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
  };

  return (
    <>
      <Button onClick={handleOpenModal}>{buttonName}</Button>
      {isModalOpen && (
        <Modal isOpen={isModalOpen} onClose={handleCloseModal} size="6xl">
          <ModalOverlay />
          <ModalContent maxW="80%">
            <ModalHeader></ModalHeader>
            <ModalCloseButton />
            <ModalBody>
              <iframe
                title="Iframe Content"
                src={iframeURL || "https://en.wikipedia.org/wiki/Avocado"}
                width="100%"
                height="600px"
              />
            </ModalBody>
            <ModalFooter>
              <Button variant="ghost" onClick={handleCloseModal}>
                Close
              </Button>
            </ModalFooter>
          </ModalContent>
        </Modal>
      )}
    </>
  );
};

export default ViewIframeModal;
