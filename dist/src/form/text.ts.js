import { waitForElement } from "/src/helpers/functions/functions.ts.js";
export class Text {
  constructor() {
  }
  executeText = async (instruction, page) => {
    try {
      const element = await waitForElement(instruction.field, instruction.timeout, page);
      await this.setElementValue(element, instruction.value);
      console.log(`Successfully executed text action: ${instruction.description}`);
    } catch (error) {
      console.error(`Error executing text action: ${instruction.description}`, error);
    }
  };
  setElementValue = async (element, value) => {
    if (element instanceof HTMLInputElement || element instanceof HTMLTextAreaElement) {
      element.focus();
      element.value = value;
      element.dispatchEvent(new Event("input", { bubbles: true, cancelable: true }));
      element.dispatchEvent(new Event("change", { bubbles: true, cancelable: true }));
    } else {
      throw new Error("Provided element is not an HTMLInputElement or HTMLTextAreaElement");
    }
  };
}
