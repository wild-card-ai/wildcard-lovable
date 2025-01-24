import { useRef, useEffect, useState } from 'react'
import { useChat } from '@/hooks/useChat'
import { ChatMessage } from './ChatMessage'
import { ExamplePrompts } from './ExamplePrompts'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Send } from 'lucide-react'

interface ChatProps {
  sessionId: string
}

export function Chat({ sessionId }: ChatProps) {
  const { messages, isProcessing, error, status, sendMessage } = useChat(sessionId)
  const [inputValue, setInputValue] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)
  const scrollRef = useRef<HTMLDivElement>(null)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const message = inputValue.trim()
    if (!message || isProcessing) return

    sendMessage(message)
    setInputValue('')
  }

  const handleExampleSelect = (prompt: string) => {
    setInputValue(prompt)
    inputRef.current?.focus()
  }

  // Auto scroll to bottom when new messages arrive
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTo({
        top: scrollRef.current.scrollHeight,
        behavior: 'smooth'
      })
    }
  }, [messages, status])

  return (
    <div className="flex flex-col h-[75vh] bg-card rounded-lg border border-border/50 shadow-lg">
      <ScrollArea ref={scrollRef} className="flex-1">
        <div className="py-8 px-4 space-y-6">
          {messages.length === 0 && (
            <div className="space-y-4">
              <h2 className="text-2xl font-medium text-center mb-6">Try an example</h2>
              <ExamplePrompts onSelect={handleExampleSelect} />
            </div>
          )}
          {messages.map(message => (
            <ChatMessage key={message.id} message={message} />
          ))}
          {(isProcessing || status) && (
            <div className="flex items-center gap-2 text-blue-600/80">
              {status ? (
                <div className="flex items-center gap-2">
                  <div className="animate-spin">◌</div>
                  <span>{status}</span>
                </div>
              ) : (
                <>
                  <div className="animate-pulse">●</div>
                  <div className="animate-pulse animation-delay-200">●</div>
                  <div className="animate-pulse animation-delay-400">●</div>
                </>
              )}
            </div>
          )}
          {error && error !== "function execution failed" && (
            <div className="p-4 rounded-md bg-red-500/10 text-red-600 border border-red-200/20">
              {error}
            </div>
          )}
        </div>
      </ScrollArea>

      <div className="border-t border-border/50 bg-muted/50 p-4">
        <form onSubmit={handleSubmit} className="flex gap-3">
          <Input
            ref={inputRef}
            placeholder="Type your message..."
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            disabled={isProcessing}
            className="flex-1 bg-background border-border/50"
          />
          <Button 
            type="submit" 
            disabled={isProcessing || !inputValue.trim()}
            className="bg-blue-700 hover:bg-blue-800 text-white"
          >
            <Send className="h-4 w-4" />
          </Button>
        </form>
      </div>
    </div>
  )
} 
