/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_AES_PUBLIC_KEY: string;
  readonly VITE_AES_IV_KEY: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
