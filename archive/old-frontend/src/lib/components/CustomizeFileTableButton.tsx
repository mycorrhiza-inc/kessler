import React, { useState } from "react";
import {
  Box,
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
  IconButton,
  Text,
} from "@chakra-ui/react";
import { FaWindowClose } from "react-icons/fa";
import {
  DndContext,
  closestCenter,
  useSensor,
  useSensors,
  MouseSensor,
  TouchSensor,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  verticalListSortingStrategy,
  useSortable,
  sortableKeyboardCoordinates,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";

import { TableLayout } from "../interfaces";
interface Column {
  key: string;
  label: string;
  width: string;
  enabled: boolean;
}

interface CustomizeFileTableButtonProps {
  layout: TableLayout;
  setLayout: React.Dispatch<React.SetStateAction<TableLayout>>;
}

const CustomizeFileTableButton: React.FC<CustomizeFileTableButtonProps> = ({
  layout,
  setLayout,
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [tempLayout, setTempLayout] = useState(layout);

  const sensors = useSensors(useSensor(MouseSensor), useSensor(TouchSensor));

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

  const onChangeColumnWidthAndLabel = (
    index: number,
    field: string,
    value: string | boolean,
  ) => {
    const updatedColumns = [...tempLayout.columns];
    updatedColumns[index] = {
      ...updatedColumns[index],
      [field]: value,
    };
    setTempLayout({
      ...tempLayout,
      columns: updatedColumns,
    });
  };

  const onRemoveColumn = (index: number) => {
    const updatedColumns = tempLayout.columns.filter((_, i) => i !== index);
    setTempLayout({
      ...tempLayout,
      columns: updatedColumns,
    });
  };

  const handleDragEnd = (event: any) => {
    const { active, over } = event;
    if (active.id !== over.id) {
      const oldIndex = tempLayout.columns.findIndex(
        (column) => column.key === active.id,
      );
      const newIndex = tempLayout.columns.findIndex(
        (column) => column.key === over.id,
      );

      setTempLayout({
        ...tempLayout,
        columns: arrayMove(tempLayout.columns, oldIndex, newIndex),
      });
    }
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
            <DndContext
              sensors={sensors}
              collisionDetection={closestCenter}
              onDragEnd={handleDragEnd}
            >
              <SortableContext
                items={tempLayout.columns.map((column) => column.key)}
                strategy={verticalListSortingStrategy}
              >
                <Stack>
                  {tempLayout.columns.map((column, index) => (
                    <SortableItem key={column.key} id={column.key}>
                      <Text>{column.label}</Text>
                      <Box
                        p={4}
                        shadow="md"
                        borderWidth="1px"
                        borderColor="green.200"
                        // backgroundColor={column.enabled ? "white" : "gray.200"}
                        display="flex"
                        alignItems="center"
                        justifyContent="space-between"
                      >
                        <FormControl>
                          <FormLabel>Column Width</FormLabel>
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
                            isChecked={column.enabled}
                            onChange={(e) =>
                              onChangeColumnWidthAndLabel(
                                index,
                                "enabled",
                                e.target.checked,
                              )
                            }
                          >
                            Enable
                          </Checkbox>
                        </FormControl>
                        <IconButton
                          aria-label="Delete"
                          icon={<FaWindowClose />}
                          size="sm"
                          onClick={() => onRemoveColumn(index)}
                        />
                      </Box>
                    </SortableItem>
                  ))}
                  <Checkbox
                    isChecked={tempLayout.showExtraFeatures}
                    onChange={onToggleExtraFeatures}
                  >
                    Show Extra Features
                  </Checkbox>
                </Stack>
              </SortableContext>
            </DndContext>
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

type SortableItemProps = {
  id: string;
  children: React.ReactNode;
};

const SortableItem: React.FC<SortableItemProps> = ({ id, children }) => {
  const { attributes, listeners, setNodeRef, transform, transition } =
    useSortable({ id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div ref={setNodeRef} style={style} {...attributes} {...listeners}>
      {children}
    </div>
  );
};

export default CustomizeFileTableButton;

//   <FormLabel>Column Label</FormLabel>
//   <Input
//     type="text"
//     placeholder="Column Label"
//     value={column.label}
//     onChange={(e) =>
//       onChangeColumnWidthAndLabel(
//         index,
//         "label",
//         e.target.value,
//       )
//     }
//   />
