import { Button } from "@/components/ui/button"
import { Example } from "@/types/chat"

const EXAMPLES: Example[] = [
  {
    label: "Create Product Tiers",
    description: "Stripe",
    prompt: "For my SaaS product create a recurring pricing model with a basic, standard, and pro tier at $10, $50, and $300 per month"
  },
  {
    label: "Customer Count",
    description: "Stripe",
    prompt: "Add some dummy customer data and then look up how many customers I have"
  },
  {
    label: "Update Price",
    description: "Stripe",
    prompt: "Create a product called 'Hamburger' at $100/month and then update the price to $150/month"
  },
  {
    label: "Create Single Product",
    description: "Stripe",
    prompt: "Create a new recurring product for $20/month as the Starter plan"
  }
]

interface ExamplePromptsProps {
  onSelect: (prompt: string) => void
}

export function ExamplePrompts({ onSelect }: ExamplePromptsProps) {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
      {EXAMPLES.map((example) => (
        <Button
          key={example.label}
          variant="outline"
          className="relative h-auto p-4 flex flex-col items-start gap-2 group hover:bg-accent/50 bg-muted/50 border-border/50 transition-colors"
          onClick={() => onSelect(example.prompt)}
        >
          <div className="flex flex-col gap-1 w-full">
            <div className="font-medium text-foreground">{example.label}</div>
            <div className="text-xs text-muted-foreground font-light">{example.description}</div>
            <div className="absolute right-4 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 transition-opacity">
              <span className="text-xs text-primary">Try →</span>
            </div>
          </div>
        </Button>
      ))}
    </div>
  )
} 