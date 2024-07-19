import { useState, useEffect } from "react";
import { FileType } from "../interfaces/file";
import { LoadingSpinner } from "@saas-ui/react";
import DocumentViewer from "./DocumentViewer";
import { IoMdCheckmarkCircleOutline } from "react-icons/io";
import { ImCross } from "react-icons/im";
import FileTable from "./FileTable";
interface RowData {
  selected: boolean;
  data: FileType;
}

interface FileTableProps {
  files: FileType[];
}
// TODO: Fix so that this actually recives a string of uuids instead of this BS
const FileTableUUIDOnly = ({ uuid_list }: { uuid_list: object[] }) => {
  console.log(uuid_list);
  // @ts-ignore
  const actual_uuid_list = uuid_list.map((obj) => obj.uuid);
  console.log(actual_uuid_list);

  const [files, setFiles] = useState<FileType[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchFiles = async () => {
    setLoading(true);
    const promises = actual_uuid_list.map((uuid) => fetchFile(uuid));
    const results = await Promise.all(promises);
    setFiles(results);
    setLoading(false);
  };

  useEffect(() => {
    fetchFiles();
  }, [uuid_list]);

  return <FileTable files={files} />;
};

const fetchFile = async (uuid: string) => {
  const response = await fetch(`/api/files/` + uuid, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });
  const result = await response.json();
  return result;
};

export default FileTableUUIDOnly;
