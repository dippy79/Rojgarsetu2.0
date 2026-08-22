import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { Mic, Video, PhoneOff } from 'lucide-react';

export default function VideoCallPage() {
  const router = useRouter();
  const { id } = router.query;
  const [roomUrl, setRoomUrl] = useState('');

  useEffect(() => {
    if (id && typeof window !== 'undefined') {
      const token = localStorage.getItem('rojgar_token') || localStorage.getItem('token');
      fetch(`http://localhost:3001/api/v1/interviews/${id}`, {
        headers: { Authorization: `Bearer ${token}` }
      })
        .then(r => r.ok ? r.json() : {})
        .then(data => setRoomUrl(data.room_url || ''));
    }
  }, [id]);

  if (!roomUrl) return <div className="p-8 text-center text-slate-500 min-h-screen flex items-center justify-center">Connecting to video interview room...</div>;

  return (
    <div className="h-screen w-screen bg-slate-900 flex flex-col">
      <div className="flex-1">
        <iframe
          src={roomUrl}
          allow="camera; microphone; fullscreen; display-capture"
          className="w-full h-full border-none"
          title="Interview Video Stream"
        />
      </div>
      <div className="h-20 bg-slate-800 flex items-center justify-center gap-6">
        <button className="p-4 bg-slate-700 hover:bg-slate-600 rounded-full text-white"><Mic className="w-5 h-5" /></button>
        <button className="p-4 bg-slate-700 hover:bg-slate-600 rounded-full text-white"><Video className="w-5 h-5" /></button>
        <button onClick={() => router.back()} className="p-4 bg-rose-600 hover:bg-rose-700 rounded-full text-white">
          <PhoneOff className="w-5 h-5" />
        </button>
      </div>
    </div>
  );
}
