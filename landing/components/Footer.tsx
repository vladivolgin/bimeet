"use client";
import { useLanguage } from "@/lib/language-context";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL ?? "#";

export default function Footer() {
  const { t } = useLanguage();
  const f = t.footer;
  return (
    <footer className="bg-white border-t border-slate-100 py-10">
      <div className="max-w-5xl mx-auto px-6 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div className="flex items-center gap-2">
          <span className="text-xl font-extrabold text-indigo-600">Bimeet</span>
          <span className="text-slate-400 text-sm">{f.tagline}</span>
        </div>

        <div className="flex items-center gap-6 text-sm">
          <a
            href={`${APP_URL}/login`}
            className="text-slate-500 hover:text-indigo-600 transition-colors font-medium"
          >
            {f.login}
          </a>
          <a
            href={`${APP_URL}/register`}
            className="inline-flex items-center px-4 py-2 rounded-xl bg-indigo-600 text-white font-semibold hover:bg-indigo-700 transition-colors"
          >
            {f.register}
          </a>
        </div>
      </div>

      <div className="max-w-5xl mx-auto px-6 mt-6 text-center text-xs text-slate-400">
        © {new Date().getFullYear()} Bimeet. {f.copyright}
      </div>
    </footer>
  );
}
