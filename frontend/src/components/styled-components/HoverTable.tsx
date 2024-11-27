import clsx from "clsx";

const HoverTable = ({
  children,
  className,
  onClick,
}: {
  children: React.ReactNode;
  className?: string;
  onClick?: () => void;
}) => {
  return (
    <table
      className={clsx(
        "table table-pin-rows hover:bg-base-200 transition duration-500 ease-out",
        className,
      )}
      onClick={onClick}
    >
      {children}
    </table>
  );
};

export default HoverTable;
