package stripe

import (
	"fmt"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/balance"
	portalconfig "github.com/stripe/stripe-go/v81/billingportal/configuration"
	portalsession "github.com/stripe/stripe-go/v81/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/invoice"
	"github.com/stripe/stripe-go/v81/invoiceitem"
	"github.com/stripe/stripe-go/v81/paymentlink"
	"github.com/stripe/stripe-go/v81/price"
	"github.com/stripe/stripe-go/v81/product"
	"github.com/stripe/stripe-go/v81/refund"
)

// Executor handles Stripe API operations
type Executor struct {
	apiKey string
}

// NewExecutor creates a new Stripe executor
func NewExecutor(apiKey string) *Executor {
	stripe.Key = apiKey
	return &Executor{
		apiKey: apiKey,
	}
}

// FunctionMap maps operation IDs to their corresponding functions
var FunctionMap = map[string]interface{}{
	"stripe_post_customers":                     (*Executor).CreateCustomer,
	"stripe_get_customers":                      (*Executor).ListCustomers,
	"stripe_post_products":                      (*Executor).CreateProduct,
	"stripe_get_products":                       (*Executor).ListProducts,
	"stripe_post_prices":                        (*Executor).CreatePrice,
	"stripe_get_prices":                         (*Executor).ListPrices,
	"stripe_post_payment_links":                 (*Executor).CreatePaymentLink,
	"stripe_post_invoices":                      (*Executor).CreateInvoice,
	"stripe_post_invoiceitems":                  (*Executor).CreateInvoiceItem,
	"stripe_post_invoices_invoice_finalize":     (*Executor).FinalizeInvoice,
	"stripe_get_balance":                        (*Executor).GetBalance,
	"stripe_post_refunds":                       (*Executor).CreateRefund,
	"stripe_post_products_id":                   (*Executor).UpdateProduct,
	"stripe_get_products_id":                    (*Executor).GetProduct,
	"stripe_post_checkout_sessions":             (*Executor).CreateCheckoutSession,
	"stripe_post_billing_portal_sessions":       (*Executor).CreateBillingPortalSession,
	"stripe_get_prices_price":                   (*Executor).GetPrice,
	"stripe_post_prices_price":                  (*Executor).UpdatePrice,
	"stripe_get_customers_search":               (*Executor).SearchCustomers,
	"stripe_get_customers_customer":             (*Executor).GetCustomer,
	"stripe_get_billing_portal_configurations":  (*Executor).ListBillingPortalConfigurations,
	"stripe_post_billing_portal_configurations": (*Executor).CreateBillingPortalConfiguration,
}

// ExecuteFunction executes a Stripe function by name with given arguments
func (e *Executor) ExecuteFunction(name string, args map[string]interface{}) (interface{}, error) {
	fn, exists := FunctionMap[name]
	if !exists {
		return nil, fmt.Errorf("unknown function: %s", name)
	}

	method := fn.(func(*Executor, map[string]interface{}) (interface{}, error))
	return method(e, args)
}

func (e *Executor) CreateCustomer(params map[string]interface{}) (interface{}, error) {
	p := &stripe.CustomerParams{}
	if email, ok := params["email"].(string); ok {
		p.Email = stripe.String(email)
	}
	if name, ok := params["name"].(string); ok {
		p.Name = stripe.String(name)
	}
	if description, ok := params["description"].(string); ok {
		p.Description = stripe.String(description)
	}
	return customer.New(p)
}

func (e *Executor) ListCustomers(params map[string]interface{}) (interface{}, error) {
	p := &stripe.CustomerListParams{}
	if limit, ok := params["limit"].(float64); ok {
		p.Limit = stripe.Int64(int64(limit))
	}
	i := customer.List(p)
	return collectResults(i)
}

func (e *Executor) CreateProduct(params map[string]interface{}) (interface{}, error) {
	p := &stripe.ProductParams{}
	if name, ok := params["name"].(string); ok {
		p.Name = stripe.String(name)
	}
	if description, ok := params["description"].(string); ok {
		p.Description = stripe.String(description)
	}
	if active, ok := params["active"].(bool); ok {
		p.Active = stripe.Bool(active)
	}
	return product.New(p)
}

func (e *Executor) ListProducts(params map[string]interface{}) (interface{}, error) {
	p := &stripe.ProductListParams{}
	if active, ok := params["active"].(bool); ok {
		p.Active = stripe.Bool(active)
	}
	i := product.List(p)
	return collectResults(i)
}

