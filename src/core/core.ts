import { Instruction } from "../types/types";

chrome.runtime.onInstalled.addListener(() => {
  console.log("Extension installed");
});

export class InstructionProcessor {
  constructor(private rawInstructions: string) {}

  public start() {
    const instructions: Instruction[] = JSON.parse(this.rawInstructions);
    instructions.forEach(this.processInstruction);
  }

  private processInstruction = async (instruction: Instruction) => {
    for (const config of instruction.instructions) {
      try {
        // const current_tab = await chrome.tabs.create({ url: config.startingUrl });
        const [currentTab] = await InstructionProcessor.getCurrentlyActiveTab();

        if (currentTab.id) {
          const response = await chrome.tabs.sendMessage(currentTab.id, {
            action: "executeTemplate",
            template: config.template,
          });
          chrome.storage.session.setAccessLevel({ accessLevel: "TRUSTED_AND_UNTRUSTED_CONTEXTS" });
          await chrome.storage.session.set({ tabID: currentTab.id });

          console.log("Performed actions successfully", response);
        }
      } catch (error) {
        console.error("Error processing instruction:", error);
      }
    }
  };

  static getCurrentlyActiveTab = async () => {
    return chrome.tabs.query({
      active: true,
      currentWindow: true,
    });
  };

  static createNotification = (title: string, message: string) => {
    chrome.notifications.create({
      type: "basic",
      iconUrl: "./icons/approval80.jpg",
      title: title,
      message: message,
      priority: 2,
    });
  };
}

chrome.runtime.onMessage.addListener((request, _sender, _sendResponse) => {
  if (request.action === "notify") {
    InstructionProcessor.createNotification(request.title, request.message);
  }
});

chrome.webNavigation.onCommitted.addListener(function (details) {
  if (
    details.transitionType === "auto_subframe" ||
    details.transitionType === "form_submit" ||
    details.transitionType === "reload"
  ) {
    const continueProcess = async () => {
      const storageRetrievalResult = await chrome.storage.session.get(["args"]);
      const [currentTab] = await InstructionProcessor.getCurrentlyActiveTab();

      if (storageRetrievalResult !== undefined && currentTab.id) {
        const response = await chrome.tabs.sendMessage(currentTab.id, {
          action: "continueExecutingTemplate",
          template: storageRetrievalResult.args,
        });

        return response;
      }
    };

    continueProcess()
      .then((result) => {
        console.log(result);
      })
      .catch((error) => {
        console.log("An error occurred:", error);
      });
  }
});
