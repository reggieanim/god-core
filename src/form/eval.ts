import { waitForElement } from "../helpers/functions/functions.ts";
import { FormInstructions } from "../types/types.ts";
import { Form } from "./form.ts";

export class Eval {
  constructor() {}

  public conditionalEvaluate = async (instruction: FormInstructions, page: Document, templateName: string) => {
    const result = await this.detectFieldPresence(instruction, page);
    if (result !== undefined && result) {
      await new Form(instruction.body, templateName).start();
      return;
    }

    if (instruction.fallback !== undefined && Object.keys(instruction.fallback).length !== 0) {
      await new Form(instruction.fallback, templateName).start();
      return;
    }
  };

  private detectFieldPresence = async (instruction: FormInstructions, page: Document) => {
    try {
      const element = await waitForElement(instruction.field, instruction.timeout, page);

      if (element instanceof HTMLElement) {
        const propertyOrMethod = instruction.evalExpression;

        if (propertyOrMethod in element) {
          const result =
            typeof element[propertyOrMethod] === "function"
              ? element[propertyOrMethod]()
              : element[propertyOrMethod];

          return result === instruction.value;
        } else {
          throw new Error(`Property or method ${propertyOrMethod} does not exist on the element`);
        }
      }
    } catch (err) {
      console.error(err);
      console.log(
        `Error checking presence of element: ${instruction.field}, with property: ${instruction.evalExpression}, that has the value: ${instruction.value}`
      );
      return false;
    }
    return false;
  };
}
