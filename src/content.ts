import { Form } from "./form/form";

console.log("Content script loaded");

chrome.runtime.onMessage.addListener((message, _sender, sendResponse) => {
  let result;
  switch (message.action) {
    case "ping":
      sendResponse({ status: "ready" });
      return true;
    case "executeTemplate":
      result = executeTemplate(message.template, message.templateUrl);
      break;
    case "continueExecutingTemplate":
      result = continueExecutingTemplate(message.template, message.templateUrl);
      break;
    default:
      result = { error: "Unknown action" };
  }
  sendResponse({ result });
});

let isTemplateExecuting = false;

const executeTemplate = async (template: any[], templateUrl: string): Promise<void> => {
  if (isTemplateExecuting) {
    console.log("Template is already executing. Skipping execution.");
    return;
  }
  isTemplateExecuting = true;

  try {
    for (const item of template) {
      if (Array.isArray(item)) {
        const [action, ...args] = item;
        switch (action) {
          case "form":
            const storageRetrievalResult = (await chrome.storage.session.get(["args"])) || {};
            storageRetrievalResult[templateUrl] = [...args];
            await chrome.storage.session.set({ args: storageRetrievalResult });
            await new Form(args, templateUrl).start();
            break;
          case "print":
            console.log("Printing args", args);
            break;
          default:
            console.log(`Unknown action: ${action}`);
        }
      }
    }
  } catch (error) {
    console.error("Error executing template:", error);
  } finally {
    isTemplateExecuting = false;
  }
};

let isExecuting = false;
const continueExecutingTemplate = async (template: any[], templateUrl: string): Promise<void> => {
  if (isExecuting) {
    return;
  }
  isExecuting = true;
  if (template !== undefined) {
    await new Form(template, templateUrl).start();
  }
  isExecuting = false;
};
