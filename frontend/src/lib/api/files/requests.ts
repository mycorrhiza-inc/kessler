const AddLocalFile = () => {};
const DeleteLocalFile = () => {};

export const AddLink = async (formData: any) => {
    let link = formData.link
    console.log(`link:\n${link}`);

    await fetch("http://localhost:5000/api/files/add", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
        "Access-Control-Allow-Origin": "*",
      },
      body: JSON.stringify({ url: link, title: "Textual content" }),
    })
      .then((e) => {
        console.log(`successfully added link "${link}":\n${e}`);
      })
      .catch((e) => {
        console.log(`error adding links:\n${e}`);
      });
  };
const RemoveLink = () => {};