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
