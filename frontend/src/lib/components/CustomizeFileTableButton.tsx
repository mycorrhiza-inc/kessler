import React, { useState } from "react";
import {
  Button,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Checkbox,
  Input,
  Stack,
  FormControl,
  FormLabel,
} from "@chakra-ui/react";

const CustomizeFileTableButton = ({ layout, setLayout }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [tempLayout, setTempLayout] = useState(layout);

  const onSave = () => {
    setLayout(tempLayout);
    setIsOpen(false);
  };

  const onToggleExtraFeatures = () => {
    setTempLayout({
      ...tempLayout,
      showExtraFeatures: !tempLayout.showExtraFeatures,
    });
  };

  const onChangeColumnWidthAndLabel = (index, field, value) => {
    const updatedColumns = [...tempLayout.columns];
    updatedColumns[index][field] = value;
    setTempLayout({
      ...tempLayout,
      columns: updatedColumns,
    });
  };

  return (
    <>
      <Button onClick={() => setIsOpen(true)}>Customize Table</Button>
      <Modal isOpen={isOpen} onClose={() => setIsOpen(false)}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Customize Table Layout</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Stack>
              {tempLayout.columns.map((column, index) => (
                <FormControl key={index}>
                  <FormLabel>{column.label}</FormLabel>
                  <Input
                    type="text"
                    placeholder="Column Label"
                    value={column.label}
                    onChange={(e) =>
                      onChangeColumnWidthAndLabel(
                        index,
                        "label",
                        e.target.value,
                      )
                    }
                  />
                  <Input
                    type="text"
                    placeholder="Column Width"
                    value={column.width}
                    onChange={(e) =>
                      onChangeColumnWidthAndLabel(
                        index,
                        "width",
                        e.target.value,
                      )
                    }
                  />
                  <Checkbox
                    isChecked={column.truncate}
                    onChange={(e) =>
                      onChangeColumnWidthAndLabel(
                        index,
                        "truncate",
                        e.target.checked,
                      )
                    }
                  >
                    Truncate
                  </Checkbox>
                </FormControl>
              ))}
              <Checkbox
                isChecked={tempLayout.showExtraFeatures}
                onChange={onToggleExtraFeatures}
              >
                Show Extra Features
              </Checkbox>
            </Stack>
          </ModalBody>
          <ModalFooter>
            <Button colorScheme="blue" onClick={onSave}>
              Save
            </Button>
            <Button variant="ghost" onClick={() => setIsOpen(false)}>
              Cancel
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
};

export default CustomizeFileTableButton;
