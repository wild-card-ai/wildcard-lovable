import { useState, useCallback } from 'react'
import { v4 as uuidv4 } from 'uuid'
import { StreamEvent, ChatState, Message } from '../types/chat'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'
const DELAY_BETWEEN_EVENTS = 500

// Helper to add delay between state updates
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

// Parse a single SSE line into event and data
const parseSSELine = (line: string): { eventType: StreamEvent['type']; data: any } | null => {
  const parts = line.split('\n')
  const eventLine = parts.find(p => p.startsWith('data: '))
  if (!eventLine) return null

  try {
    const eventData = JSON.parse(eventLine.replace('data: ', ''))
    return {
      eventType: eventData.type as StreamEvent['type'],
      data: eventData.data
    }
  } catch (e) {
    console.error('Failed to parse SSE line:', e)
    return null
  }
}

export const useChat = (sessionId: string) => {
  const [state, setState] = useState<ChatState>({
    messages: [],
    isProcessing: false,
    error: null,
    status: [],
    statusMessagesFolded: true
  })

  const toggleStatusFold = useCallback(() => {
    setState(prev => ({
      ...prev,
      statusMessagesFolded: !prev.statusMessagesFolded
    }))
  }, [])

  // Filter status messages based on fold state
  const getFilteredMessages = useCallback((messages: Message[]) => {
    if (!state.statusMessagesFolded) return messages

    // Find consecutive groups of status messages
    const result: Message[] = []
    let currentStatusGroup: Message[] = []
    
    for (const message of messages) {
      if (message.type === 'status') {
        currentStatusGroup.push(message)
      } else {
        // When we hit a non-status message, process any accumulated status group
        if (currentStatusGroup.length > 0) {
          if (currentStatusGroup.length > 2) {
            result.push({
              id: 'folded-indicator',
              type: 'status',
              content: `${currentStatusGroup.length - 2} earlier steps`,
              timestamp: currentStatusGroup[0].timestamp
            })
            result.push(...currentStatusGroup.slice(-2))
          } else {
            result.push(...currentStatusGroup)
          }
          currentStatusGroup = []
        }
        result.push(message)
      }
    }

    // Handle any remaining status group at the end
    if (currentStatusGroup.length > 0) {
      if (currentStatusGroup.length > 2) {
        result.push({
          id: 'folded-indicator',
          type: 'status',
          content: `${currentStatusGroup.length - 2} earlier steps`,
          timestamp: currentStatusGroup[0].timestamp
        })
        result.push(...currentStatusGroup.slice(-2))
      } else {
        result.push(...currentStatusGroup)
      }
    }

    return result
  }, [state.statusMessagesFolded])

  // Handle a single event update
  const handleEventUpdate = useCallback(async (eventType: StreamEvent['type'], data: any) => {
    await delay(DELAY_BETWEEN_EVENTS)

    setState(prev => {
      switch (eventType) {
        case 'start':
        case 'progress':
          return {
            ...prev,
            messages: [...prev.messages, {
              id: uuidv4(),
              type: 'status',
              content: data.message,
              timestamp: new Date()
            }],
            error: null
          }

        case 'error':
          return {
            ...prev,
            error: data.error,
            status: []
          }

        case 'complete':
          if (!data.message) return prev

          return {
            ...prev,
            status: [],
            error: null,
            isProcessing: false,
            messages: [...prev.messages, {
              id: uuidv4(),
              type: 'assistant',
              content: data.message,
              timestamp: new Date()
            }]
          }

        default:
          return prev
      }
    })
  }, [])

  // Process the SSE stream
  const processStream = useCallback(async (response: Response) => {
    const reader = response.body?.getReader()
    if (!reader) throw new Error('No reader available')

    const decoder = new TextDecoder()
    let buffer = ''

    try {
      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n\n')
        buffer = lines.pop() || ''

        for (const line of lines) {
          if (!line.trim()) continue
          
          const event = parseSSELine(line)
          if (event && event.data) {
            await handleEventUpdate(event.eventType, event.data)
          }
        }
      }
    } catch (e) {
      console.error('Error processing stream:', e)
      // Handle stream disconnection error
      setState(prev => ({
        ...prev,
        error: 'Connection lost. Please try again.',
        status: [],
        isProcessing: false
      }))
      throw e
    } finally {
      reader.cancel() // Ensure reader is properly closed
    }
  }, [handleEventUpdate])

  // Send a message and handle the response
  const sendMessage = useCallback(async (content: string) => {
    setState(prev => ({
      ...prev,
      isProcessing: true,
      error: null,
      status: [],
      messages: [...prev.messages, {
        id: uuidv4(),
        type: 'user',
        content,
        timestamp: new Date()
      }]
    }))

    try {
      const response = await fetch(`${API_URL}/process-stream`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          message: content, 
          user_id: sessionId
        })
      })

      if (!response.ok) throw new Error('Failed to process message')
      await processStream(response)
    } catch (err) {
      setState(prev => ({
        ...prev,
        error: err instanceof Error ? err.message : 'An error occurred',
        status: ['Failed to process message'],
        isProcessing: false
      }))
    }
  }, [processStream, sessionId])

  return {
    messages: getFilteredMessages(state.messages),
    isProcessing: state.isProcessing,
    error: state.error,
    status: state.status.join('\n\n'),
    sendMessage,
    toggleStatusFold,
    statusMessagesFolded: state.statusMessagesFolded
  }
} 