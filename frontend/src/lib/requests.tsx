import { FileType } from "./interfaces";

export default function AuthenticatedFetch() {
  const authenticatedFetch = async (
    // same type signature as fetch
    resource: string | URL | Request,
    args: RequestInit,
  ) => {
    return fetch(resource, {
      /* default args for fetch */
      method: "POST", // *GET, POST, PUT, DELETE, etc.
      mode: "cors", // no-cors, *cors, same-origin
      cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
      credentials: "same-origin", // include, *same-origin, omit
      headers: {
        "Content-Type": "application/json",
        "Access-Control-Allow-Origin": "*",
        // clerky bearer token
        // Authorization: `Bearer ${await getToken()}`,
        // 'Content-Type': 'application/x-www-form-urlencoded',
      },
      redirect: "follow", // manual, *follow, error
      referrerPolicy: "no-referrer", // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
      // body: JSON.stringify({}), // body data type must match "Content-Type" header
      ...args,
    });
  };

  return authenticatedFetch;
}

export const GetAllFiles = async (): Promise<any> => {
  let result = await fetch("/api/files/all", {
    method: "GET",
    headers: {
      // "Content-Type": "application/json",
      Accept: "application/json",
      "Access-Control-Allow-Origin": "*",
    },
    // body: JSON.stringify({ url: link, title: "Textual content", isUrl: true }),
  }).then((e) => {
    console.log("completed request");
    console.log(e);
    if (e.status < 200 || e.status > 299) {
      console.log(`error adding links:\n${e}`);
      return "failed request";
    }
    let res_j = e.json();
    if (res_j == undefined) { return [] }
    else { return res_j }
  });
  if (result == undefined) return [];
  let out = result.map((f: any) => {
    return {
      id: f.id.toString(),
      url: "",
      title: f.name,
    };
  });
  console.log(out);
  return out;
};
