import axios from "axios";

const chatUrl = "/api/chat";
const modelToSend = "llama-3.1-8b";

export interface ChatMessageInterface {
  role: string;
  content: string;
  key: symbol;
}

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
        model: modelToSend,
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
