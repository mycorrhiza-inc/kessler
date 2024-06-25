"use client";
import {
  Table,
  Thead,
  Tbody,
  Tfoot,
  Tr,
  Th,
  Td,
  TableCaption,
  TableContainer,
  Center,
  Checkbox,
  Box,
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  Button,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { FileType } from "../interfaces/file";
import { LoadingSpinner } from "@saas-ui/react";
import {
  FiInfo,
  FiTrash2,
  FiAward,
  FiFeather,
  FiChevronDown,
  FiDownload,
} from "react-icons/fi";
import { GetAllFiles } from "../requests";
import FileUploadButton from "./FileUploader";

enum sorts {
  none,
  SortPublished,
  SortPublishedDec,
  SortModified,
  SortModifiedDec,
  SortStatus,
  SortStatusBy,
}
interface RowData {
  selected: boolean;
  data: FileType;
}
interface FileExploreState {
  files: RowData[];
  // which pagination page
  page: number;
  // number of rows per page of the file
  resultsPerPage: number;
  // the id order to be shown
  sortStyle: sorts;
  // index of the last selected item
  lastSelected: number[];
}

const initState: FileExploreState = {
  files: [
    {
      selected: false,
      data: {
        id: "1",
        url: "",
        title: "",
      },
    },
    {
      selected: false,
      data: {
        id: "2",
        url: "",
        title: "",
      },
    },
    {
      selected: false,
      data: {
        id: "3",
        url: "",
        title: "",
      },
    },
    {
      selected: false,
      data: {
        id: "4",
        url: "",
        title: "",
      },
    },
  ],
  // which pagination page
  page: 0,
  // number of rows per page of the file
  resultsPerPage: 20,
  // the id order to be shown
  sortStyle: sorts.SortModified,
  lastSelected: [],
};

interface RowProps {
  file: RowData;
  updateSelected: (id: string) => void;
}

export default function FileExplorer() {
  const [state, setState] = useState<FileExploreState>(initState);
  const [loading, setLoading] = useState(true);
  const loadFiles = async () => {
    let files = await GetAllFiles();
    const stateFiles: RowData[] = files.map((f: FileType) => {
      return {
        selected: false,
        data: f,
      };
    });
    console.log(stateFiles);
    setState({ ...state, files: stateFiles });
    setLoading(false);
  };
  useEffect(() => {
    if (loading) loadFiles();
  });

  const getSelectedIndex = (id: string) => {
    for (let i = 0; i < state.files.length; i++) {
      if (state.files[i].data.id == id) return i;
    }
    return -1;
  };

  const getAllSelected = () => {
    let s = state.files.filter((file) => file.selected == true);
    // get an array of indexes of selected elements
    return s.map((_file, index) => {
      return index;
    });
  };
  const pushLast = (id: string) => {
    let last = state.lastSelected;
    const index = getSelectedIndex(id);
    last.push(index);
    setState({
      ...state,
      lastSelected: last,
    });
  };

  const popLast = (id: string) => {
    let last = state.lastSelected;
    const index = getSelectedIndex(id);
    let _ = last.pop();
    setState({
      ...state,
      lastSelected: last,
    });
  };

  const updateSelected = (id: string) => {
    let newFiles = state.files.map((file) => {
      if (file.data.id == id) {
        // set the last selected
        file.selected = !file.selected;
      }
      return file;
    });
    setState({
      ...state,
      files: newFiles,
    });
  };

  function ExplorerRow({ file, updateSelected }: RowProps) {
    const [loading, setLoading] = useState(false);
    return (
      <Tr key={file.data.id}>
        {/* Select */}
        <Td>
          <Box
            onClick={(e) => {
              const meta = e.nativeEvent.metaKey;
              const shift = e.nativeEvent.shiftKey;

              if (meta && shift) {
                // re-set the checked box
                updateSelected(file.data.id);
                return;
              }

              if (shift) {
                // get the index of the current id
                const last = state.lastSelected.pop();

                if (last == undefined) {
                  updateSelected(file.data.id);
                  return;
                }

                pushLast(file.data.id);

                const current = getSelectedIndex(file.data.id);
                let low;
                let high;
                if (current > last) {
                  low = last;
                  high = current;
                } else {
                  high = last;
                  low = current;
                }

                let newFiles = state.files.map((file, index) => {
                  if (index >= low && index <= high) {
                    // set the last selected
                    file.selected = true;
                  }
                  return file;
                });
                setState({
                  ...state,
                  files: newFiles,
                });
              }
              console.log(`last selected: ${state.lastSelected}`);
            }}
          >
            <Checkbox
              isChecked={file.selected}
              onChange={(e) => {
                console.log(`box checked`);
                console.log(e);
                updateSelected(file.data.id);
                if (file.selected) pushLast(file.data.id);
                else if (
                  state.lastSelected[state.lastSelected.length - 1] ==
                  getSelectedIndex(file.data.id)
                )
                  popLast(file.data.id);
                console.log(`last selected: ${state.lastSelected}`);
              }}
            />
          </Box>
        </Td>
        {/* Filename */}
        <Td>{file.data.title}</Td>
        {/* <Td>{file.data.datePublished.toString()}</Td>
        <Td>{file.data.dateAdded.toString()}</Td>
        <Td>{file.data.dateModified.toString()}</Td> */}
        {/* Status */}
        <Td>
          {loading && <LoadingSpinner />}
          {!loading && <FiInfo />}
        </Td>
      </Tr>
    );
  }

  return (
    <TableContainer>
      <Table>
        <Thead>
          <Tr>
            <Th width="2%"></Th>
            <Th width="96%"></Th>
            {/* <Th width="10%"></Th>
            <Th width="10%"></Th>
            <Th width="10%"></Th> */}
            <Th width="2%">
              <FileUploadButton />
              <Menu>
                <MenuButton as={Button} rightIcon={<FiChevronDown />}>
                  Actions
                </MenuButton>
                <MenuList>
                  <MenuItem icon={<FiDownload />}>Download</MenuItem>
                  <MenuItem icon={<FiAward />}>Add tag</MenuItem>
                  <MenuItem icon={<FiFeather />}> Add to Project</MenuItem>
                  <MenuItem icon={<FiTrash2 />}>Delete</MenuItem>
                </MenuList>
              </Menu>
            </Th>
          </Tr>
          <Tr>
            <Th width="2%">Select</Th>
            <Th width="96%">Filename</Th>
            {/* <Th width="10%">Published</Th>
            <Th width="10%">Added</Th>
            <Th width="10%">Modified</Th> */}
            <Th width="2%">Status</Th>
          </Tr>
        </Thead>
        <Tbody>
          {loading && (
            <Tr>
              <Td></Td>
              <Td>
                <Center>
                  <LoadingSpinner />
                </Center>
              </Td>
            </Tr>
          )}
          {!loading &&
            state.files.map((file, index) => {
              if (file.data.title == "") return;

              return (
                <ExplorerRow
                  key={file.data.id}
                  file={file}
                  updateSelected={updateSelected}
                />
              );
            })}
        </Tbody>
      </Table>
    </TableContainer>
  );
}
