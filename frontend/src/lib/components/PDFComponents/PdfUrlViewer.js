import React, { useEffect, useState, useCallback, useRef } from "react";
import PropTypes from "prop-types";
import pdfjs from "pdfjs-dist";
import PdfViewer from "./PdfViewer";

const PdfUrlViewer = props => {
  const { url, ...others } = props;

  const pdfRef = useRef();

  const [itemCount, setItemCount] = useState(0);

  useEffect(() => {
    var loadingTask = pdfjs.getDocument(url);
    loadingTask.promise.then(
      pdf => {
        pdfRef.current = pdf;

        setItemCount(pdf._pdfInfo.numPages);

        // Fetch the first page
        var pageNumber = 1;
        pdf.getPage(pageNumber).then(function(page) {
          console.log("Page loaded");
        });
      },
      reason => {
        // PDF loading error
        console.error(reason);
      }
    );
  }, [url]);

  const handleGetPdfPage = useCallback(index => {
    return pdfRef.current.getPage(index + 1);
  }, []);

  return (
    <PdfViewer
      {...others}
      itemCount={itemCount}
      getPdfPage={handleGetPdfPage}
    />
  );
};

PdfUrlViewer.propTypes = {
  url: PropTypes.string.isRequired
};

export default PdfUrlViewer;
