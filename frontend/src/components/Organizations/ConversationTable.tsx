"use client";
import { apiURL } from "@/lib/env_variables";
import axios from "axios";
import Link from "next/link";
import useSWRImmutable from "swr/immutable";
import LoadingSpinner from "../styled-components/LoadingSpinner";

const conversationsListAll = (redundant_key: string) => {
  const cleanData = (response: any) => {
    console.log(response.data);
    const return_data: any[] = response.data;
    if (return_data.length == 0 || return_data == undefined) {
      return [];
    }
    return return_data;
  };
  return axios
    .get(`${apiURL}/v2/public/conversations/list`)
    .then((res) => cleanData(res));
};

type ConversationTableSchema = {
  Name: string;
  DocketID: string;
  DocumentCount: number;
  State: string;
};

const ConversationTable = ({
  convoList,
}: {
  convoList: ConversationTableSchema[];
}) => {
  return (
    <table className="table table-pin-rows">
      {/* disable pinned rows due to the top row overlaying the filter sidebar */}
      <thead>
        <tr>
          <td className="w-[60%]">Name</td>
          <td className="w-[20%]">ID</td>
          <td className="w-[10%]">Document Count</td>
          <td className="w-[10%]">State</td>
        </tr>
      </thead>
      <tbody>
        {convoList.map((convo: any) => (
          <tr
            key={convo.DocketID}
            className="border-base-300 hover:bg-base-200 transition duration-500 ease-out"
          >
            <td colSpan={4} className="p-0">
              <Link
                href={`/proceedings/${convo.DocketID}`}
                className="flex w-full"
              >
                <div className="w-[60%] px-4 py-3">{convo.Name}</div>
                <div className="w-[20%] px-4 py-3">{convo.DocketID}</div>
                <div className="w-[10%] px-4 py-3">{convo.DocumentCount}</div>
                <div className="w-[10%] px-4 py-3">{convo.State}</div>
              </Link>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};
const ConversationTableSimple = () => {
  const { data, error, isLoading } = useSWRImmutable(
    `redudant_key`,
    conversationsListAll,
  );
  console.log("Convo List:", data);
  return (
    <>
      {isLoading && <LoadingSpinner loadingText="Loading Conversations" />}
      {error && <p>Failed to load conversations {String(error)}</p>}
      {!isLoading && !error && data != undefined && (
        <ConversationTable convoList={data} />
      )}
    </>
  );
};

export default ConversationTableSimple;
