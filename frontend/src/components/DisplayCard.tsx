const IntroductionCard = () => {
  return <div>This is a test for new users</div>;
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
