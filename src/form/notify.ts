import { FormInstructions } from "../types/types.ts";

export class Notify {
  constructor() {}

  public sendNotification = async (instruction: FormInstructions) => {
    await chrome.runtime.sendMessage({
      action: "notify",
      title: "My Approval Extension",
      message: instruction.value,
    });
  };
}
