import React, { useState, useCallback, useRef, useEffect, RefObject } from "react";
import { VariableSizeList } from "react-window";
import useResizeObserver from "use-resize-observer";
import PdfPage from "./PdfPage";
import Page from "./Page";

// Define the shape of the props
interface PdfViewerProps {
  width: number | string;
  height: number | string;
  itemCount: number;
  getPdfPage: (index: number) => Promise<any>;
  scale: number;
  gap: number;
  windowRef?: RefObject<VariableSizeList>;
}

const PdfViewer: React.FC<PdfViewerProps> = ({
  width = "100%",
  height = "400px",
  itemCount,
  getPdfPage,
  scale = 1,
  gap = 40,
  windowRef
}) => {

  const [pages, setPages] = useState<Array<any>>([]);

  const listRef = useRef<VariableSizeList>(null);

  const { ref, width: internalWidth = 400, height: internalHeight = 600 } = useResizeObserver<HTMLDivElement>();

  const fetchPage = useCallback(
    async (index: number) => {
      if (!pages[index]) {
        const page = await getPdfPage(index);
        setPages(prev => {
          const next = [...prev];
          next[index] = page;
          return next;
        });
        listRef.current?.resetAfterIndex(index);
      }
    },
    [getPdfPage, pages]
  );

  const handleItemSize = useCallback(
    (index: number) => {
      const page = pages[index];
      if (page) {
        const viewport = page.getViewport({ scale });
        return viewport.height + gap;
      }
      return 50;
    },
    [pages, scale, gap]
  );

  const handleListRef = useCallback(
    (elem: VariableSizeList | null) => {
      if (elem) {
        listRef.current = elem;
        if (windowRef) {
          windowRef.current = elem;
        }
      }
    },
    [windowRef]
  );

  useEffect(() => {
    listRef.current?.resetAfterIndex(0);
  }, [scale]);

  const style: React.CSSProperties = {
    width,
    height,
    border: "1px solid #ccc",
    background: "#ddd"
  };

  return (
    <div ref={ref} style={style}>
      <VariableSizeList
        ref={handleListRef}
        width={internalWidth}
        height={internalHeight}
        itemCount={itemCount}
        itemSize={handleItemSize}
      >
        {({ index, style }) => {
          fetchPage(index);
          return (
            <Page style={style}>
              <PdfPage page={pages[index]} scale={scale} />
            </Page>
          );
        }}
      </VariableSizeList>
    </div>
  );
};

export default PdfViewer;

