"use client";
import { useLanguage } from "@/lib/language-context";

export default function Features() {
  const { t } = useLanguage();
  const f = t.features;
  return (
    <section id="features" className="bg-slate-50 py-24">
      <div className="max-w-5xl mx-auto px-6">
        <div className="text-center mb-16">
          <h2 className="text-4xl font-extrabold text-slate-900 tracking-tight mb-4">
            {f.title}
          </h2>
          <p className="text-lg text-slate-500 max-w-xl mx-auto">{f.subtitle}</p>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-8">
          {f.items.map((item) => (
            <div
              key={item.title}
              className="bg-white rounded-3xl p-8 shadow-sm border border-slate-100 hover:shadow-md hover:-translate-y-0.5 transition-all duration-200"
            >
              <div className="text-4xl mb-4">{item.emoji}</div>
              <h3 className="text-xl font-bold text-slate-900 mb-2">{item.title}</h3>
              <p className="text-slate-500 leading-relaxed text-sm">{item.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
