import { apiURL } from "@/lib/env_variables";
import axios from "axios";
import Link from "next/link";
import useSWRImmutable from "swr/immutable";
import LoadingSpinner from "../styled-components/LoadingSpinner";

const conversationsListAll = async () => {
  const response = await axios.get(`${apiURL}/v2/public/conversations/list`);
  console.log(response.data);
  return response.data;
};

const ConversationTable = () => {
  const { data, error, isLoading } = useSWRImmutable(conversationsListAll);
  const convoList = data;
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
