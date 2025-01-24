import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { StripeChatPage } from '@/pages/stripe';

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/stripe/:sessionId" element={<StripeChatPage />} />
        <Route path="*" element={<Navigate to={`/stripe/${crypto.randomUUID()}`} replace />} />
      </Routes>
    </BrowserRouter>
  );
}
