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

const ConversationTable = () => {
  const { data, error, isLoading } = useSWRImmutable(
    `redudant_key`,
    conversationsListAll,
  );
  const convoList = data;
  console.log("Convo List:", convoList);
  return (
    <>
      {isLoading && <LoadingSpinner loadingText="Loading Conversations" />}
      {error && <p>Failed to load conversations {error}</p>}
      {!isLoading && !error && convoList != undefined && (
        <table className="table table-pin-rows">
          <thead>
            <tr>
              <td>Name</td>
              <td>State</td>
            </tr>
          </thead>
          <tbody>
            {convoList.map((convo: any) => (
              <tr key={convo.DocketID}>
                <td colSpan={2} className="p-0">
                  <Link
                    href={`/proceedings/${convo.DocketID}`}
                    className="flex w-full"
                  >
                    <div className="flex-1 px-4 py-3">{convo.DocketID}</div>
                    <div className="flex-1 px-4 py-3">{convo.State}</div>
                  </Link>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </>
  );
};

export default ConversationTable;
