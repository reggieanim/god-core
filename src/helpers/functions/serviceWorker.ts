import _ from "underscore";

import { FunctionMap, Instruction } from "../../types/types";
import { clearStorage, insertCustomBanner, removeCustomBanner, setWindowToFalse } from "./functions";

export const getCurrentlyActiveTab = async (): Promise<chrome.tabs.Tab[]> => {
  return chrome.tabs.query({
    active: true,
    currentWindow: true,
  });
};

export const createNotification = (title: string, message: string): void => {
  try {
    chrome.notifications.create({
      type: "basic",
      iconUrl: "./icons/approval48.png",
      title: title,
      message: message,
      priority: 2,
    });
  } catch (error) {
    console.error("Error creating a notification");
  }
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
    .then(() => console.log("Executed function:", funcName));
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

export const getStartingURLs = (rawInstructions: string): string[] => {
  const parsedInstructions: Instruction[] = JSON.parse(rawInstructions);

  return parsedInstructions.flatMap((instructionSet) =>
    instructionSet.instructions.map((instruction) => instruction.startingUrl).filter(Boolean)
  );
};

export const CreateNewWindowOrTab = async (rawInstructions: string) => {
  const parsedInstructions: Instruction[] = JSON.parse(rawInstructions);

  if (Array.isArray(parsedInstructions)) {
    let idx = 0;
    let url: string;
    for (const instructionSet of parsedInstructions) {
      url = idx == 0 ? instructionSet.instructions[0].startingUrl : "";
      const windowID = await createNewWindow(url);

      if (windowID !== undefined) {
        const instructions = instructionSet.instructions;

        if (Array.isArray(instructions)) {
          for (const instruction of instructions) {
            if (idx === 0) {
              idx++;
              continue;
            }
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

export const processInstruction = async (instruction: Instruction, tabID: number, url: string): Promise<void> => {
  for (const config of instruction.instructions) {
    if (config.startingUrl === url) {
      try {
        await chrome.tabs.sendMessage(tabID, {
          action: "executeTemplate",
          template: config.template,
          templateUrl: new URL(config.startingUrl).host,
        });
        await chrome.storage.session.setAccessLevel({ accessLevel: "TRUSTED_AND_UNTRUSTED_CONTEXTS" });

        console.log("Performed actions successfully");
      } catch (error) {
        console.error("Error processing instruction:");
      }
    }
  }
};

export async function onTabCreatedListener(tab: chrome.tabs.Tab): Promise<void> {
  const storageRetrievalResult = await chrome.storage.session.get(["startingUrls", "instructions"]);

  const startingUrls: string[] = storageRetrievalResult.startingUrls as string[];

  if (
    startingUrls !== undefined &&
    startingUrls.length > 0 &&
    startingUrls.includes(tab.pendingUrl!) &&
    storageRetrievalResult.instructions !== undefined &&
    tab.id !== undefined
  ) {
    const instructions: Instruction[] = JSON.parse(storageRetrievalResult.instructions);
    const maxAttempts = 75;
    let attempts = 0;

    const checkReadiness = async () => {
      if (attempts >= maxAttempts) {
        console.error("Max attempts reached. Page not ready.");
        return;
      }

      try {
        const response = await chrome.tabs.sendMessage(tab.id!, { action: "ping" });
        if (response && response.status === "ready") {
          await processInstructions(instructions, tab);
          await clearStorage(["startingUrls", "instructions"]);
        } else {
          attempts++;
          setTimeout(checkReadiness, 400);
        }
      } catch (error) {
        console.error("Error checking readiness:", error);
        attempts++;
        setTimeout(checkReadiness, 400);
      }
    };

    checkReadiness();
  }
}

async function processInstructions(instructions: Instruction[], tab: chrome.tabs.Tab) {
  instructions.forEach((instruction) => processInstruction(instruction, tab.id!, tab.pendingUrl!));
}
