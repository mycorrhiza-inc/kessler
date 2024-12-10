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

type OrganizationTableSchema = {
  Name: string;
  DocumentCount: number;
};

const OrganizationTable = ({
  orgList,
}: {
  orgList: OrganizationTableSchema[];
}) => {
  return (
    <table className="table table-pin-rows">
      <thead>
        <tr>
          <td className="w-[80%]">Name</td>
          <td className="w-[20%]">Documents Authored</td>
          {/* <td>Description</td> */}
        </tr>
      </thead>
      <tbody>
        {orgList.map((org: any) => (
          <tr
            key={org.DocketID}
            className="border-base-300 hover:bg-base-200 transition duration-500 ease-out"
          >
            <td colSpan={2} className="p-0">
              <Link href={`/orgs/${org.ID}`} className="flex w-full">
                <div className="w-[80%] px-4 py-3">{org.Name}</div>
                <div className="w-[20%] px-4 py-3">{org.DocumentCount}</div>
                {/* <div className="flex-1 px-4 py-3">{org.Description}</div> */}
              </Link>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};
const OrganizationTableSimple = () => {
  const { data, error, isLoading } = useSWRImmutable(
    `${apiURL}/v2/public/organizations/list`,
    organizationsListGet,
  );
  return (
    <>
      {isLoading && <LoadingSpinner loadingText="Loading Organizations" />}
      {error && <p>Failed to load organizations {String(error)}</p>}
      {!isLoading && !error && data != undefined && (
        <OrganizationTable orgList={data} />
      )}
    </>
  );
};

export default OrganizationTableSimple;
