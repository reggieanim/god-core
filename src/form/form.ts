import { FormInstructions, Options } from "../types/types";
import { validate } from "../helpers/validation/validate";
import { Block, Click, Eval, Notify, Select, Text, Wait } from "./index";
import { clearStorage } from "../helpers/functions/functions";

export class Form {
  private readonly instructions: FormInstructions[];
  private options: Options;
  private readonly templateUrl: string;

  constructor(data: unknown, templateUrl: string = "") {
    if (!Array.isArray(data)) {
      throw new Error("Wrong instructions format in form");
    }

    this.instructions = data;
    this.options = data.slice(0, -1) as Options;
    this.templateUrl = templateUrl;
  }

  public async start() {
    let countRetrys = 0;

    const retry = this.options.retry ?? 1;
    const scroll = this.options.scroll ?? 0;
    const skip = this.options.skip ?? "";

    while (true) {
      if (retry === countRetrys) {
        await clearStorage([this.templateUrl]);
        break;
      }
      if (skip === "true") {
        break;
      }

      let contextDocument: Document = document;

      if (this.options.iframeSelector) {
        const iframe = document.querySelector(this.options.iframeSelector);
        if (!iframe) {
          console.log("Error in iframe: element not found");
          continue;
        }

        // @ts-ignore
        const iframeDocument = iframe.contentDocument || iframe.contentWindow.document;
        if (!iframeDocument) {
          console.log("Error in iframe: document not found");
          continue;
        }
        contextDocument = iframeDocument;
      }

      for (const instruction of [...this.instructions]) {
        await this.runForm(instruction, contextDocument);
        this.instructions.shift();

        await chrome.storage.session.set({ [this.templateUrl]: [...this.instructions] });
      }

      // @ts-ignore
      contextDocument.defaultView.scrollBy(0, scroll);
      await new Promise((resolve) => setTimeout(resolve, 2000));
      countRetrys++;
    }
  }

  private async runForm(instruction: FormInstructions, contextDocument: Document) {
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
          await new Eval().conditionalEvaluate(instruction, contextDocument, this.templateUrl);
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
    }
  }
}
