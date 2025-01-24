import { Message } from '@/types/chat'
import { cn } from '@/lib/utils'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Card } from '@/components/ui/card'
import { AlertCircle, Info } from 'lucide-react'

interface ChatMessageProps {
  message: Message
}

export function ChatMessage({ message }: ChatMessageProps) {
  const isUser = message.type === 'user'
  const isStatus = message.type === 'status'
  const isError = message.type === 'error'

  if (isStatus || isError) {
    return (
      <div className="flex items-center gap-2 text-sm">
        {isStatus ? (
          <div className="flex items-center gap-2 text-primary/60">
            <Info className="h-4 w-4" />
            <span>{message.content}</span>
          </div>
        ) : (
          <div className="flex items-center gap-2 text-destructive">
            <AlertCircle className="h-4 w-4" />
            <span>{message.content}</span>
          </div>
        )}
      </div>
    )
  }

  return (
    <div className={cn(
      'flex gap-3 mb-4',
      isUser ? 'flex-row-reverse' : 'flex-row'
    )}>
      <Avatar className="h-8 w-8 border border-border/50">
        <AvatarFallback className={cn(
          "text-sm font-medium",
          isUser ? "bg-primary/10 text-primary" : "bg-muted text-muted-foreground"
        )}>
          {isUser ? 'U' : 'A'}
        </AvatarFallback>
      </Avatar>
      <div className={cn(
        'flex flex-col gap-1 max-w-[80%]',
        isUser ? 'items-end' : 'items-start'
      )}>
        <Card className={cn(
          'p-4 shadow-md',
          isUser ? 'bg-primary text-primary-foreground' : 'bg-muted/50 text-foreground'
        )}>
          <div className="whitespace-pre-wrap text-sm">{message.content}</div>
        </Card>
        <span className="text-xs text-muted-foreground px-2">
          {new Date(message.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
        </span>
      </div>
    </div>
  )
} 