func (e *Executor) CreatePrice(params map[string]interface{}) (interface{}, error) {
	p := &stripe.PriceParams{}
	if currency, ok := params["currency"].(string); ok {
		p.Currency = stripe.String(currency)
	}
	if productID, ok := params["product"].(string); ok {
		p.Product = stripe.String(productID)
	}
	if unitAmount, ok := params["unit_amount"].(float64); ok {
		p.UnitAmount = stripe.Int64(int64(unitAmount))
	}
	return price.New(p)
}

func (e *Executor) ListPrices(params map[string]interface{}) (interface{}, error) {
	p := &stripe.PriceListParams{}
	if active, ok := params["active"].(bool); ok {
		p.Active = stripe.Bool(active)
	}
	i := price.List(p)
	return collectResults(i)
}

func (e *Executor) CreatePaymentLink(params map[string]interface{}) (interface{}, error) {
	p := &stripe.PaymentLinkParams{}
	if lineItems, ok := params["line_items"].([]interface{}); ok {
		p.LineItems = make([]*stripe.PaymentLinkLineItemParams, len(lineItems))
		for i, item := range lineItems {
			if itemMap, ok := item.(map[string]interface{}); ok {
				lineItem := &stripe.PaymentLinkLineItemParams{}
				if price, ok := itemMap["price"].(string); ok {
					lineItem.Price = stripe.String(price)
				}
				if quantity, ok := itemMap["quantity"].(float64); ok {
					lineItem.Quantity = stripe.Int64(int64(quantity))
				}
				p.LineItems[i] = lineItem
			}
		}
	}
	return paymentlink.New(p)
}

func (e *Executor) CreateInvoice(params map[string]interface{}) (interface{}, error) {
	p := &stripe.InvoiceParams{}
	if customerID, ok := params["customer"].(string); ok {
		p.Customer = stripe.String(customerID)
	}
	if collectionMethod, ok := params["collection_method"].(string); ok {
		p.CollectionMethod = stripe.String(collectionMethod)
	}
	return invoice.New(p)
}

func (e *Executor) CreateInvoiceItem(params map[string]interface{}) (interface{}, error) {
	p := &stripe.InvoiceItemParams{}
	if customerID, ok := params["customer"].(string); ok {
		p.Customer = stripe.String(customerID)
	}
	if amount, ok := params["amount"].(float64); ok {
		p.Amount = stripe.Int64(int64(amount))
	}
	if currency, ok := params["currency"].(string); ok {
		p.Currency = stripe.String(currency)
	}
	return invoiceitem.New(p)
}

func (e *Executor) FinalizeInvoice(params map[string]interface{}) (interface{}, error) {
	id, ok := params["invoice"].(string)
	if !ok {
		return nil, fmt.Errorf("invoice ID is required")
	}
	return invoice.FinalizeInvoice(id, nil)
}

func (e *Executor) GetBalance(params map[string]interface{}) (interface{}, error) {
	return balance.Get(nil)
}

func (e *Executor) CreateRefund(params map[string]interface{}) (interface{}, error) {
	p := &stripe.RefundParams{}
	if charge, ok := params["charge"].(string); ok {
		p.Charge = stripe.String(charge)
	}
	if amount, ok := params["amount"].(float64); ok {
		p.Amount = stripe.Int64(int64(amount))
	}
	return refund.New(p)
}

func (e *Executor) UpdateProduct(params map[string]interface{}) (interface{}, error) {
	id, ok := params["id"].(string)
	if !ok {
		return nil, fmt.Errorf("product ID is required")
	}
	p := &stripe.ProductParams{}
	if name, ok := params["name"].(string); ok {
		p.Name = stripe.String(name)
	}
	if description, ok := params["description"].(string); ok {
		p.Description = stripe.String(description)
	}
	return product.Update(id, p)
}

func (e *Executor) GetProduct(params map[string]interface{}) (interface{}, error) {
	id, ok := params["id"].(string)
	if !ok {
		return nil, fmt.Errorf("product ID is required")
	}
	return product.Get(id, nil)
}

