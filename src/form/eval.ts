import { waitForElement } from "../helpers/functions/functions.ts";
import { FormInstructions } from "../types/types.ts";

export class Eval {
  constructor() {}

  public detectFieldPresence = async (instruction: FormInstructions, page: Document) => {
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
  };

  private evaluate = async (instruction: FormInstructions, page: Document) => {
    function evaluate(element: HTMLElement, codeToExecute: string) {
      const func = new Function("return " + codeToExecute).bind(element);
      return func();
    }

    try {
      const element = await waitForElement(instruction.field, instruction.timeout, page);

      if (element instanceof HTMLElement) {
        const evaluation = evaluate(element, instruction.evalExpression);
        return evaluation();
      }
    } catch (err) {
      console.error(err);
      console.log(
        `Error executing code on element: ${instruction.field} with function: ${instruction.evalExpression}`
      );
    }
  };
}
