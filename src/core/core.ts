import {
  createNotification,
  executeScriptInActiveTab,
  // onTabCreatedListener,
} from "../helpers/functions/serviceWorker";

chrome.runtime.onInstalled.addListener(() => {
  console.log("Extension installed");
});

chrome.runtime.onMessage.addListener(async (request, sender, sendResponse) => {
  if (request.action === "notify" && sender.tab?.id !== undefined) {
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
});

chrome.storage.session.setAccessLevel({ accessLevel: "TRUSTED_AND_UNTRUSTED_CONTEXTS" });
