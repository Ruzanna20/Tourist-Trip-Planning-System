import React, { useState } from 'react';
import Register from './Register'; 
import Login from './Login';

export default function Landing() {
  const [authMode, setAuthMode] = useState(null);

  return (
    <div className="flex h-screen bg-slate-950 text-white overflow-hidden font-sans transition-all duration-500">
      
      <div className={`relative flex flex-col transition-all duration-500 ease-in-out ${authMode ? 'w-full md:w-1/2' : 'w-full'}`}>
        
        <div className="absolute top-[-10%] left-[-10%] w-[60%] h-[60%] bg-blue-600/10 rounded-full blur-[120px] pointer-events-none"></div>

        <nav className="relative z-10 flex justify-between items-center px-10 py-8">
          <div className="flex items-center gap-2 text-2xl font-black tracking-tighter">
            🗺️
          </div>
          {!authMode && (
            <div className="flex gap-4 animate-fade-in">
              <button onClick={() => setAuthMode('login')} className="px-6 py-2.5 font-bold text-sm hover:text-blue-400 transition-all">Մուտք</button>
              <button onClick={() => setAuthMode('register')} className="bg-blue-600 hover:bg-blue-700 px-6 py-2.5 rounded-full font-bold text-sm shadow-lg shadow-blue-600/20 transition-all">Գրանցվել</button>
            </div>
          )}
        </nav>

        <main className="relative z-10 flex-grow flex flex-col items-center justify-center text-center px-6 pb-20">
          <h1 className={`font-black mb-6 leading-tight tracking-tighter italic transition-all duration-500 ${authMode ? 'text-4xl md:text-6xl' : 'text-6xl md:text-8xl'}`}>
            Պլանավորիր <br /> <span className="text-blue-500">Ուղևորությունդ</span>
          </h1>
        </main>
      </div>

      <div className={`relative flex flex-col bg-white text-slate-900 transition-all duration-500 ease-in-out border-l border-white/10 ${authMode ? 'w-full md:w-1/2 opacity-100' : 'w-0 opacity-0 overflow-hidden'}`}>
        
        {authMode && (
          <button 
            onClick={() => setAuthMode(null)}
            className="absolute top-8 right-8 p-2 text-slate-400 hover:text-slate-900 transition-all z-[110]"
          >
            <span className="text-2xl">✕</span>
          </button>
        )}
        
        <div className="h-full overflow-y-auto p-8 md:p-16 flex flex-col justify-center bg-slate-50">
          <div className="max-w-md mx-auto w-full">
            {authMode === 'login' ? (
                <Login onSwitch={() => setAuthMode('register')} />
                ) : (
                <Register onSwitch={() => setAuthMode('login')} />
                )}
            </div>
          </div>
        </div>
      </div>
  );
}