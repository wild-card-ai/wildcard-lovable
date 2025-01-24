import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { StripeChatPage } from './pages/stripe';
import { Toaster } from './components/ui/toaster';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/stripe/:sessionId?" element={<StripeChatPage />} />
      </Routes>
      <Toaster />
    </Router>
  );
}

export default App;
