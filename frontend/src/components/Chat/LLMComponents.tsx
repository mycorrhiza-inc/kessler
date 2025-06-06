import { DocketPill, TextPill } from "../Tables/TextPills";

export const LinkFile = ({ uuid, text }: { uuid: string; text?: string }) => {
  return (
    <TextPill
      text={text || uuid}
      href={`/files/${uuid}`}
      seed={uuid}
    ></TextPill>
  );
};

export const LinkDocket = ({
  docket_named_id,
  text,
}: {
  docket_named_id: string;
  text?: string;
}) => {
  return <DocketPill text={text} docketId={docket_named_id} />;
};

// export const LinkOrg = ({
//   org_name,
//   text,
// }: {
//   docket_named_id: string;
//   text?: string;
// }) => {}
