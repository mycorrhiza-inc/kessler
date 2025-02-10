import clsx from "clsx";

export const TableStyled = ({
  header_row_content,
  table_content,
  colgroup,
}: {
  header_row_content: React.ReactNode;
  table_content: React.ReactNode;
  colgroup?: React.ReactNode;
}) => {
  return (
    <div
      className={
        true
          ? "max-h-[1000px] overflow-y-scroll overflow-x-scroll"
          : "overflow-y-auto"
      }
    >
      <table
        className={clsx(
          "w-full divide-y divide-gray-200  border-collaps table lg:table-fixed md:table-auto sm:table-auto z-1",
          // pinClassName,
        )}
      >
        {/* disable pinned rows due to the top row overlaying the filter sidebar */}
        <colgroup>{colgroup}</colgroup>
        <thead>
          <tr className="border-b border-gray-200">{header_row_content}</tr>
        </thead>

        <tbody>{table_content}</tbody>
      </table>
    </div>
  );
};
