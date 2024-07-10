import React, { useEffect, useState, useCallback, useRef, RefObject } from "react";
import PropTypes from "prop-types";
import pdfjs from "pdfjs-dist";
import PdfViewer from "./PdfViewer";

interface PdfUrlViewerProps {
  url: string;
  scale: number;
  windowRef: RefObject<any>;
  [key: string]: any;
}

const PdfUrlViewer: React.FC<PdfUrlViewerProps> = (props) => {
  const { url, scale, windowRef, ...others } = props;

  const pdfRef = useRef<pdfjs.PDFDocumentProxy | null>(null);

  const [itemCount, setItemCount] = useState(0);

  useEffect(() => {
    const loadingTask = pdfjs.getDocument(url);
    loadingTask.promise.then(
      (pdf) => {
        pdfRef.current = pdf;
        setItemCount(pdf.numPages);

        // Fetch the first page
        const pageNumber = 1;
        pdf.getPage(pageNumber).then(function (page) {
          console.log("Page loaded");
        });
      },
      (reason) => {
        // PDF loading error
        console.error(reason);
      }
    );
  }, [url]);

  const handleGetPdfPage = useCallback((index: number) => {
    if (!pdfRef.current) {
      throw new Error("PDF document not loaded");
    }
    return pdfRef.current.getPage(index + 1);
  }, []);

  return (
    <PdfViewer
      {...others}
      itemCount={itemCount}
      getPdfPage={handleGetPdfPage}
      scale={scale}
      windowRef={windowRef}
    />
  );
};

//PdfUrlViewer.propTypes = {
//  url: PropTypes.string.isRequired,
//  scale: PropTypes.number.isRequired,
//  windowRef: PropTypes.object.isRequired,
//};

export default PdfUrlViewer;


