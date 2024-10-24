import { FormInstructions } from "../../types/types";

const kindMap: { [key: string]: boolean } = {
  text: true,
  notify: true,
  select: true,
  leftClick: true,
  rightClick: true,
  wait: true,
  eval: true,
  condEval: true,
  block: true,
};

class Validator {
  errors: string[] = [];

  check(condition: boolean, field: string, message: string) {
    if (!condition) {
      this.errors.push(`${field}: ${message}`);
    }
  }

  valid(): boolean {
    return this.errors.length === 0;
  }
}

export const validate = (instructions: FormInstructions): boolean => {
  if (Object.keys(instructions).length === 0) return true;

  const isValidDescription = typeof instructions.description === "string";
  const isValidField = typeof instructions.field === "string";
  const isValidValue = typeof instructions.value === "string";
  const insKind = instructions.kind;
  const isValidKind = typeof insKind === "string";
  const doesNotMatchExpectedKind = insKind ? kindMap[insKind] : false;

  const v = new Validator();
  v.check(isValidDescription, "description", "Needs a 'description' property of string");
  v.check(isValidField, "field", "Needs a 'field' property of 'text'");
  v.check(isValidValue, "value", "Needs a 'value' property of string");
  v.check(isValidDescription, "description", "Needs a description property");
  v.check(isValidKind, "kind", "Needs a 'kind' property");
  v.check(
    doesNotMatchExpectedKind,
    "noKind",
    "Needs a 'kind' property of 'text' || condEval || select || leftClick || rightClick' || 'wait' || 'notify' || 'eval' || block"
  );

  if (!v.valid()) {
    for (const error of v.errors) {
      console.error(error);
    }
    console.log("Could not validate");
  }

  return v.valid();
};
