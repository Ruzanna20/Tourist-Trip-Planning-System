import { useTranslation } from 'react-i18next';

export default function LanguageSelector() {
  const { i18n } = useTranslation();

  return (
    <div className="flex gap-2 p-2 bg-gray-800 rounded-lg">
      <button 
        onClick={() => i18n.changeLanguage('hy')}
        className={`px-2 py-1 text-xs rounded ${i18n.language === 'hy' ? 'bg-brand-600' : 'hover:bg-gray-700'}`}
      >
        🇦🇲 HY
      </button>
      <button 
        onClick={() => i18n.changeLanguage('en')}
        className={`px-2 py-1 text-xs rounded ${i18n.language === 'en' ? 'bg-brand-600' : 'hover:bg-gray-700'}`}
      >
        🇺🇸 EN
      </button>
    </div>
  );
}