import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Chat } from '@/components/chat/Chat'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card } from '@/components/ui/card'
import { v4 as uuidv4 } from 'uuid'

export default function StripeChatPage() {
  const { sessionId } = useParams()
  const navigate = useNavigate()
  const [apiKey, setApiKey] = useState('')
  const [isConfigured, setIsConfigured] = useState(false)

  useEffect(() => {
    // If no session ID, generate one and redirect
    if (!sessionId) {
      navigate(`/stripe/${uuidv4()}`)
      return
    }

    // Check for stored API key
    const storedKey = localStorage.getItem(`stripe_key_set_${sessionId}`)
    if (storedKey === 'true') {
      setIsConfigured(true)
    }
  }, [sessionId, navigate])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!apiKey || !sessionId) return

    localStorage.setItem(`stripe_key_set_${sessionId}`, 'true')
    setIsConfigured(true)
  }

  if (!sessionId) {
    return null // Wait for redirect
  }

  if (!isConfigured) {
    return (
      <div className="min-h-screen bg-background flex flex-col items-center justify-center p-4">
        <div className="text-center mb-8">
          <h1 className="mb-3">
            Stripe Chat
          </h1>
          <p className="subtitle max-w-lg">
            Connect with your Stripe account to manage products, prices, and subscriptions using natural language.
          </p>
        </div>
        <Card className="w-full max-w-md p-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <h2>Enter your API Key</h2>
              <p>
                Your API key will be stored locally in your browser.
              </p>
            </div>
            <Input
              type="password"
              placeholder="sk_test_..."
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
            />
            <Button 
              type="submit" 
              className="w-full" 
              disabled={!apiKey}
            >
              Continue
            </Button>
          </form>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background p-4">
      <div className="max-w-4xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1>
            Stripe Chat
          </h1>
          <Button 
            variant="outline"
            onClick={() => {
              localStorage.removeItem(`stripe_key_set_${sessionId}`)
              setIsConfigured(false)
              setApiKey('')
            }}
          >
            Reset API Key
          </Button>
        </div>
        <Chat sessionId={sessionId} />
      </div>
    </div>
  )
} 