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
        const [currentTab] = await chrome.tabs.query({
          active: true,
          currentWindow: true,
        });

        if (currentTab.id) {
          const response = await chrome.tabs.sendMessage(currentTab.id, {
            action: "executeTemplate",
            template: config.template,
          });
          await chrome.storage.local.set({ tabID: currentTab.id });

          console.log("Performed actions successfully", response);
        }
      } catch (error) {
        console.error("Error processing instruction:", error);
      }
    }
  };
}
