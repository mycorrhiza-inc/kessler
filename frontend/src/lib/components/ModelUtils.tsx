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
import ReactDOM from "react-dom";
import MarkdownRenderer from "./MarkdownRenderer";
import PDFViewer from "./PDFViewer";
import { FileType } from "../interfaces";

import { Table, Thead, Tbody, Tr, Th, Td } from "@chakra-ui/react";
const DynamicModalIframe: React.FC<{
  iframeURL: string;
  onClose: () => void;
}> = ({ iframeURL, onClose }) => {
  console.log(iframeURL);
  const {
    isOpen,
    onOpen,
    onClose: closeModal,
  } = useDisclosure({ defaultIsOpen: true });

  const toast = useToast();

  useEffect(() => {
    if (isOpen) {
      (async () => {
        // Do stuff when the modal opens
      })();
    }
  }, [isOpen]);

  const handleClose = () => {
    closeModal();
    onClose();
  };

  return ReactDOM.createPortal(
    <Modal isOpen={isOpen} onClose={handleClose} size="6xl">
      <ModalOverlay />
      <ModalContent maxW="80%">
        <ModalHeader></ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <iframe
            title="Iframe Content"
            src={iframeURL || "https://en.wikipedia.org/wiki/Avocado"}
            width="100%"
            height="1000px"
          />
        </ModalBody>
        <ModalFooter>
          <Button variant="ghost" onClick={handleClose}>
            Close
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>,
    document.body,
  );
};

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
        <DynamicModalIframe iframeURL={iframeURL} onClose={handleCloseModal} />
      )}
    </>
  );
};

// PDF Viewer Coming soon
// <Document
//   file="https://raw.githubusercontent.com/mycorrhizainc/examples/main/CO%20Clean%20Energy%20Plan%20Info%20Sheet.pdf"
//   onLoadSuccess={onDocumentLoadSuccess}
// >
//   <Page pageNumber={pageNumber} />
// </Document>
export default ViewIframeModal;
