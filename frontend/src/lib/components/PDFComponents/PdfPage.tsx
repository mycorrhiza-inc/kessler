import React, { useRef, useEffect } from "react";
// import pdfjs, { PDFPageProxy, PDFRenderParams } from "pdfjs-dist";
import pdfjs from "pdfjs-dist";
import "./PdfPage.css";

interface PdfPageProps {
  page: any | null;
  scale: number;
}

const PdfPage: React.FC<PdfPageProps> = React.memo(({ page, scale }) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const textLayerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!page) {
      return;
    }

    const viewport = page.getViewport({ scale });

    // Prepare canvas using PDF page dimensions
    const canvas = canvasRef.current;
    if (canvas) {
      const context = canvas.getContext("2d");
      if (context) {
        canvas.height = viewport.height;
        canvas.width = viewport.width;

        // Render PDF page into canvas context
        const renderContext: any = {
          canvasContext: context,
          viewport: viewport
        };
        const renderTask = page.render(renderContext);
        renderTask.promise.then(() => {
          // console.log("Page rendered");
        });

        page.getTextContent().then((textContent : any) => {
          // console.log(textContent);
          if (!textLayerRef.current) {
            return;
          }

          // Pass the data to the method for rendering of text over the pdf canvas.
          // @ts-ignore
          pdfjs.renderTextLayer({
          // @ts-ignore
            textContent,
            container: textLayerRef.current!,
            viewport: viewport,
            textDivs: []
          });
        });
      }
    }
  }, [page, scale]);

  return (
    <div className="PdfPage">
      <canvas ref={canvasRef} />
      <div ref={textLayerRef} className="PdfPage__textLayer" />
    </div>
  );
});

export default PdfPage;
