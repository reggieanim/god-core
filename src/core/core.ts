import {
  createNotification,
  executeScriptInActiveTab,
  onTabCreatedListener,
  removeTabUpdateListener,
  webNavigationOnCommittedListener,
} from "../helpers/functions/serviceWorker";

chrome.runtime.onInstalled.addListener(() => {
  console.log("Extension installed");
});

chrome.runtime.onMessage.addListener(async (request, sender, sendResponse) => {
  if (request.action === "notify" && sender.tab?.id !== undefined) {
    await createNotification(request.title, request.message);
  }

  if (request.action === "finished" && sender.tab?.id !== undefined) {
    chrome.tabs.onCreated.removeListener(onTabCreatedListener);
    removeTabUpdateListener();
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
});

export function checkHasListener() {
  if (!chrome.tabs.onCreated.hasListener(onTabCreatedListener)) {
    console.log("No Listeners registered");
    chrome.tabs.onCreated.addListener(onTabCreatedListener);
    chrome.webNavigation.onCommitted.addListener(webNavigationOnCommittedListener);
  }
}

// chrome.webNavigation.onCommitted.removeListener(webNavigationOnCommittedListener);
