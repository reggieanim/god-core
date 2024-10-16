export const waitForElement = async (selector: string, timeout: number = 30, page: Document): Promise<Element> => {
  const startTime = Date.now();
  while (Date.now() - startTime < timeout * 1000) {
    const element = page.querySelector(selector);
    if (element) {
      return element;
    }
    await new Promise((resolve) => setTimeout(resolve, timeout));
  }
  throw new Error(`Timeout waiting for element: ${selector}`);
};

export const observeAndWaitForElement = async (
  selector: string,
  timeout: number = 5000,
  page: Document
): Promise<Element> => {
  return new Promise((resolve, reject) => {
    const element = page.querySelector(selector);
    if (element) {
      resolve(element);
    }

    const observer = new MutationObserver((_mutationsList, observer) => {
      const element = page.querySelector(selector);
      if (element) {
        observer.disconnect();
        resolve(element);
      }
    });

    observer.observe(page, { childList: true, subtree: true });

    setTimeout(() => {
      observer.disconnect();
      reject(new Error(`Timeout waiting for element: ${selector}`));
    }, timeout * 1000);
  });
};
