const IntroductionCard = () => {
  return (
    <div className="card card-side bg-primary border-primary outline-secondary border-2 text-neutral-content shadow-xl w-5/12">
      <figure>
        <img src="/logo-big.png" alt="Kessler Logo" />
      </figure>
      <div className="card-body text-primary-content">
        <h2 className="card-title">Welcome To Kessler!</h2>
        <p>Use the buttons at the bottom to start.</p>
        <p>
          When searching and use the filters to search any goverment database
          that we support.
        </p>
        <p>
          When chatting, the chatbot has access to all of our data. Feel free to
          ask it any questions!
        </p>
        <p>
          Tip: If you want to limit what documents the chatbot can see, use the
          filters in the searchbox to limit the scope to a specific docket or
          author.
        </p>
        <div className="card-actions justify-end">
          <button className="btn btn-accent">Start Tour</button>
        </div>
      </div>
    </div>
  );
};
const ReturningCard = () => {
  return <div>This is a test for returning users</div>;
};

const DisplayCard = ({ cardType }: { cardType: string }) => {
  if (cardType === "introduction") {
    return <IntroductionCard />;
  }
  if (cardType === "returning") {
    return <ReturningCard />;
  } else {
    return <></>;
  }
};
export default DisplayCard;
