export type Message = {
  id: string
  type: 'user' | 'assistant' | 'status' | 'error'
  content: string
  timestamp: Date
}

export type StreamEvent = {
  type: 'start' | 'progress' | 'complete' | 'error'
  data: {
    message?: string
    result?: any
    error?: string
  }
}

export type ChatState = {
  messages: Message[]
  isProcessing: boolean
  error: string | null
  status: string[]
  statusMessagesFolded: boolean
}

export type Example = {
  label: string
  description: string
  prompt: string
} 