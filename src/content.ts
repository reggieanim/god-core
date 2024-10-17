import { Form } from "./form/form";

console.log("Content script loaded");

chrome.runtime.onMessage.addListener((message, _sender, sendResponse) => {
  let result;
  switch (message.action) {
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

const executeTemplate = (template: any[]): any => {
  return template.map(async (item) => {
    if (Array.isArray(item)) {
      const [action, ...args] = item;
      switch (action) {
        case "form":
          chrome.storage.session.set({ args: [args] });
          await new Form(args).start();
          break;
        case "print":
          console.log("Printing");
          break;
        default:
          console.log(`Unknown action: ${action}`);
      }
    }
    return item;
  });
};

const continueExecutingTemplate = async (template: any[]): Promise<void> => {
  await new Form(template).start();
};
