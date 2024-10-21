import { Cipher, createCipheriv, createDecipheriv } from "crypto";

export class AesEncryptUtil {
  private static publicKey: string;
  private static ivKey: string;

  private static getAesIvKey(): string {
    if (AesEncryptUtil.ivKey) return AesEncryptUtil.ivKey;
    if (!process.env.VITE_AES_IV_KEY) throw new Error("AES IV KEY SHOULD BE SET");
    AesEncryptUtil.ivKey = process.env.VITE_AES_IV_KEY;
    return AesEncryptUtil.ivKey;
  }

  private static getAesPublicKey(): string {
    if (AesEncryptUtil.publicKey) return AesEncryptUtil.publicKey;
    if (!process.env.VITE_AES_PUBLIC_KEY) throw new Error("AES PUBLIC KEY SHOULD BE SET");
    AesEncryptUtil.publicKey = process.env.VITE_AES_PUBLIC_KEY;
    return AesEncryptUtil.publicKey;
  }

  static async aesEncrypt(dataToEncrypt: any): Promise<string> {
    dataToEncrypt = JSON.stringify(dataToEncrypt);
    const iv: Buffer = Buffer.from(AesEncryptUtil.getAesIvKey(), "hex");
    const key: Buffer = Buffer.from(AesEncryptUtil.getAesPublicKey(), "hex");
    const cipher: Cipher = createCipheriv("aes-256-ctr", key, iv);
    const encryptedBuffer = Buffer.concat([cipher.update(dataToEncrypt), cipher.final()]);
    return encryptedBuffer.toString("base64");
  }

  static async aesDecrypt(dataToDecrypt: string): Promise<string> {
    const iv: Buffer = Buffer.from(AesEncryptUtil.getAesIvKey(), "hex");
    const key: Buffer = Buffer.from(AesEncryptUtil.getAesPublicKey(), "hex");
    const decipher = createDecipheriv("aes-256-ctr", key, iv);
    const decryptedBuffer = Buffer.concat([
      decipher.update(Buffer.from(dataToDecrypt, "base64")),
      decipher.final(),
    ]);
    return decryptedBuffer.toString();
  }
}
