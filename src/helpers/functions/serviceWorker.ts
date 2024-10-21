import { FunctionMap, Instruction } from "../../types/types";
import { insertCustomBanner, removeCustomBanner, setWindowToFalse } from "./functions";

export const getCurrentlyActiveTab = async (): Promise<chrome.tabs.Tab[]> => {
  return chrome.tabs.query({
    active: true,
    currentWindow: true,
  });
};

export const createNotification = (title: string, message: string): void => {
  chrome.notifications.create({
    type: "basic",
    iconUrl: "./icons/approval80.jpg",
    title: title,
    message: message,
    priority: 2,
  });
};

export const executeScriptInActiveTab = async (funcName: string, args: string, tabID: number): Promise<void> => {
  const functionMap: FunctionMap = {
    insertCustomBanner: insertCustomBanner,
    removeCustomBanner: removeCustomBanner,
    setWindowToFalse: setWindowToFalse,
  };

  chrome.scripting
    .executeScript({
      target: { tabId: tabID },
      func: functionMap[funcName],
      args: [args],
    })
    .then(() => console.log("Executed Banner"));
};

export const createNewTab = async (startingUrl: string, windowID?: number): Promise<number | undefined> => {
  const createTabOptions: chrome.tabs.CreateProperties = { url: startingUrl, active: true };

  if (windowID !== undefined) {
    createTabOptions.windowId = windowID;
  }

  const currentTab = await chrome.tabs.create(createTabOptions);

  return currentTab?.id;
};

export const createNewWindow = async (startingUrl: string): Promise<number | undefined> => {
  const createWindowOptions: chrome.windows.CreateData = {};

  if (startingUrl !== undefined && startingUrl !== "") {
    createWindowOptions.url = startingUrl;
  }

  const newWindow = await chrome.windows.create(createWindowOptions);
  return newWindow?.id;
};

export const getStartingURL = (rawInstructions: string): string => {
  const parsedInstructions: Instruction[] = JSON.parse(rawInstructions);
  return parsedInstructions[0].instructions[0].startingUrl;
};

export const CreateNewWindowOrTab = async (rawInstructions: string) => {
  const parsedInstructions: Instruction[] = JSON.parse(rawInstructions);

  if (Array.isArray(parsedInstructions)) {
    let idx = 0;
    let url: string;
    for (const instructionSet of parsedInstructions) {
      url = idx == 0 ? instructionSet.instructions[0].startingUrl : "";
      const windowID = await createNewWindow(url);
      if (idx == 0) {
        idx++;
        continue;
      }
      if (windowID !== undefined) {
        const instructions = instructionSet.instructions;

        if (Array.isArray(instructions)) {
          for (const instruction of instructions) {
            const instructionStartingURL = instruction.startingUrl;
            if (instructionStartingURL && instructionStartingURL !== "") {
              await createNewTab(instructionStartingURL, windowID);
            }
          }
        }
      }
    }
  }
};

export const processInstruction = async (instruction: Instruction, tabID: number): Promise<void> => {
  for (const config of instruction.instructions) {
    try {
      await chrome.tabs.sendMessage(tabID, {
        action: "executeTemplate",
        template: config.template,
      });
      await chrome.storage.session.setAccessLevel({ accessLevel: "TRUSTED_AND_UNTRUSTED_CONTEXTS" });

      console.log("Performed actions successfully");
    } catch (error) {
      console.error("Error processing instruction:");
    }
  }
};
