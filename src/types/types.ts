interface Instructions {
  name: string;
  startingUrl: string;
  template: any[];
}

export interface Instruction {
  headless: boolean;
  lender: string;
  inBrowser: boolean;
  saveState: boolean;
  stealth: boolean;
  slowMotion: number;
  trace: boolean;
  close: boolean;
  instructions: Instructions[];
}

export interface LogConfig {
  webhookURL: string;
}

export interface FormInstructions {
  description: string;
  field: string;
  value: string;
  shdType: boolean;
  kind: string;
  evalExpression: string;
  iframeSelector: string;
  timeout: number;
  skip: string;
  body: Object;
  fallback: Object;
  mute: boolean;
}

export interface Options {
  retry?: number;
  scroll?: number;
  skip?: string;
  iframeSelector?: string;
}

export type FunctionMap = {
  [key: string]: (value: string) => void;
};
