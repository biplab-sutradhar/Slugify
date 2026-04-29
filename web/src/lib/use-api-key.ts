"use client";

import { useEffect, useState } from "react";

const STORAGE_KEY = "slugify_api_key";

export function useApiKey() {
  const [apiKey, setApiKey] = useState("");

  useEffect(() => {
    const k = sessionStorage.getItem(STORAGE_KEY);
    if (k) setApiKey(k);
  }, []);

  const save = (value: string) => {
    setApiKey(value);
    sessionStorage.setItem(STORAGE_KEY, value.trim());
  };

  const clear = () => {
    setApiKey("");
    sessionStorage.removeItem(STORAGE_KEY);
  };

  return { apiKey, setApiKey, save, clear };
}