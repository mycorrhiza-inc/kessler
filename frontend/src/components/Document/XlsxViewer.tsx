const XlsxViewer = ({ file }: { file: string }) => {
  return (
    <div className="flex flex-col items-center justify-center p-8 m-4 rounded-lg bg-success/10 text-success-content">
      <div className="text-5xl mb-4">⚠️</div>
      <h3 className="text-xl font-bold mb-2">
        Unforuntately, we dont have support for viewing Excel Spreadsheets in
        the browser yet :(
      </h3>
      <p className="text-center mb-4">
        We are working on it though, in the meantime you can just download the
        file.
      </p>
      <a
        href={file}
        target="_blank"
        onClick={() => window.location.reload()}
        className="btn btn-success-content btn-outline"
      >
        Download File
      </a>
    </div>
  );
};

export default XlsxViewer;
