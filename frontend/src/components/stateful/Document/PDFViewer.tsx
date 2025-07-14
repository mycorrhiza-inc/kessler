"use client";

import { useCallback, useEffect, useState } from "react";
import { pdfjs, Document, Page } from "react-pdf";
import "react-pdf/dist/esm/Page/AnnotationLayer.css";
import "react-pdf/dist/esm/Page/TextLayer.css";
import dynamic from "next/dynamic";

import "./PDFViewer.css";
import LoadingSpinner from "@/components/style/misc/LoadingSpinner";
import ErrorMessage from "@/components/style/messages/ErrorMessage";

import type { PDFDocumentProxy } from "pdfjs-dist";

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
  const [networkErr, setNetworkErr] = useState("");
  const [renderErr, setRenderErr] = useState("");
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
      console.log(errString);
      setNetworkErr(errString);
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
  if (networkErr || renderErr) {
    var message;
    if (networkErr) {
      message = "Sorry, our servers couldnt fetch the pdf :(";
    }
    if (renderErr) {
      message = "Sorry, the client couldnt render the pdf :(";
    }
    return <ErrorMessage message={message} error={networkErr || renderErr} />;
  }
  if (fileData === undefined) {
    return <ErrorMessage error="PDF Data is Undefined" />;
  }
  const onError = (err: Error) => {
    console.log(err);
    setRenderErr(err.message);
  };
  return (
    <PDFViewerRaw file={fileData} setLoading={() => {}} onLoadError={onError} />
  );
};

function PDFViewerRaw({
  file,
  setLoading,
  onLoadError,
}: {
  file: PDFFile;
  setLoading: (val: boolean) => void;
  onLoadError: (val: Error) => void;
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

  // throw new Error("panicing now to avoid having to throw on a new dependancy, check PDFViewer.tsx for more details")
  // ResizeObserver(containerRef, resizeObserverOptions, onResize);

  function onDocumentLoadSuccess({ numPages: nextNumPages }: any): void {
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
            onLoadError={onLoadError}
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
