# Stripe Agent Sandbox

This is a sandbox application for interacting with Stripe through a chat interface. 
It consists of a React frontend and Go backend.

## Features

The agent supports:
- Creating subscription products with multiple tiers
- Generating checkout links
- Updating product prices
- Creating individual subscription plans

Example commands:
- "Create a recurring pricing model with basic, standard, and pro tier at $10, $50, and $300 per month"
- "Get a stripe checkout link for the standard tier"
- "Change the price of the standard tier to $100/month"
- "Create a new recurring product for $20/month as the Starter plan"

## Setup

1. Install dependencies:
```bash
pnpm install
```

2. Start the development server:
```bash
pnpm dev
```

Note: Ensure the Go backend is running and configured with Stripe API keys.