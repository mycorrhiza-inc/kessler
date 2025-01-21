import axios from "axios";
import { FileQueryFilterFields, backendFilterGenerate } from "./filters";

export interface Message {
  role: string;
  content: string;
  citations: any[]; // Define a format for search results and include them here
  key: symbol;
}
export const exampleChatHistory: Message[] = [
  {
    role: "user",
    content: "What is a black hole?",
    key: Symbol(),
    citations: [],
  },
  {
    role: "assistant",
    content:
      "A black hole is a region of space where the gravitational pull is so strong that not even light can escape from it.",
    key: Symbol(),
    citations: [],
  },
  {
    role: "user",
    content: "How are black holes formed?",
    key: Symbol(),
    citations: [],
  },
  {
    role: "assistant",
    content:
      "Black holes are formed when massive stars collapse under their own gravity at the end of their life cycle.",
    key: Symbol(),
    citations: [],
  },
  {
    role: "user",
    content: "Can black holes be seen?",
    key: Symbol(),
    citations: [],
  },
  {
    role: "assistant",
    content:
      "No, black holes cannot be seen directly because their gravitational pull prevents light from escaping, but their presence can be inferred by observing the effects on nearby objects.",
    key: Symbol(),
    citations: [],
  },
  {
    role: "user",
    content: "What would happen if you fell into a black hole?",
    key: Symbol(),
    citations: [],
  },
  {
    role: "assistant",
    content:
      "If you fell into a black hole, you would experience extreme gravitational forces and time dilation. Ultimately, you would be stretched and compressed in a process known as spaghettification.",
    key: Symbol(),
    citations: [],
  },
];

export interface ChatMessageInterface {
  role: string;
  content: string;
  key: symbol;
}

export const getUpdatedChatHistory = async (
  chatHistory: Message[],
  ragFilters: FileQueryFilterFields,
  chatUrl: string,
  model?: string,
) => {
  const backendFilters = backendFilterGenerate(ragFilters);
  let result_message = await fetch(chatUrl, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      accept: "application/json",
    },
    body: JSON.stringify({
      model: model,
      chat_history: chatHistory,
      filters: backendFilters,
    }),
  })
    .then((resp) => {
      if (resp.status < 200 || resp.status > 299) {
        console.log("failed request with status " + resp.status);
        console.log(resp);
        return "failed request";
      }
      return resp.json();
    })
    .then((data) => {
      if (!data.message) {
        console.log("no message in data");
        console.log(data);
        return "failed request";
      }
      console.log("got data");
      console.log(data);
      console.log("Returning Message:");
      console.log(data.message);
      return data.message;
    })
    .catch((e) => {
      console.log("error making request");
      console.log(JSON.stringify(e));
      return "encountered exception while fetching data";
    });
  let chat_response: Message;

  if (typeof result_message === "string") {
    chat_response = {
      role: "assistant",
      key: Symbol(),
      content: result_message,
      citations: [],
    };
  } else {
    chat_response = {
      role: "assistant",
      key: Symbol(),
      content: result_message.content,
      citations: result_message.citations,
    };
  }
  return [...chatHistory, chat_response];
};

export interface ChatLogInterface {
  id: string;
  loaded: boolean;
  title?: string;
  timestamp: string;
  messages: ChatMessageInterface[];
  addMessage: (message: ChatMessageInterface) => void;
  RecomputeMessage: (id: string) => void;
  loadLog: () => void;
  new: () => void;
}

export class ChatLog implements ChatLogInterface {
  id: string = "";
  title?: string | undefined;
  loaded: boolean = false;
  timestamp: string = "";
  messages: ChatMessageInterface[] = [];

  addMessage: (message: ChatMessageInterface) => void = (
    message: ChatMessageInterface,
  ) => {};
  RecomputeMessage: (id: string) => void = () => {};
  loadLog: () => void = () => {
    axios
      .get("/api/chat/?id=" + this.id, {})
      .then((result: any) => {
        if (result.data) {
          this.messages = result.data.chat_history;
        }
      })
      .catch((err) => {
        console.log(err);
      });
    this.loaded = true;
  };
  new: () => void = () => {
    axios
      .post("/api/chat/new", {})
      .then((result: any) => {
        if (result.data) {
          this.id = result.data.id;
        }
      })
      .catch((err) => {
        console.log(err);
      });
  };
  getHistory = () => {
    if (!this.messages) {
      return "";
    }

    let chat_hist = this.messages.map((m) => {
      let { key, ...rest } = m;
      return rest;
    });
    return chat_hist;
  };

  async sendMessage() {
    const chatUrl = "/api/chat/?id=" + this.id;
    let result = await fetch(chatUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        accept: "application/json",
      },
      body: JSON.stringify({
        model: "llama-70b",
        chat_history: this.getHistory(),
      }),
    })
      .then((resp) => {
        if (resp.status < 200 || resp.status > 299) {
          return "failed request";
        }
        return resp.json();
      })
      .then((data) => {
        return data;
      })
      .catch((e) => {
        console.log("error making request");
        console.log(JSON.stringify(e));
      });

    this.messages = [...this.messages, result.message];
  }
}
