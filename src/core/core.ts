import { clearStorage } from "../helpers/functions/functions";
import {
  createNotification,
  executeScriptInActiveTab,
  processInstruction,
} from "../helpers/functions/serviceWorker";
import { Instruction } from "../types/types";

chrome.runtime.onInstalled.addListener(() => {
  console.log("Extension installed");
});

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === "notify") {
    createNotification(request.title, request.message);
  }

  if (
    request.action === "scripting" &&
    request.function !== undefined &&
    request.args !== undefined &&
    sender.tab?.id !== undefined
  ) {
    executeScriptInActiveTab(request.function, request.args, sender.tab?.id)
      .then(() => {
        sendResponse({ status: "success" });
      })
      .catch((error) => {
        console.error("An error occurred", error);
      });
  }
  return true;
});

chrome.webNavigation.onCommitted.addListener(function listener(
  details: chrome.webNavigation.WebNavigationTransitionCallbackDetails
): boolean {
  if (details.transitionType === "auto_subframe" || details.transitionType === "form_submit") {
    const continueProcess = async () => {
      const storageRetrievalResult = await chrome.storage.session.get(["args"]);

      if (storageRetrievalResult !== undefined && details.tabId) {
        const response = await chrome.tabs.sendMessage(details.tabId, {
          action: "continueExecutingTemplate",
          template: storageRetrievalResult.args,
        });

        await new Promise((resolve) => setTimeout(resolve, 90 * 1000));
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

  chrome.webNavigation.onCommitted.removeListener(listener);
  return true;
});

chrome.tabs.onCreated.addListener(async function listener(tab: chrome.tabs.Tab): Promise<boolean> {
  const storageRetrievalResult = await chrome.storage.session.get(["startingUrl", "instructions"]);
  if (
    tab.pendingUrl === storageRetrievalResult.startingUrl &&
    storageRetrievalResult.instructions !== undefined &&
    tab.id !== undefined
  ) {
    const instructions: Instruction[] = JSON.parse(storageRetrievalResult.instructions);

    chrome.tabs.onUpdated.addListener(function updatedListener(tabId, changeInfo, _updatedTab) {
      if (tabId === tab.id && changeInfo.status === "complete") {
        const maxAttempts = 75;
        let attempts = 0;

        const checkReadiness = async () => {
          if (attempts >= maxAttempts) {
            console.error("Max attempts reached. Page not ready.");
            chrome.tabs.onUpdated.removeListener(updatedListener);
            return;
          }

          try {
            const response = await chrome.tabs.sendMessage(tabId, { action: "ping" });
            if (response && response.status === "ready") {
              await processInstructions();
              await clearStorage();
              await chrome.storage.session.set({ tabID: tabId });
              chrome.tabs.onUpdated.removeListener(updatedListener);
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

        const processInstructions = async () => {
          instructions.forEach((instruction) => processInstruction(instruction, tab.id!));
        };

        checkReadiness();
      }
    });
  }

  chrome.tabs.onCreated.removeListener(listener);
  return true;
});
