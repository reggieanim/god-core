export class Wait {
  constructor() {}

  public async wait(timeToWait: string): Promise<void> {
    return new Promise((resolve) => {
      setTimeout(resolve, Number(timeToWait));
    });
  }
}
