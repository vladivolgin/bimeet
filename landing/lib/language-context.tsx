"use client";
import { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { ru } from "@/locales/ru";
import { en } from "@/locales/en";

export type Locale = "ru" | "en";

type Dict = typeof ru;

const dictionaries: Record<Locale, Dict> = { ru, en };

const LanguageContext = createContext<{
  locale: Locale;
  t: Dict;
  setLocale: (l: Locale) => void;
}>({ locale: "ru", t: ru, setLocale: () => {} });

export function LanguageProvider({ children }: { children: ReactNode }) {
  const [locale, setLocaleState] = useState<Locale>("ru");

  useEffect(() => {
    const saved = localStorage.getItem("lang") as Locale | null;
    if (saved === "ru" || saved === "en") {
      setLocaleState(saved);
      document.documentElement.lang = saved;
    }
  }, []);

  function setLocale(l: Locale) {
    setLocaleState(l);
    localStorage.setItem("lang", l);
    document.documentElement.lang = l;
  }

  return (
    <LanguageContext.Provider value={{ locale, t: dictionaries[locale], setLocale }}>
      {children}
    </LanguageContext.Provider>
  );
}

export const useLanguage = () => useContext(LanguageContext);
