import { until } from "/src/helpers/functions/functions.ts.js";
export class Block {
  constructor() {
  }
  executeBlock = async (instruction) => {
    await chrome.runtime.sendMessage({
      action: "scripting",
      function: "insertCustomBanner",
      args: instruction.value
    });
    await chrome.runtime.sendMessage({
      action: "scripting",
      function: "removeCustomBanner",
      args: ""
    });
    await until(() => window.startAutofill == true);
    await chrome.runtime.sendMessage({
      action: "scripting",
      function: "setWindowToFalse",
      args: ""
    });
  };
}
