package models 


// HospitalInventoryItem represents a single item in a hospital's inventory,
// enriched with details from the 'supplies' and 'vendors' tables.
type HospitalInventoryItem struct {
	SupplyID          string    `db:"supply_id" json:"supply_id"`
	SupplyName        string    `db:"supply_name" json:"supply_name"`
	SupplySKU         string    `db:"supply_sku" json:"supply_sku"`
	SupplyCategory    string    `db:"supply_category" json:"supply_category"`
	IsVital           bool      `db:"is_vital" json:"is_vital"`
	CurrentStock      float64   `db:"current_stock" json:"current_stock"`
	ReorderThreshold  float64   `db:"reorder_threshold" json:"reorder_threshold"`
	VendorName        string    `db:"vendor_name" json:"vendor_name"` // From the vendor who supplied it (if you track this per inventory item, otherwise remove)
}