import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { StripeChatPage } from './pages/stripe';
import { Toaster } from './components/ui/toaster';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Navigate to="/stripe" replace />} />
        <Route path="/stripe/:sessionId?" element={<StripeChatPage />} />
      </Routes>
      <Toaster />
    </Router>
  );
}

export default App;
