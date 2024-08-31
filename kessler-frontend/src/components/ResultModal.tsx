// Modal.js
import React from "react";
import { Box } from "@mui/joy";
import "./Modal.css"; // For styling

type ModalProps = {
  isOpen: boolean;
  onClose: () => void;
  children: React.ReactNode;
};

const ResultModal = ({ isOpen, onClose, children }: ModalProps) => {
  if (!isOpen) return null;

  return (
    <Box className="modal-content standard-box" sx={{}}>
      {/* box containing the  */}
      {/* children are components passed to result modals from special searches */}
      {children}
      <div> {/* header for modal */}</div>
    </Box>
  );
};

export default ResultModal;
