package stripe

import (
	"fmt"
	"reflect"
	"strings"

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
	keyStore StripeKeyStoreInterface
}

// StripeKeyStoreInterface defines the interface for storing and retrieving Stripe API keys
type StripeKeyStoreInterface interface {
	GetStripeKey(userID string) (string, error)
}

// NewExecutor creates a new Stripe executor
func NewExecutor(keyStore StripeKeyStoreInterface) *Executor {
	return &Executor{
		keyStore: keyStore,
	}
}

// setStripeKey sets the Stripe API key for the current operation
func (e *Executor) setStripeKey(userID string) error {
	key, err := e.keyStore.GetStripeKey(userID)
	if err != nil {
		return fmt.Errorf("failed to get Stripe API key for user %s: %v", userID, err)
	}
	stripe.Key = key
	return nil
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
func (e *Executor) ExecuteFunction(userID string, name string, args map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	fn, exists := FunctionMap[name]
	if !exists {
		return nil, fmt.Errorf("unknown function: %s", name)
	}

	method := fn.(func(*Executor, string, map[string]interface{}) (interface{}, error))
	return method(e, userID, args)
}

// convertToStripeParams converts a map[string]interface{} to a Stripe params struct using reflection
func convertToStripeParams(params map[string]interface{}, target interface{}) error {
	targetValue := reflect.ValueOf(target).Elem()
	targetType := targetValue.Type()

	// Handle metadata if it exists
	if metadata, ok := params["metadata"].(map[string]interface{}); ok {
		if method := targetValue.MethodByName("AddMetadata"); method.IsValid() {
			for k, v := range metadata {
				if strVal, ok := v.(string); ok {
					method.Call([]reflect.Value{
						reflect.ValueOf(k),
						reflect.ValueOf(strVal),
					})
				}
			}
		}
	}

	// Handle expand fields if they exist
	if expand, ok := params["expand"].([]interface{}); ok {
		if method := targetValue.MethodByName("AddExpand"); method.IsValid() {
			for _, e := range expand {
				if strVal, ok := e.(string); ok {
					method.Call([]reflect.Value{reflect.ValueOf(strVal)})
				}
			}
		}
	}

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		formTag := field.Tag.Get("form")
		if formTag == "" {
			continue
		}
		formName := strings.Split(formTag, ",")[0]
		if formName == "*" || formName == "-" {
			continue
		}

		if value, ok := params[formName]; ok && value != nil {
			fieldValue := targetValue.Field(i)
			if !fieldValue.CanSet() {
				continue
			}

			switch fieldValue.Type().String() {
			case "*string":
				if strVal, ok := value.(string); ok {
					fieldValue.Set(reflect.ValueOf(stripe.String(strVal)))
				}
			case "*int64":
				switch v := value.(type) {
				case float64:
					fieldValue.Set(reflect.ValueOf(stripe.Int64(int64(v))))
				case int:
					fieldValue.Set(reflect.ValueOf(stripe.Int64(int64(v))))
				}
			case "*bool":
				if boolVal, ok := value.(bool); ok {
					fieldValue.Set(reflect.ValueOf(stripe.Bool(boolVal)))
				}
			case "[]*string":
				if arr, ok := value.([]interface{}); ok {
					strArr := make([]*string, len(arr))
					for i, v := range arr {
						if strVal, ok := v.(string); ok {
							strArr[i] = stripe.String(strVal)
						}
					}
					fieldValue.Set(reflect.ValueOf(strArr))
				}
			case "map[string]string":
				if mapVal, ok := value.(map[string]interface{}); ok {
					strMap := make(map[string]string)
					for k, v := range mapVal {
						if strVal, ok := v.(string); ok {
							strMap[k] = strVal
						}
					}
					fieldValue.Set(reflect.ValueOf(strMap))
				}
			default:
				// Handle nested structs
				if fieldValue.Kind() == reflect.Ptr {
					if nestedMap, ok := value.(map[string]interface{}); ok {
						nestedType := fieldValue.Type().Elem()
						nestedValue := reflect.New(nestedType)
						if err := convertToStripeParams(nestedMap, nestedValue.Interface()); err != nil {
							return err
						}
						fieldValue.Set(nestedValue)
					}
				} else if fieldValue.Kind() == reflect.Struct {
					if nestedMap, ok := value.(map[string]interface{}); ok {
						nestedValue := reflect.New(fieldValue.Type())
						if err := convertToStripeParams(nestedMap, nestedValue.Interface()); err != nil {
							return err
						}
						fieldValue.Set(nestedValue.Elem())
					}
				}
			}
		}
	}
	return nil
}

func (e *Executor) CreateCustomer(userID string, params map[string]interface{}) (interface{}, error) {
	p := &stripe.CustomerParams{}
	if err := convertToStripeParams(params, p); err != nil {
		return nil, err
	}
	return customer.New(p)
}

