"use client";
import { useLanguage } from "@/lib/language-context";

export default function LanguageSwitcher() {
  const { locale, setLocale } = useLanguage();
  return (
    <button
      onClick={() => setLocale(locale === "ru" ? "en" : "ru")}
      className="text-sm font-semibold text-slate-500 hover:text-indigo-600 transition-colors px-2 py-1 rounded-lg hover:bg-slate-100"
    >
      {locale === "ru" ? "EN" : "RU"}
    </button>
  );
}
