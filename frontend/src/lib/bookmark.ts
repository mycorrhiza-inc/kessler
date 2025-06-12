import axios from "axios";

type BookmarkType =
  | "document"
  | "organization"
  | "maybe"
  | "campaign"
  | "encounter"
  | undefined;

export interface BookmarkInterface {
  id: string;
  type: BookmarkType;
  title?: string;
  UpdateBookmarkTitle: (bookmarkId: string) => void;
}

export default class Bookmark implements BookmarkInterface {
  id: string = "";
  type: BookmarkType = undefined;
  title?: string | undefined;
  UpdateBookmarkTitle = async (bookmarkId: string) => {
    console.error(
      "Inserting to quickwit directly seems bad, currently it only does localhost, if you see this error and need to use it, its probably a good idea to refactor it to use the backend go api - Nicole",
    );
    throw new Error(
      "Inserting to quickwit directly seems bad, currently it only does localhost, if you see this error and need to use it, its probably a good idea to refactor it to use the backend go api - Nicole",
    );
    axios
      .post("http://localhost:4041/bookmarks/", {
        id: bookmarkId,
      })
      .then((result: any) => {
        if (result.data) {
          this.title = result.data.title;
        }
      })
      .catch((err) => {
        console.log(err);
      });
    return this;
  };

  constructor(id: string, type: BookmarkType, title?: string) {
    this.id = id;
    this.type = type;
    this.title = title;
  }
}
