"use client";
import React, { useState, useRef } from "react";
import PdfUrlViewer from "./PdfUrlViewer";

const PDFApplet: React.FC = () => {
  const [scale, setScale] = useState<number>(1);
  const [page, setPage] = useState<number>(1);
  const windowRef = useRef<any>(null);
  const url = "https://raw.githubusercontent.com/mycorrhizainc/examples/main/CO%20Clean%20Energy%20Plan%20Info%20Sheet.pdf";

  const scrollToItem = () => {
    windowRef.current && windowRef.current.scrollToItem(page - 1, "start");
  };

  return (
    <div className="App">
      <h1>Pdf Viewer</h1>
      <div>
        <input
          type="number"
          value={page}
          onChange={(e) => setPage(Number(e.target.value))}
        />
        <button type="button" onClick={scrollToItem}>
          goto
        </button>
        Zoom
        <button type="button" onClick={() => setScale((v) => v + 0.1)}>
          +
        </button>
        <button type="button" onClick={() => setScale((v) => v - 0.1)}>
          -
        </button>
      </div>
      <br />
      <PdfUrlViewer url={url} scale={scale} windowRef={windowRef} />
      <p>
        https://mozilla.github.io/pdf.js/examples/index.html#interactive-examples
      </p>
      <p>https://react-window.now.sh/#/examples/list/variable-size</p>
    </div>
  );
};

export default PDFApplet;


