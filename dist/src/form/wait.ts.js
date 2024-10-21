export class Wait {
  constructor() {
  }
  async wait(timeToWait) {
    return new Promise((resolve) => {
      setTimeout(resolve, Number(timeToWait) * 1e3);
    });
  }
}
