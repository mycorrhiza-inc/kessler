"use client";
import { apiURL } from "@/lib/env_variables";
import axios from "axios";
import Link from "next/link";
import useSWRImmutable from "swr/immutable";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import HoverTable from "../styled-components/HoverTable";

const organizationsListAll = (redundant_key: string) => {
  const cleanData = (response: any) => {
    console.log(response.data);
    const return_data: any[] = response.data;
    if (return_data.length == 0 || return_data == undefined) {
      return [];
    }
    return return_data;
  };
  return axios
    .get(`${apiURL}/v2/public/organizations/list`)
    .then((res) => cleanData(res));
};

const OrganizationTable = () => {
  const { data, error, isLoading } = useSWRImmutable(
    `redundant_key`,
    organizationsListAll,
  );
  const convoList = data;
  console.log("Convo List:", convoList);
  return (
    <>
      {isLoading && <LoadingSpinner loadingText="Loading Organizations" />}
      {error && <p>Failed to load organizations {error}</p>}
      {!isLoading && !error && convoList != undefined && (
        <HoverTable>
          <thead>
            <tr>
              <td>Name</td>
              <td>Description</td>
            </tr>
          </thead>
          <tbody>
            {convoList.map((convo: any) => (
              <tr key={convo.DocketID}>
                <td colSpan={2} className="p-0">
                  <Link href={`/orgs/${convo.ID}`} className="flex w-full">
                    <div className="flex-1 px-4 py-3">{convo.Name}</div>
                    <div className="flex-1 px-4 py-3">{convo.Description}</div>
                  </Link>
                </td>
              </tr>
            ))}
          </tbody>
        </HoverTable>
      )}
    </>
  );
};

export default OrganizationTable;