func (e *Executor) ListCustomers(userID string, params map[string]interface{}) (interface{}, error) {
	p := &stripe.CustomerListParams{}
	if limit, ok := params["limit"].(float64); ok {
		p.Limit = stripe.Int64(int64(limit))
	}
	i := customer.List(p)
	return collectResults(i)
}

func (e *Executor) CreateProduct(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	p := &stripe.ProductParams{}
	if err := convertToStripeParams(params, p); err != nil {
		return nil, err
	}
	fmt.Printf("Product params: %+v\n", p)
	return product.New(p)
}

func (e *Executor) ListProducts(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	p := &stripe.ProductListParams{}
	if active, ok := params["active"].(bool); ok {
		p.Active = stripe.Bool(active)
	}
	i := product.List(p)
	return collectResults(i)
}

func (e *Executor) CreatePrice(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	p := &stripe.PriceParams{}
	if err := convertToStripeParams(params, p); err != nil {
		return nil, err
	}
	return price.New(p)
}

func (e *Executor) ListPrices(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	p := &stripe.PriceListParams{}
	if err := convertToStripeParams(params, p); err != nil {
		return nil, err
	}
	i := price.List(p)
	return collectResults(i)
}

func (e *Executor) CreatePaymentLink(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	p := &stripe.PaymentLinkParams{}
	if err := convertToStripeParams(params, p); err != nil {
		return nil, err
	}
	return paymentlink.New(p)
}

func (e *Executor) CreateInvoice(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	p := &stripe.InvoiceParams{}
	if err := convertToStripeParams(params, p); err != nil {
		return nil, err
	}
	return invoice.New(p)
}

func (e *Executor) CreateInvoiceItem(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	p := &stripe.InvoiceItemParams{}
	if err := convertToStripeParams(params, p); err != nil {
		return nil, err
	}
	return invoiceitem.New(p)
}

func (e *Executor) FinalizeInvoice(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	id, ok := params["invoice"].(string)
	if !ok {
		return nil, fmt.Errorf("invoice ID is required")
	}
	return invoice.FinalizeInvoice(id, nil)
}

func (e *Executor) GetBalance(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	return balance.Get(nil)
}

func (e *Executor) CreateRefund(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	p := &stripe.RefundParams{}
	if err := convertToStripeParams(params, p); err != nil {
		return nil, err
	}
	return refund.New(p)
}

func (e *Executor) UpdateProduct(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	id, ok := params["id"].(string)
	if !ok {
		return nil, fmt.Errorf("product ID is required")
	}
	delete(params, "id")

	p := &stripe.ProductParams{}
	if err := convertToStripeParams(params, p); err != nil {
		return nil, err
	}
	return product.Update(id, p)
}

func (e *Executor) GetProduct(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	id, ok := params["id"].(string)
	if !ok {
		return nil, fmt.Errorf("product ID is required")
	}
	return product.Get(id, nil)
}

func (e *Executor) CreateCheckoutSession(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	p := &stripe.CheckoutSessionParams{}
	if err := convertToStripeParams(params, p); err != nil {
		return nil, err
	}
	return checkoutsession.New(p)
}

func (e *Executor) CreateBillingPortalSession(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	p := &stripe.BillingPortalSessionParams{}
	if err := convertToStripeParams(params, p); err != nil {
		return nil, err
	}
	return portalsession.New(p)
}

func (e *Executor) GetPrice(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	id, ok := params["price"].(string)
	if !ok {
		return nil, fmt.Errorf("price ID is required")
	}
	return price.Get(id, nil)
}

func (e *Executor) UpdatePrice(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	id, ok := params["price"].(string)
	if !ok {
		return nil, fmt.Errorf("price ID is required")
	}
	delete(params, "price")

	p := &stripe.PriceParams{}
	if err := convertToStripeParams(params, p); err != nil {
		return nil, err
	}
	return price.Update(id, p)
}

func (e *Executor) SearchCustomers(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	p := &stripe.CustomerSearchParams{}
	if query, ok := params["query"].(string); ok {
		p.Query = query
	}
	i := customer.Search(p)
	return collectResults(i)
}

func (e *Executor) GetCustomer(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	id, ok := params["customer"].(string)
	if !ok {
		return nil, fmt.Errorf("customer ID is required")
	}
	return customer.Get(id, nil)
}

func (e *Executor) ListBillingPortalConfigurations(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	p := &stripe.BillingPortalConfigurationListParams{}
	if err := convertToStripeParams(params, p); err != nil {
		return nil, err
	}
	i := portalconfig.List(p)
	return collectResults(i)
}

func (e *Executor) CreateBillingPortalConfiguration(userID string, params map[string]interface{}) (interface{}, error) {
	if err := e.setStripeKey(userID); err != nil {
		return nil, err
	}

	p := &stripe.BillingPortalConfigurationParams{}
	if err := convertToStripeParams(params, p); err != nil {
		return nil, err
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
