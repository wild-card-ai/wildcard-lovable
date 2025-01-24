import React, { useState, useEffect, FormEvent } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { v4 as uuidv4 } from 'uuid';
import { Chat } from '@/components/chat/Chat';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card } from '@/components/ui/card';

export function StripeChatPage() {
  const navigate = useNavigate();
  const { sessionId } = useParams();
  const [apiKey, setApiKey] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isKeySet, setIsKeySet] = useState(false);

  useEffect(() => {
    if (!sessionId) {
      const newSessionId = uuidv4();
      navigate(`/stripe/${newSessionId}`);
    } else {
      const storedApiKey = localStorage.getItem(`apiKey_${sessionId}`);
      if (storedApiKey) {
        setApiKey(storedApiKey);
        setIsKeySet(true);
      }
    }
  }, [sessionId, navigate]);

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsSubmitting(true);
    localStorage.setItem(`apiKey_${sessionId}`, apiKey);
    setIsKeySet(true);
    setIsSubmitting(false);
  };

  const handleApiKeyChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setApiKey(e.target.value);
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-background py-12">
      <div className="w-full max-w-5xl px-4 md:px-8 space-y-8">
        <div className="text-center space-y-2">
          <h1 className="text-4xl font-bold bg-gradient-to-r from-blue-800 to-blue-600 text-transparent bg-clip-text">
            Stripe Agent Sandbox
          </h1>
          <p className="text-lg text-muted-foreground">
            Interact with your Stripe account using natural language.
          </p>
        </div>

        {!isKeySet ? (
          <Card className="p-6 border shadow-sm">
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <h2 className="text-xl font-semibold">Enter your Stripe API Key</h2>
                <p className="text-sm text-muted-foreground">
                  Your API key will be stored locally and never shared
                </p>
              </div>
              <div className="flex gap-2">
                <Input
                  type="password"
                  placeholder="sk_test_..."
                  value={apiKey}
                  onChange={handleApiKeyChange}
                  className="flex-1"
                />
                <Button 
                  type="submit" 
                  disabled={isSubmitting || !apiKey}
                  className="bg-blue-700 hover:bg-blue-800 text-white"
                >
                  {isSubmitting ? 'Saving...' : 'Save Key'}
                </Button>
              </div>
            </form>
          </Card>
        ) : (
          <Chat sessionId={sessionId!} />
        )}
      </div>
    </div>
  );
} 
