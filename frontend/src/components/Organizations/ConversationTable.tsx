import { apiURL } from "@/lib/env_variables";
import axios from "axios";
import Link from "next/link";
import useSWRImmutable from "swr/immutable";
import LoadingSpinner from "../styled-components/LoadingSpinner";

const conversationsListAll = async (url: string) => {
  const response = await axios.get(url);
  console.log(response.data);
  const return_data: any[] = response.data;
  if (return_data.length == 0 || return_data == undefined) {
    return [];
  }
  return return_data;
};

const ConversationTable = () => {
  const { data, error, isLoading } = useSWRImmutable(
    `${apiURL}/v2/public/conversations/list`,
    conversationsListAll,
  );
  const convoList = data;
  console.log("convo list: " + convoList);
  return (
    <>
      {isLoading && <LoadingSpinner loadingText="Loading Conversations" />}
      {error && <p>Failed to load conversations {error}</p>}
      {!isLoading && !error && (
        <table className="table table-pin-rows">
          <thead>
            <tr>
              <td>Name</td>
              <td>Description</td>
            </tr>
          </thead>
          <tbody>
            {convoList.map((convo: any) => (
              <Link href={`/proceedings/${convo.name}`}>
                <tr>
                  <td>{convo.name}</td>
                  <td>{convo.description}</td>
                </tr>
              </Link>
            ))}
          </tbody>
        </table>
      )}
    </>
  );
};

export default ConversationTable;
