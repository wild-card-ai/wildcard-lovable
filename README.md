# wildcard-lovable

Hosted with love at [lovable.wild-card.ai](https://lovable.wild-card.ai).

## Overview

This sandbox shows how to integrate with Wildcard's Stripe APIs.
To run the go backend server, see [here](go-server/README.md).
To run the frontend, see [here](sandbox/README.md).

## Supported Tools

This demo provides a testing environment for the following Stripe API operations:

### Customers
- stripe_post_customers: Create a customer
- stripe_get_customers: List all customers
- stripe_get_customers_search: Find customers by search
- stripe_get_customers_customer: Get customer details

### Products
- stripe_post_products: Create a product
- stripe_get_products: List all products
- stripe_post_products_id: Modify an existing product
- stripe_get_products_id: Get product details

### Pricing
- stripe_post_prices: Create a price
- stripe_get_prices: List all prices
- stripe_get_prices_price: Get details of a price
- stripe_post_prices_price: Modify an existing price

### Payment & Checkout
- stripe_post_payment_links: Create a payment link
- stripe_post_checkout_sessions: Create a new Checkout Session
- stripe_get_balance: Retrieve balance
- stripe_post_refunds: Create a refund

### Invoices
- stripe_post_invoices: Create an invoice
- stripe_post_invoiceitems: Create an invoice item
- stripe_post_invoices_invoice_finalize: Finalize an invoice

### Billing Portal
- stripe_post_billing_portal_sessions: Create customer portal session
- stripe_get_billing_portal_configurations: Get portal configurations list
- stripe_post_billing_portal_configurations: Create new portal configuration 
