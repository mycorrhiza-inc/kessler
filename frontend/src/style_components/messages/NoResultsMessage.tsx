const NoResultsMessage = () => {
  return (
    <div className="flex flex-col items-center justify-center p-8 m-4 rounded-lg bg-success/10 text-success-content">
      <div className="text-5xl mb-4"></div>
      <h3 className="text-xl font-bold mb-2">No Results Found...</h3>
      <p className="text-center mb-4">
        Perhaps you could try a different search query?
      </p>
    </div>
  );
};

export default NoResultsMessage;
