"use client";
import { useLanguage } from "@/lib/language-context";

export default function Header() {
  const { locale, setLocale } = useLanguage();

  return (
    <div className="relative max-w-5xl mx-auto px-6 pt-6 flex items-center justify-between">
      <div className="flex items-center gap-2">
        <div className="w-7 h-7 rounded-lg bg-white/20 flex items-center justify-center text-xs text-white">
          ✦
        </div>
        <span className="text-white font-bold text-base tracking-tight">Bimeet</span>
      </div>

      <button
        onClick={() => setLocale(locale === "ru" ? "en" : "ru")}
        className="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-white/15 hover:bg-white/25 text-white/90 hover:text-white text-sm font-semibold transition-all duration-200 border border-white/20"
      >
        <span className="text-xs opacity-60">{locale === "ru" ? "🇬🇧" : "🇷🇺"}</span>
        {locale === "ru" ? "EN" : "RU"}
      </button>
    </div>
  );
}
