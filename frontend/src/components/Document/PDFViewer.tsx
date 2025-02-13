"use client";

import { useCallback, useEffect, useState } from "react";
import { useResizeObserver } from "@wojtekmaj/react-hooks";
import { pdfjs, Document, Page } from "react-pdf";
import "react-pdf/dist/esm/Page/AnnotationLayer.css";
import "react-pdf/dist/esm/Page/TextLayer.css";
import dynamic from "next/dynamic";

import "./PDFViewer.css";

import type { PDFDocumentProxy } from "pdfjs-dist";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import ErrorMessage from "../ErrorMessage";
import { fi } from "date-fns/locale";

// TODO : Inline at some point so we dont get screwed by a malicious cdn.
pdfjs.GlobalWorkerOptions.workerSrc =
  "https://cdn.jsdelivr.net/npm/pdfjs-dist@4.8.69/build/pdf.worker.mjs";
// pdfjs.GlobalWorkerOptions.workerSrc = new URL('pdfjs-dist/build/pdf.worker.min.mjs', import.meta.url, ).toString();

const options = {
  cMapUrl: "/cmaps/",
  standardFontDataUrl: "/standard_fonts/",
};

const resizeObserverOptions = {};

const maxWidth = 800;

type PDFFile = string | File | null;

const PDFViewer = ({ file }: { file: string }) => {
  const [downloading, setDownloading] = useState(true);
  const [err, setErr] = useState("");
  const [fileData, setFileData] = useState<File>();

  const getPDFData = async (url: string): Promise<File> => {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(
        `Server Returned Bad Error Code: ${response.status} ${response.statusText}`,
      );
    }
    const blob = await response.blob();
    return new File([blob], "document.pdf", { type: "application/pdf" });
  };
  const getPDFDataHandleErrors = async (url: string): Promise<File> => {
    setDownloading(true);
    try {
      const result = await getPDFData(url);
      setDownloading(false);
      return result;
    } catch (error) {
      setDownloading(false);
      const errString = `Failed to fetch PDF: ${error}`;
      setErr(errString);
      throw new Error(errString);
    }
  };
  useEffect(() => {
    let active = true;
    load();
    return () => {
      active = false;
    };
    async function load() {
      const res = await getPDFDataHandleErrors(file);
      if (!active) {
        return;
      }
      setFileData(res);
    }
  }, [file]);
  if (downloading) {
    return <LoadingSpinner loadingText="Loading PDF Data" />;
  }
  if (err) {
    return <ErrorMessage error={err} />;
  }
  if (fileData === undefined) {
    return <ErrorMessage error="PDF Data is Undefined" />;
  }
  return <PDFViewerRaw file={fileData} setLoading={() => {}} />;
};

function PDFViewerRaw({
  file,
  setLoading,
}: {
  file: PDFFile;
  setLoading: (val: boolean) => void;
}) {
  const [numPages, setNumPages] = useState<number>();
  const [containerRef, setContainerRef] = useState<HTMLElement | null>(null);
  const [containerWidth, setContainerWidth] = useState<number>();

  const onResize = useCallback<ResizeObserverCallback>((entries) => {
    const [entry] = entries;

    if (entry) {
      setContainerWidth(entry.contentRect.width);
    }
  }, []);

  useResizeObserver(containerRef, resizeObserverOptions, onResize);

  function onDocumentLoadSuccess({
    numPages: nextNumPages,
  }: PDFDocumentProxy): void {
    setLoading(false);
    setNumPages(nextNumPages);
  }

  return (
    <div className="Example">
      <div className="Example__container">
        <div className="Example__container__document" ref={setContainerRef}>
          <Document
            file={file}
            onLoadSuccess={onDocumentLoadSuccess}
            options={options}
            onLoadError={console.error}
          >
            {Array.from(new Array(numPages), (el, index) => (
              <Page
                key={`page_${index + 1}`}
                pageNumber={index + 1}
                width={
                  containerWidth ? Math.min(containerWidth, maxWidth) : maxWidth
                }
              />
            ))}
          </Document>
        </div>
      </div>
    </div>
  );
}
export default dynamic(() => Promise.resolve(PDFViewer), {
  ssr: false,
});
