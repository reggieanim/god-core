import { Form } from "/src/form/form.ts.js";
console.log("Content script loaded");
chrome.runtime.onMessage.addListener((message, _sender, sendResponse) => {
  let result;
  switch (message.action) {
    case "ping":
      sendResponse({ status: "ready" });
      return true;
    case "executeTemplate":
      result = executeTemplate(message.template);
      break;
    case "continueExecutingTemplate":
      result = continueExecutingTemplate(message.template);
      break;
    default:
      result = { error: "Unknown action" };
  }
  sendResponse({ result });
});
const executeTemplate = (template) => {
  return template.map(async (item) => {
    if (Array.isArray(item)) {
      const [action, ...args] = item;
      switch (action) {
        case "form":
          chrome.storage.session.set({ args: [args] });
          await new Form(args).start();
          break;
        case "print":
          console.log("Printing args", args);
          break;
        default:
          console.log(`Unknown action: ${action}`);
      }
    }
    return item;
  });
};
let isExecuting = false;
const continueExecutingTemplate = async (template) => {
  if (isExecuting) {
    return;
  }
  isExecuting = true;
  if (template !== void 0) {
    await new Form(template).start();
  }
  isExecuting = false;
};
