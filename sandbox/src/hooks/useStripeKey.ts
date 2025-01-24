import { useState } from 'react';
import { useToast } from '@/components/ui/use-toast';

export function useStripeKey(sessionId: string) {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isKeySet, setIsKeySet] = useState(false);
  const { toast } = useToast();

  const registerStripeKey = async (apiKey: string) => {
    setIsSubmitting(true);
    
    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL}/register-stripe`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          userId: sessionId,
          apiKey: apiKey,
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to register Stripe key');
      }

      localStorage.setItem(`apiKey_${sessionId}`, apiKey);
      setIsKeySet(true);
      toast({
        title: "Success",
        description: "Your Stripe API key has been registered successfully.",
      });
    } catch (error) {
      console.error('Error registering Stripe key:', error);
      toast({
        title: "Error",
        description: "Failed to register your Stripe API key. Please try again.",
        variant: "destructive",
      });
      throw error;
    } finally {
      setIsSubmitting(false);
    }
  };

  const checkExistingKey = () => {
    const storedApiKey = localStorage.getItem(`apiKey_${sessionId}`);
    if (storedApiKey) {
      setIsKeySet(true);
      return storedApiKey;
    }
    return null;
  };

  return {
    isSubmitting,
    isKeySet,
    registerStripeKey,
    checkExistingKey,
  };
} 