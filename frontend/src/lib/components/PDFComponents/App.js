import React, { useState, useRef } from "react";
import "./styles.css";
import PdfUrlViewer from "./PdfUrlViewer";

export default function App() {
  const [scale, setScale] = useState(1);
  const [page, setPage] = useState(1);
  const windowRef = useRef();
  const url = "oodmetrics.pdf";

  const scrollToItem = () => {
    windowRef.current && windowRef.current.scrollToItem(page - 1, "start");
  };

  return (
    <div className="App">
      <h1>Pdf Viewer</h1>
      <div>
        <input value={page} onChange={e => setPage(e.target.value)} />
        <button type="button" onClick={scrollToItem}>
          goto
        </button>
        Zoom
        <button type="button" onClick={() => setScale(v => v + 0.1)}>
          +
        </button>
        <button type="button" onClick={() => setScale(v => v - 0.1)}>
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
}
