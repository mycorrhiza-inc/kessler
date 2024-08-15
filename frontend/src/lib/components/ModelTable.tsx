import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableContainer,
  Box,
  Spinner,
  Select,
  Text,
  Button,
  Center,
} from "@chakra-ui/react";
import { ModelType } from "../interfaces";

import { useEffect, useState } from "react";

import Paginator from "./Paginator";
import ViewIframeModal from "./ModelUtils";
interface FilePageBrowserProps {
  fileUrl: string;
  data: any;
}

interface RowData {
  data: ModelType;
}

interface Layout {
  columns: {
    key: string;
    label: string;
    width: string;
    enabled: boolean;
  }[];
  showExtraFeatures: boolean;
  showDisplayText: boolean;
}

export const defaultLayout: Layout = {
  columns: [
    { key: "name", label: "Name", width: "60%", enabled: true },
    { key: "date", label: "Date", width: "20%", enabled: true },
  ],
  showExtraFeatures: true,
  showDisplayText: true,
};

interface ModelTableProps {
  models: ModelType[];
  layout: Layout;
}

const ModelTable: React.FC<ModelTableProps> = ({ models, layout }) => {
  const [fileState, setFileState] = useState<RowData[]>(
    models.map((model) => ({ data: model })),
  );

  function truncateString(str: string, length = 60) {
    return str.length < length ? str : str.slice(0, length - 3) + "...";
  }
  function getFieldFromFile(key: string, model: ModelType): string {
    // Please shut up, I know what I'm doing
    // @ts-ignore
    let result = model[key];
    return result !== undefined ? String(result) : "Unknown";
  }

  const layoutFiltered: Layout = {
    ...layout,
    columns: layout.columns.filter((column) => column.enabled),
  };

  return (
    <TableContainer>
      <Table>
        <Thead>
          <Tr>
            {layoutFiltered.columns.map((col) => (
              <Th key={col.key} width={col.width}>
                {col.label}
              </Th>
            ))}
            {layoutFiltered.showExtraFeatures && (
              <>
                <Th width="6%">View</Th>
                <Th width="2%">Status</Th>
              </>
            )}
          </Tr>
        </Thead>
        <Tbody>
          {fileState.map((model) => (
            <>
              <Tr key={model.data.id}>
                {layoutFiltered.columns.map((col) => (
                  <Td key={col.key}>
                    {truncateString(getFieldFromFile(col.key, model.data))}
                  </Td>
                ))}
                {layoutFiltered.showExtraFeatures && (
                  <>
                    <Td>
                      <ViewIframeModal
                        // @ts-ignore
                        iframeUrl={model.data.visualization_url}
                        buttonName="View Model"
                      />
                    </Td>
                    <Td>
                      <a href={model.data.data_url}>
                        <Button>Download Data</Button>
                      </a>
                    </Td>
                  </>
                )}
              </Tr>
            </>
          ))}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

interface ModelPageBrowserProps {
  modelUrl: string;
  data: any;
}

const ModelPageBrowser: React.FC<ModelPageBrowserProps> = ({
  modelUrl,
  data,
}) => {
  const [models, setModels] = useState<ModelType[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [maxPage, setMaxPage] = useState(1);
  const [numResults, setNumResults] = useState(10);
  const [layout, setLayout] = useState(defaultLayout);
  const fetchFiles = async () => {
    setLoading(true);
    const response = await fetch(
      `${modelUrl}?num_results=${numResults}&page=${page}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      },
    );
    const result = await response.json();
    console.log(result);
    setModels(result[0]);
    setMaxPage(result[1]);
    setLoading(false);
  };

  useEffect(() => {
    fetchFiles();
  }, [modelUrl, data, page, numResults]);

  return (
    <Box>
      {loading ? (
        <Center>
          <Spinner />
        </Center>
      ) : (
        <ModelTable models={models} layout={layout} />
      )}
      <Paginator
        page={page}
        setPage={setPage}
        maxPage={maxPage}
        numResults={numResults}
        setNumResults={setNumResults}
      />

      <ViewIframeModal
        // @ts-ignore
        iframeUrl="https://nimbus.kessler.xyz/schema/swagger"
        buttonName="Admin Panel"
      />
    </Box>
  );
};

export default ModelPageBrowser;
