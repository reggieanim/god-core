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
  console.log("Called executeTemplate function");

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

const continueExecutingTemplate = async (template: any[], templateUrl: string): Promise<void> => {
  console.log("Called continueExecutingTemplate function");

  if (template !== undefined) await new Form(template, templateUrl).start();
};

(function executeContent() {
  const baseUrl = window.location.host;
  chrome.storage.session.get(["args"]).then(async (storageRetrievalResult) => {
    if (storageRetrievalResult.args?.[baseUrl] !== undefined && storageRetrievalResult.args[baseUrl].length > 0) {
      continueExecutingTemplate(storageRetrievalResult.args?.[baseUrl], baseUrl);
    }
  });
})();
