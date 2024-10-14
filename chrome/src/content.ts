import { Form } from "./form/form";

console.log("Content script loaded");

chrome.runtime.onMessage.addListener((message, _sender, sendResponse) => {
  if (message.action === "executeTemplate") {
    console.log(message.tabID);

    const result = executeTemplate(message.template);
    sendResponse({ result });
  }
});

const executeTemplate = (template: any[]): any => {
  return template.map(async (item) => {
    if (Array.isArray(item)) {
      const [action, ...args] = item;
      switch (action) {
        case "form":
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
