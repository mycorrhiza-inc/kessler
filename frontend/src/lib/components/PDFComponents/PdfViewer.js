import React, { useState, useCallback, useRef, useEffect } from "react";
import PropTypes from "prop-types";
import { VariableSizeList } from "react-window";
import useResizeObserver from "use-resize-observer";
import PdfPage from "./PdfPage";
import Page from "./Page";

const PdfViewer = props => {
  const { width, height, itemCount, getPdfPage, scale, gap, windowRef } = props;

  const [pages, setPages] = useState([]);

  const listRef = useRef();

  const {
    ref,
    width: internalWidth = 400,
    height: internalHeight = 600
  } = useResizeObserver();

  const fetchPage = useCallback(
    index => {
      if (!pages[index]) {
        getPdfPage(index).then(page => {
          setPages(prev => {
            const next = [...prev];
            next[index] = page;
            return next;
          });
          listRef.current.resetAfterIndex(index);
        });
      }
    },
    [getPdfPage, pages]
  );

  const handleItemSize = useCallback(
    index => {
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
    elem => {
      listRef.current = elem;
      if (windowRef) {
        windowRef.current = elem;
      }
    },
    [windowRef]
  );

  useEffect(() => {
    listRef.current.resetAfterIndex(0);
  }, [scale]);

  const style = {
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

PdfViewer.propTypes = {
  width: PropTypes.oneOfType([PropTypes.number, PropTypes.string]),
  height: PropTypes.oneOfType([PropTypes.number, PropTypes.string]),
  itemCount: PropTypes.number.isRequired,
  getPdfPage: PropTypes.func.isRequired,
  scale: PropTypes.number,
  gap: PropTypes.number,
  windowRef: PropTypes.object
};

PdfViewer.defaultProps = {
  width: "100%",
  height: "400px",
  scale: 1,
  gap: 40
};

export default PdfViewer;
