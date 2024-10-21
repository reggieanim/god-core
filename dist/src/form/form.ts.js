import { validate } from "/src/helpers/validation/validate.ts.js";
import { Block, Click, Eval, Notify, Select, Text, Wait } from "/src/form/index.ts.js";
import { clearStorage } from "/src/helpers/functions/functions.ts.js";
export class Form {
  instructions;
  options;
  constructor(data) {
    if (!Array.isArray(data)) {
      throw new Error("Wrong instructions format in form");
    }
    this.instructions = data;
    this.options = data.slice(0, -1);
  }
  async start() {
    let countRetrys = 0;
    const retry = this.options.retry ?? 1;
    const scroll = this.options.scroll ?? 0;
    const skip = this.options.skip ?? "";
    while (true) {
      if (retry === countRetrys) {
        await clearStorage();
        break;
      }
      if (skip === "true") {
        break;
      }
      let contextDocument = document;
      if (this.options.iframeSelector) {
        const iframe = document.querySelector(this.options.iframeSelector);
        if (!iframe) {
          console.log("Error in iframe: element not found");
          continue;
        }
        const iframeDocument = iframe.contentDocument || iframe.contentWindow.document;
        if (!iframeDocument) {
          console.log("Error in iframe: document not found");
          continue;
        }
        contextDocument = iframeDocument;
      }
      for (let _ = 0; _ < this.instructions.length; _++) {
        const instruction = this.instructions.shift();
        await this.runForm(instruction, contextDocument);
        chrome.storage.session.set({ args: this.instructions });
      }
      contextDocument.defaultView.scrollBy(0, scroll);
      await new Promise((resolve) => setTimeout(resolve, 2e3));
      countRetrys++;
    }
  }
  async runForm(instruction, contextDocument) {
    console.log("Running Form with instruction", instruction);
    if (!validate(instruction)) {
      throw new Error("Invalid fields");
    }
    if (instruction.skip === "true") {
      return;
    }
    try {
      switch (instruction.kind) {
        case "text":
          await new Text().executeText(instruction, contextDocument);
          break;
        case "leftClick":
          await new Click().executeLeftClick(instruction, contextDocument);
          break;
        case "rightClick":
          await new Click().executeRightClick(instruction, contextDocument);
          break;
        case "wait":
          await new Wait().wait(instruction.value);
          break;
        case "condEval":
          await new Eval().detectFieldPresence(instruction, contextDocument);
          break;
        case "notify":
          await new Notify().sendNotification(instruction);
          break;
        case "select":
          await new Select().inputSelect(instruction, contextDocument);
          break;
        case "block":
          await new Block().executeBlock(instruction);
          break;
        default:
          break;
      }
    } catch (error) {
      console.error("An error occurred:", error);
      return;
    } finally {
      return;
    }
  }
}
