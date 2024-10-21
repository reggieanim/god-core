import { until } from "../helpers/functions/functions";
import { FormInstructions } from "../types/types";

export class Block {
  constructor() {}

  public executeBlock = async (instruction: FormInstructions): Promise<void> => {
    await chrome.runtime.sendMessage({
      action: "scripting",
      function: "insertCustomBanner",
      args: instruction.value,
    });
    await chrome.runtime.sendMessage({
      action: "scripting",
      function: "removeCustomBanner",
      args: "",
    });

    await until(() => window.startAutofill == true);

    await chrome.runtime.sendMessage({
      action: "scripting",
      function: "setWindowToFalse",
      args: "",
    });
  };
}
