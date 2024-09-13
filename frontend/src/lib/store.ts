import { create } from "zustand";
import Bookmark from "@/lib/bookmark";
import { ChatLog } from "@/lib/chat";
import omit from "lodash-es/omit";
import { DocidMap, Document } from "@/lib/document";

interface KesslerState {
  bookmarks: Bookmark[]; // array of ids
  chats: { [key: string]: ChatLog }; // array of ids
  mainChatId: string;
  chatList: string[];
  documents: DocidMap;
  addBookmark: (id: string, type: string) => void;
  newChat: () => Promise<void>;
  addMessage: (chat_id: string, message: string) => void;
  setMainChat: (chat_id: string) => void;
}

const useKesslerStore = create<KesslerState>((set) => ({
  bookmarks: [],
  chats: {},
  chatList: [],
  mainChatId: "",
  documents: {},
  addBookmark: (id, type) => set((state) => ({})),
  newChat: async () => {
    // create new chat via API
    // add to chatList
    const chat_id = "new_chat_id";
    set((state) => ({
      ...state,
      chats: { ...state.chats, [chat_id]: new ChatLog() },
    }));
    set((state) => ({ ...state, chatList: [...state.chatList, chat_id] }));
  },
  sortChatsByDate: () => {
    set((state) => {
      const sorted = state.chatList.sort((a, b) => {
        const atime = new Date(state.chats[a].timestamp).getTime();
        const btime = new Date(state.chats[b].timestamp).getTime();
        return atime - btime;
      });
      return { ...state, chatList: sorted };
    });
  },
  addMessage: (chat_id, message) => {},
  setMainChat: (chat_id) => set((state) => ({ ...state, mainChatId: chat_id })),
  cacheDocids: (doc_ids: string[]) =>
    set((state) => {
      const updated_docs = doc_ids.reduce(
        (acc, docid) => {
          if (!state.documents[docid]) {
            acc[docid] = new Document({ docid: docid });
          }
          return acc;
        },
        { ...state.documents }
      );

      return { ...state, documents: updated_docs };
    }),
}));