func (e *Executor) CreateCheckoutSession(params map[string]interface{}) (interface{}, error) {
	p := &stripe.CheckoutSessionParams{}
	if mode, ok := params["mode"].(string); ok {
		p.Mode = stripe.String(mode)
	}
	if successURL, ok := params["success_url"].(string); ok {
		p.SuccessURL = stripe.String(successURL)
	}
	if cancelURL, ok := params["cancel_url"].(string); ok {
		p.CancelURL = stripe.String(cancelURL)
	}
	if lineItems, ok := params["line_items"].([]interface{}); ok {
		p.LineItems = make([]*stripe.CheckoutSessionLineItemParams, len(lineItems))
		for i, item := range lineItems {
			if itemMap, ok := item.(map[string]interface{}); ok {
				lineItem := &stripe.CheckoutSessionLineItemParams{}
				if price, ok := itemMap["price"].(string); ok {
					lineItem.Price = stripe.String(price)
				}
				if quantity, ok := itemMap["quantity"].(float64); ok {
					lineItem.Quantity = stripe.Int64(int64(quantity))
				}
				p.LineItems[i] = lineItem
			}
		}
	}
	return checkoutsession.New(p)
}

func (e *Executor) CreateBillingPortalSession(params map[string]interface{}) (interface{}, error) {
	p := &stripe.BillingPortalSessionParams{}
	if customerID, ok := params["customer"].(string); ok {
		p.Customer = stripe.String(customerID)
	}
	if returnURL, ok := params["return_url"].(string); ok {
		p.ReturnURL = stripe.String(returnURL)
	}
	return portalsession.New(p)
}

func (e *Executor) GetPrice(params map[string]interface{}) (interface{}, error) {
	id, ok := params["price"].(string)
	if !ok {
		return nil, fmt.Errorf("price ID is required")
	}
	return price.Get(id, nil)
}

func (e *Executor) UpdatePrice(params map[string]interface{}) (interface{}, error) {
	id, ok := params["price"].(string)
	if !ok {
		return nil, fmt.Errorf("price ID is required")
	}
	p := &stripe.PriceParams{}
	if active, ok := params["active"].(bool); ok {
		p.Active = stripe.Bool(active)
	}
	if nickname, ok := params["nickname"].(string); ok {
		p.Nickname = stripe.String(nickname)
	}
	return price.Update(id, p)
}

func (e *Executor) SearchCustomers(params map[string]interface{}) (interface{}, error) {
	p := &stripe.CustomerSearchParams{}
	if query, ok := params["query"].(string); ok {
		p.Query = query
	}
	i := customer.Search(p)
	return collectResults(i)
}

func (e *Executor) GetCustomer(params map[string]interface{}) (interface{}, error) {
	id, ok := params["customer"].(string)
	if !ok {
		return nil, fmt.Errorf("customer ID is required")
	}
	return customer.Get(id, nil)
}

func (e *Executor) ListBillingPortalConfigurations(params map[string]interface{}) (interface{}, error) {
	p := &stripe.BillingPortalConfigurationListParams{}
	if active, ok := params["active"].(bool); ok {
		p.Active = stripe.Bool(active)
	}
	i := portalconfig.List(p)
	return collectResults(i)
}

func (e *Executor) CreateBillingPortalConfiguration(params map[string]interface{}) (interface{}, error) {
	p := &stripe.BillingPortalConfigurationParams{}
	if businessProfile, ok := params["business_profile"].(map[string]interface{}); ok {
		p.BusinessProfile = &stripe.BillingPortalConfigurationBusinessProfileParams{}
		if privacyURL, ok := businessProfile["privacy_policy_url"].(string); ok {
			p.BusinessProfile.PrivacyPolicyURL = stripe.String(privacyURL)
		}
		if tosURL, ok := businessProfile["terms_of_service_url"].(string); ok {
			p.BusinessProfile.TermsOfServiceURL = stripe.String(tosURL)
		}
	}
	return portalconfig.New(p)
}

// collectResults collects all results from a list iterator
func collectResults(i interface{}) (interface{}, error) {
	var results []interface{}

	// Handle different iterator types
	switch it := i.(type) {
	case *customer.Iter:
		for it.Next() {
			results = append(results, it.Customer())
		}
		return results, it.Err()
	case *product.Iter:
		for it.Next() {
			results = append(results, it.Product())
		}
		return results, it.Err()
	case *price.Iter:
		for it.Next() {
			results = append(results, it.Price())
		}
		return results, it.Err()
	case *customer.SearchIter:
		for it.Next() {
			results = append(results, it.Customer())
		}
		return results, it.Err()
	case *portalconfig.Iter:
		for it.Next() {
			results = append(results, it.BillingPortalConfiguration())
		}
		return results, it.Err()
	default:
		return nil, fmt.Errorf("unsupported iterator type")
	}
}
