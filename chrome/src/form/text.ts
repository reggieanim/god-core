import { waitForElement } from "../helpers/functions/functions";
import { FormInstructions } from "../types/types";

export class Text {
  constructor() {}

  public executeText = async (instruction: FormInstructions, page: Document) => {
    try {
      const element = await waitForElement(instruction.field, instruction.timeout, page);

      await this.setElementValue(element, instruction.value);

      console.log(`Successfully executed text action: ${instruction.description}`);
    } catch (error) {
      console.error(`Error executing text action: ${instruction.description}`, error);
    }
  };

  private setElementValue = async (element: Element, value: string): Promise<void> => {
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
