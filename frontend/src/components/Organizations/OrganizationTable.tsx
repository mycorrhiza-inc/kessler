"use client";
import { apiURL } from "@/lib/env_variables";
import axios from "axios";
import Link from "next/link";
import useSWRImmutable from "swr/immutable";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import {
  OrganizationSchemaComplete,
  OrganizationSchemaCompleteValidator,
} from "@/lib/types/backend_schemas";

const organizationsListGet = (url: string) => {
  const cleanData = (response: any) => {
    console.log(response.data);
    const return_data: any[] = response.data;
    if (return_data.length == 0 || return_data == undefined) {
      return [];
    }
    return return_data;
    const valid_data = return_data.map(
      (item): OrganizationSchemaComplete =>
        OrganizationSchemaCompleteValidator.parse(return_data),
    );
    return valid_data;
  };
  return axios.get(url).then((res) => cleanData(res));
};

const OrganizationTable = () => {
  const { data, error, isLoading } = useSWRImmutable(
    `${apiURL}/v2/public/organizations/list`,
    organizationsListGet,
  );
  const convoList = data;
  console.log("Convo List:", convoList);
  return (
    <>
      {isLoading && <LoadingSpinner loadingText="Loading Organizations" />}
      {error && <p>Failed to load organizations {String(error)}</p>}
      {!isLoading && !error && convoList != undefined && (
        <table className="table table-pin-rows">
          <thead>
            <tr>
              <td>Name</td>
              <td>Description</td>
            </tr>
          </thead>
          <tbody>
            {convoList.map((org: any) => (
              <tr
                key={org.DocketID}
                className="border-base-300 hover:bg-base-200 transition duration-500 ease-out"
              >
                <td colSpan={2} className="p-0">
                  <Link href={`/orgs/${org.ID}`} className="flex w-full">
                    <div className="flex-1 px-4 py-3">{org.Name}</div>
                    <div className="flex-1 px-4 py-3">{org.Description}</div>
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

export default OrganizationTable;
